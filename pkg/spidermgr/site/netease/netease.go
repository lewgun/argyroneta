package netease

import (
	//"fmt"
	"bytes"
	"io/ioutil"
	//"strings"

	"github.com/oli-g/chuper"

	"github.com/PuerkitoBio/goquery"

	"github.com/lewgun/argyroneta/pkg/constants"
	"github.com/lewgun/argyroneta/pkg/spidermgr"

	//"github.com/lewgun/argyroneta/pkg/types"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	//	"github.com/lewgun/argyroneta/pkg/types"
	//"time"
	//	"time"
	//"regexp"
	//"strings"
	"fmt"
)

func gbk2Utf8(raw []byte) ([]byte, error) {
	rGBK := bytes.NewReader(bytes.TrimSpace(raw))
	rUTF8 := transform.NewReader(rGBK, simplifiedchinese.GBK.NewDecoder())
	return ioutil.ReadAll(rUTF8)

}

var handlerMap map[int]spidermgr.HTMLHandler

//热点新闻
func handlerHotNews(ctx chuper.Context, query *goquery.Document) bool {
//	// 若有多页内容，则获取阅读全文的链接并获取内容
//	if pageAll := query.Find(".ep-pages-all"); len(pageAll.Nodes) != 0 {
//		if pageAllUrl, ok := pageAll.Attr("href"); ok {
//			ctx.AddQueue(&request.Request{
//				Url:  pageAllUrl,
//				Rule: "热点新闻",
//				Temp: ctx.CopyTemps(),
//			})
//		}
//		return
//	}

	// 获取标题
	title := query.Find("#h1title").Text()

	t1, _ := gbk2Utf8([]byte(title))
	fmt.Println(string(t1))

//	// 获取内容
//	content := query.Find("#endText").Text()
//	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
//	// content = re.ReplaceAllStringFunc(content, strings.ToLower)
//	content = re.ReplaceAllString(content, "")
//
//	// 获取发布日期
//	release := query.Find(".ep-time-soure").Text()
//	release = strings.Split(release, "来源:")[0]
//	release = strings.Trim(release, " \t\n")

	//fmt.Println(title, release)


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
	newsType := doc.Find(".titleBar h2").Text()

	topUrls := map[string]string{}

	doc.Find(".tabContents").Each(func(n int, t *goquery.Selection) {
		t.Find("tr").Each(func(i int, s *goquery.Selection) {
			// 跳过标题栏
			if i == 0 {
				return
			}
			// 内容链接
			url, ok := s.Find("a").Attr("href")

			// 排名
			top := s.Find(".cBlue").Text()

			if ok  && !ctx.Cache().Has(url){
					ctx.Cache().Set(url, true)
					topUrls[url] += topTitle[n] + ":" + top + ","

			}
		})
	})

	for k, v := range topUrls {

		extras := map[string]interface{}{
			"newsType": newsType,
			"top":      v,
		}
		ctx.Queue().Enqueue(constants.HTTP_GET, k, ctx.Depth()-1, extras)

	}


	return true
}

//排行榜主页
func handlerRankIndex(ctx chuper.Context, doc *goquery.Document) bool {
	//	var selection = doc.Find("div.tabBox div.tabContents.active a")
	var selection = doc.Find(".subNav a")
	selection.Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		//println(url)

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

	spidermgr.Register(constants.NetEase, handler)
}
