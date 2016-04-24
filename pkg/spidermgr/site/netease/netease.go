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
)

func gbk2Utf8(raw []byte) ([]byte, error) {
	rGBK := bytes.NewReader(bytes.TrimSpace(raw))
	rUTF8 := transform.NewReader(rGBK, simplifiedchinese.GBK.NewDecoder())
	return ioutil.ReadAll(rUTF8)

}

var handlerMap map[int]spidermgr.HTMLHandler

//新闻排行榜
func handlerRankNews(ctx chuper.Context, doc *goquery.Document) bool {
	//	topTitle := []string{
	//		"1小时前点击排行",
	//		"24小时点击排行",
	//		"本周点击排行",
	//		"今日跟帖排行",
	//		"本周跟帖排行",
	//		"本月跟贴排行",
	//	}
	//
	//	// 获取新闻分类
	//	newsType := doc.Find(".titleBar h2").Text()
	//
	////	topURLs := map[string]string{}
	////
	////
	//	doc.Find(".tabContents").Each(func(n int, t *goquery.Selection) {
	//		t.Find("tr").Each(func(i int, s *goquery.Selection) {
	//			// 跳过标题栏
	//			if i == 0 {
	//				return
	//			}
	//			// 内容链接
	//			url, ok := s.Find("a").Attr("href")
	//
	//			// 排名
	//			top := s.Find(".cBlue").Text()
	//
	//			if ok {
	//				topURLs[url] += topTitle[n] + ":" + top + ","
	//				println(topURLs[url])
	//			}
	//		})
	//	})

	//	for k, v := range topURLs {
	//		ctx.AddQueue(&request.Request{
	//			Url:  k,
	//			Rule: "热点新闻",
	//			Temp: map[string]interface{}{
	//				"newsType": newsType,
	//				"top":      v,
	//			},
	//		})
	//	}

	return true
}

//排行榜主页
func handlerRankIndex(ctx chuper.Context, doc *goquery.Document) bool {
	//	var selection = doc.Find("div.tabBox div.tabContents.active a")
	var selection = doc.Find(".subNav a")
	selection.Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		println(url)

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

	handlerMap[1] = handlerRankNews
	handlerMap[2] = handlerRankIndex

	spidermgr.Register(constants.NetEase, handler)
}
