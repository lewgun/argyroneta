package netease

import (
	"bytes"
	"io/ioutil"
	"os"
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
	content := query.Find("#endText").Text()
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	// content = re.ReplaceAllStringFunc(content, strings.ToLower)

	cid, err := bolt.B.SaveBlob(types.Blob(re.ReplaceAll([]byte(content), []byte(""))))
	if err != nil {
		fmt.Println(err)
		//todo err check
		return false
	}

	// 获取标题
	title := query.Find(".post_content_main h1").Text()
	utf8Title, _ := gbk2Utf8([]byte(title))

	// 获取发布日期
	source := query.Find("div.post_time_source").Text()
	utf8Source, _ := gbk2Utf8([]byte(source))
	fields := bytes.Split(utf8Source, []byte("来源:"))

	publishAt, _ := time.Parse("2006-01-02 15:04:05", string(bytes.TrimSpace(fields[0])))
	entry := &types.Entry{
		URL:       ctx.URL().String(),
		Title:     string(utf8Title),
		PublishAt: publishAt,
		Origin:    string(bytes.TrimSpace(fields[1])),
		FetchAt:   time.Now(),
		ContentID: string(cid),
	}

	if params, ok := ctx.Extras().(map[string]interface{}); ok {
		entry.Origin = params["category"].(string)
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
	category := doc.Find(".titleBar h2").Text()

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
			"category": category,
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

	defer func() {
		if e := recover(); e != nil {
			f, _ := os.OpenFile("./abcd.txt", 666, os.ModePerm)
			f.WriteString(e.(error).Error())
			f.Sync()
			f.Close()

		}
	}()

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

	spidermgr.Register(constants.NetEase, handler)
}
