package netease

import (
	"github.com/oli-g/chuper"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
)

func handler(ctx chuper.Context, doc *goquery.Document) bool {
	//ctx.Log(map[string]interface{}{
	//"url":    ctx.URL().String(),
	//"source": ctx.SourceURL().String(),
	//}).Info("First processor")

	return true
}
