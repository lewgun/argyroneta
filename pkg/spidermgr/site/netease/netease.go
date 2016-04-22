package netease

import (
	"github.com/oli-g/chuper"

	"github.com/PuerkitoBio/goquery"
	"github.com/lewgun/argyroneta/pkg/spidermgr"
	//"github.com/Sirupsen/logrus"
	"github.com/lewgun/argyroneta/pkg/constants"
)

func handler(ctx chuper.Context, doc *goquery.Document) bool {
	//ctx.Log(map[string]interface{}{
	//"url":    ctx.URL().String(),
	//"source": ctx.SourceURL().String(),
	//}).Info("First processor")

	return true
}

func init() {
	spidermgr.Register(constants.NetEase, handler)
}
