package main

import (
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/lewgun/argyroneta/pkg/store"
	_ "github.com/lewgun/argyroneta/pkg/store/bolt"
    
    "github.com/lewgun/argyroneta/pkg/cmdmgr"
    

	//"github.com/oli-g/chuper"
	"github.com/lewgun/argyroneta/vendor/github.com/oli-g/chuper"
	"fmt"
)

var (
	delay = 2 * time.Second

	depth = 0


	criteria = &chuper.ResponseCriteria{
		Method:      "GET",
		ContentType: "text/html",
		Status:      200,
		Host:        "www.gazzetta.it",
	}

	firstProcessor = chuper.ProcessorFunc(func(ctx chuper.Context, doc *goquery.Document) bool {
		ctx.Log(map[string]interface{}{
			"url":    ctx.URL().String(),
			"source": ctx.SourceURL().String(),
		}).Info("First processor")
		return true
	})

	secondProcessor = chuper.ProcessorFunc(func(ctx chuper.Context, doc *goquery.Document) bool {
		ctx.Log(map[string]interface{}{
			"url":    ctx.URL().String(),
			"source": ctx.SourceURL().String(),
		}).Info("Second processor")
		return false

	})

	thirdProcessor = chuper.ProcessorFunc(func(ctx chuper.Context, doc *goquery.Document) bool {
		ctx.Log(map[string]interface{}{
			"url":    ctx.URL().String(),
			"source": ctx.SourceURL().String(),
		}).Info("Third processor")
		return true
	})
)

func bootUp(s store.Store, q *chuper.Queue) error {

	rules, err := s.Rules()
	if err != nil {
		return err
	}
    
    err = cmdmgr.BootUp(rules)
    if err != nil {
        return err 
    }

	for ident, seed := range seeds {


	}

	for _, u := range seeds {
		q.Enqueue("GET", u, "www.google.com", depth)
		depth++
	}

	return nil,
}
func main() {
	crawler := chuper.New()
	crawler.CrawlDelay = delay

	s, err := store.Open( store.Bolt, "./test.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer s.Close()

	crawler.Register(criteria, firstProcessor, secondProcessor, thirdProcessor)
	q := crawler.Start()

	bootUp(s, q)

	crawler.Finish()
}
