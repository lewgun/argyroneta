package netease

import (
	"github.com/oli-g/chuper"

	"github.com/PuerkitoBio/goquery"
	"github.com/lewgun/argyroneta/pkg/spidermgr"
	//"github.com/Sirupsen/logrus"
	"github.com/lewgun/argyroneta/pkg/constants"

	"github.com/lewgun/argyroneta/pkg/types"
)

var handlerMap map[int]spidermgr.HTMLHandler

func handlerDepth0(ctx chuper.Context, doc *goquery.Document) bool {

	ctx.Log(nil).Info("handlerDepth0")

	rule, ok := ctx.Extra().(*types.Site)
	if !ok {
		return false
	}

	ctx.Queue().Enqueue(constants.HTTP_GET, rule.Seed, "", ctx.Depth()-1)

	return true
}
func handlerDepth1(ctx chuper.Context, doc *goquery.Document) bool {
	ctx.Log(nil).Info("handlerDepth1")

	rule, ok := ctx.Extra().(*types.Site)
	if !ok {
		return false
	}

	ctx.Queue().Enqueue(constants.HTTP_GET, rule.Seed, "", ctx.Depth()-1)

	return true
}
func handlerDepth2(ctx chuper.Context, doc *goquery.Document) bool {
	ctx.Log(nil).Info("handlerDepth2")

	rule, ok := ctx.Extra().(*types.Site)
	if !ok {
		return false
	}
	
	var selection = doc.Find("div.tabBox div.tabContents.active a")
	selection.Each(func(i int, s *goquery.Selection) {
		val, _ := s.Attr("href")
		println(val)
	}
	
	

	ctx.Queue().Enqueue(constants.HTTP_GET, rule.Seed, "", ctx.Depth()-1)
	return true
}

/*
	r := &Rule{
		Ident:    "tmcnet",
		Start:    "http://www.tmcnet.com/voip/",
		CSSLinks: "a[href]",
		CSSTitle: "title",
		Restart:  30 * time.Minute,
		Accept: []*regexp.Regexp{
			regexp.MustCompile(`^/voip/(departments|columns|features)/articles/`),
		},
		Reject: []*regexp.Regexp{
			regexp.MustCompile(`bad link`),
		},
	}
*/
// func handler(ctx chuper.Context, doc *goquery.Document) bool {
// 	ctx.Log(map[string]interface{}{
// 	"url":    ctx.URL().String(),
// 	"source": ctx.SourceURL().String(),
// 	"depth": ctx.Depth(),
// 	}).Info("First processor")

// 	//已经挖到极限深度
// 	if ctx.Depth() <= 0 {
// 		return true
// 	}

// 	rule, ok := ctx.Extra.(*types.Site)
// 	if !ok {
// 		return false
// 	}

// 	rule.ExtractLinks()

// 	return true
// }

func handler(ctx chuper.Context, doc *goquery.Document) bool {
	ctx.Log(map[string]interface{}{
		"url":    ctx.URL().String(),
		"source": ctx.SourceURL().String(),
		"depth":  ctx.Depth(),
	}).Info("First processor")

	h, ok := handlerMap[ctx.Depth()]
	if !ok {
		return false
	}

	return h(ctx, doc)

}

func init() {

	handlerMap = make(map[int]spidermgr.HTMLHandler)

	handlerMap[0] = handlerDepth0
	handlerMap[1] = handlerDepth1
	handlerMap[2] = handlerDepth2

	spidermgr.Register(constants.NetEase, handler)
}
