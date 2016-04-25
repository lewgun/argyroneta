package netease

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"time"

	"github.com/oli-g/chuper"

	"github.com/PuerkitoBio/goquery"

	"github.com/lewgun/argyroneta/pkg/constants"
	"github.com/lewgun/argyroneta/pkg/deepcopy"
	"github.com/lewgun/argyroneta/pkg/spidermgr"
	"github.com/lewgun/argyroneta/pkg/store/bolt"
	"github.com/lewgun/argyroneta/pkg/store/mysql"
	"github.com/lewgun/argyroneta/pkg/types"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func gbk2Utf8(raw []byte) ([]byte, error) {
	rGBK := bytes.NewReader(bytes.TrimSpace(raw))
	rUTF8 := transform.NewReader(rGBK, simplifiedchinese.GBK.NewDecoder())
	return ioutil.ReadAll(rUTF8)

}

var handlerMap map[int]spidermgr.HTMLHandler
var reNewsContent *regexp.Regexp

//热点新闻
func handlerHotNews(ctx chuper.Context, query *goquery.Document) bool {

	// 若有多页内容，则获取阅读全文的链接并获取内容
	if pageAll := query.Find(".ep-pages-all"); len(pageAll.Nodes) != 0 {
		if pageAllUrl, ok := pageAll.Attr("href"); ok {

			extras := make(map[string]interface{})
			if err := deepcopy.Copy(&extras, ctx.Extras()); err == nil {
				ctx.Queue().Enqueue(constants.HTTP_GET, pageAllUrl, ctx.Depth(), extras)
			}

		}
		return true
	}

	// 获取内容
	// content = re.ReplaceAllStringFunc(content, strings.ToLower)
	content, _ := gbk2Utf8([]byte(query.Find("#endText").Text()))

	cid, err := bolt.B.SaveBlob(types.Blob(reNewsContent.ReplaceAll(content, []byte(""))))
	if err != nil {
		//todo err check
		return false
	}

	// 获取标题
	title := query.Find(".post_content_main h1").Text()
	utf8Title, _ := gbk2Utf8([]byte(title))

	//过滤掉空title的文章
	if utf8Title == nil {
		return false
	}

	// 获取发布日期
	source := query.Find("div.post_time_source").Text()
	utf8Source, _ := gbk2Utf8([]byte(source))
	fields := bytes.Split(utf8Source, []byte("来源:"))

	var (
		publishAt time.Time
		origin    string
	)
	if len(fields) >= 2 {
		publishAt, _ = time.Parse("2006-01-02 15:04:05", string(bytes.TrimSpace(fields[0])))
		origin = string(bytes.TrimSpace(fields[1]))
	}

	entry := &types.Entry{
		URL:       ctx.URL().String(),
		Title:     string(utf8Title),
		PublishAt: publishAt,
		Origin:    origin,
		FetchAt:   time.Now(),
		ContentID: string(cid),
	}

	if params, ok := ctx.Extras().(map[string]interface{}); ok {
		entry.Category = params["category"].(string)
		entry.Summary = params["top"].(string)
	}

	mysql.M.AddEntry(entry)

	return true
}

//新闻排行榜
func handlerRankNews(ctx chuper.Context, doc *goquery.Document) bool {
	topTitle := []string{
		"1小时前点击排行",
		"24小时点击排行",
		"本周点击排行",
		"今日跟帖排行",
		"本周跟帖排行",
		"本月跟贴排行",
	}
	// 获取新闻分类
	category, _ := gbk2Utf8([]byte(doc.Find(".titleBar h2").Text()))

	topUrls := map[string]string{}

	doc.Find(".tabContents").Each(func(n int, t *goquery.Selection) {
		t.Find("tr").Each(func(i int, s *goquery.Selection) {
			// 跳过标题栏
			if i == 0 {
				return
			}

			// 内容链接
			url, ok := s.Find("a").Attr("href")
			if !ok {
				return
			}

			// 排名
			top := s.Find(".cBlue").Text()

			//同一url在不同榜单中的排位
			topUrls[url] += topTitle[n] + ":" + top + ","

		})
	})

	for url, top := range topUrls {

		extras := map[string]interface{}{
			"category": string(category),
			"top":      top,
		}

		if ctx.Cache().Has(url) {
			continue
		}

		ctx.Cache().Set(url, true)

		ctx.Queue().Enqueue(constants.HTTP_GET, url, ctx.Depth()-1, extras)

	}

	return true
}

//排行榜主页
func handlerRankIndex(ctx chuper.Context, doc *goquery.Document) bool {
	//	var selection = doc.Find("div.tabBox div.tabContents.active a")
	var selection = doc.Find(".subNav a")
	selection.Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")

		if !ctx.Cache().Has(url) {
			ctx.Cache().Set(url, true)
			ctx.Queue().Enqueue(constants.HTTP_GET, url, ctx.Depth()-1)

		}

	})

	return true
}

func handler(ctx chuper.Context, doc *goquery.Document) bool {

	h, ok := handlerMap[ctx.Depth()]
	if !ok {
		return false
	}

	return h(ctx, doc)

}

func init() {

	handlerMap = make(map[int]spidermgr.HTMLHandler)

	handlerMap[0] = handlerHotNews
	handlerMap[1] = handlerRankNews
	handlerMap[2] = handlerRankIndex

	reNewsContent = regexp.MustCompile("\\<[\\S\\s]+?\\>")

	spidermgr.Register(constants.NetEase, handler)
}
