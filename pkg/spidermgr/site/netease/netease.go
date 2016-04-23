package netease

import (
	"bytes"
	"io/ioutil"

	"github.com/oli-g/chuper"

	"github.com/PuerkitoBio/goquery"

	"github.com/lewgun/argyroneta/pkg/constants"
	"github.com/lewgun/argyroneta/pkg/spidermgr"

	//"github.com/lewgun/argyroneta/pkg/types"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func gbk2Utf8(raw []byte) ([]byte, error) {
	rGBK := bytes.NewReader(raw)
	rUTF8 := transform.NewReader(rGBK, simplifiedchinese.GBK.NewDecoder())
	return ioutil.ReadAll(rUTF8)

}

var handlerMap map[int]spidermgr.HTMLHandler

func handlerDepth1(ctx chuper.Context, doc *goquery.Document) bool {
	ctx.Log(nil).Info("handlerDepth1")

	// rule, ok := ctx.Extra().(*types.Site)
	// if !ok {
	// 	return false
	// }

	var selection = doc.Find("div.endContent")
	selection.Each(func(i int, s *goquery.Selection) {
		title := s.Find("h1#h1title").Text()
		// url, _ := s.Attr("href")
		// println(url)

		utf8Txt, _ := gbk2Utf8([]byte(title))
		println(string(utf8Txt))
		// println(title)

		//ctx.Queue().Enqueue(constants.HTTP_GET, url, "", ctx.Depth()-1)
	})

	//ctx.Queue().Enqueue(constants.HTTP_GET, rule.Seed, "", ctx.Depth()-1)

	return true
}
func handlerDepth2(ctx chuper.Context, doc *goquery.Document) bool {
	ctx.Log(nil).Info("handlerDepth2")

	// rule, ok := ctx.Extra().(*types.Site)
	// if !ok {
	// 	return false
	// }

	var selection = doc.Find("div.tabBox div.tabContents.active a")
	selection.Each(func(i int, s *goquery.Selection) {

		url, _ := s.Attr("href")
		println(url)

		utf8Txt, _ := gbk2Utf8([]byte(s.Text()))
		println(string(utf8Txt))
		println()

		if !ctx.Cache().Has(url) {
			ctx.Cache().Set(url, true)
			ctx.Queue().Enqueue(constants.HTTP_GET, url, "", ctx.Depth()-1)
		} else {
			ctx.Log(nil).Errorf("dup url: %s ", url)
		}
	})

	return true
}

func handler(ctx chuper.Context, doc *goquery.Document) bool {
	// ctx.Log(map[string]interface{}{
	// 	"url":    ctx.URL().String(),
	// 	"source": ctx.SourceURL().String(),
	// 	"depth":  ctx.Depth(),
	// }).Info("First processor")

	h, ok := handlerMap[ctx.Depth()]
	if !ok {
		return false
	}

	return h(ctx, doc)

}

func init() {

	handlerMap = make(map[int]spidermgr.HTMLHandler)

	handlerMap[1] = handlerDepth1
	handlerMap[2] = handlerDepth2

	spidermgr.Register(constants.NetEase, handler)
}
