//Package spidermgr
package spidermgr

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/lewgun/argyroneta/pkg/types"
	"github.com/lewgun/argyroneta/pkg/cache"
	"github.com/lewgun/argyroneta/pkg/cache/memory"

	"github.com/oli-g/chuper"
	"github.com/lewgun/argyroneta/pkg/errutil"
	"github.com/Sirupsen/logrus"
	"github.com/PuerkitoBio/goquery"
)

const (
	NetEase types.Domain = "netease"
)

const (
	HTTP_POST = "POST"
	HTTP_GET  = "GET"
)

type HTMLHandler func(chuper.Context, *goquery.Document) bool

type Spider struct {
	*chuper.Crawler
	chuper.Enqueuer
}

var SM *SpiderMgr

func init() {
	SM = &SpiderMgr{}
}

type SpiderMgr struct {
	//about spiders
	spiders  map[types.Domain]*Spider
	handlers map[types.Domain]HTMLHandler
	rules    map[types.Domain]*types.Site

	urlPool cache.Cache
	logger  *logrus.Logger

	mu sync.Mutex
	wg sync.WaitGroup
}

func (sm *SpiderMgr) register(domain types.Domain, maker HTMLHandler) error {

	sm.mu.Lock()
	defer sm.mu.Unlock()
	if _, ok := sm.handlers[domain]; ok {
		return errutil.ErrAlreadyExisted
	}
	sm.handlers[domain] = maker
	return nil

}
func (sm *SpiderMgr) PowerOff() error {

	for _, spider := range sm.spiders {
		sm.wg.Add(1)
		go func() {
			defer sm.wg.Done()
			spider.Crawler.Finish()
			spider.Crawler.Block()
		}()

	}

	sm.wg.Wait()
	return nil

}

func (sm *SpiderMgr) Init(rules map[types.Domain]*types.Site, logger *logrus.Logger) []error {
	if rules == nil {
		return errutil.ErrInvalidParameter
	}

	sm.urlPool = memory.New()

	//todo immutable
	sm.rules = rules

	sm.logger = logger

	//todo err
	var (
		ok bool
		h  HTMLHandler

		errs = make([]error, 0, len(rules))
	)

	for domain, _ := range rules {
		if h, ok = sm.handlers[domain]; !ok {
			errs = append(errs, fmt.Errorf("no html handler for '%s' ", domain))
			continue
		}

		spider, err := sm.spiderMaker(domain, h)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		sm.spiders[domain] = spider

	}

	return errs

}

func (sm *SpiderMgr) spiderMaker(domain string, h HTMLHandler) (*Spider, error) {

	//rule *types.Site, urlPool cache.Cache, logger *logrus.Logger,

	rule, ok := sm.rules[domain]
	if !ok {
		return nil, errutil.ErrNotFound
	}

	parsedURL, err := url.Parse(rule.Seed)
	if err != nil {
		return nil, err
	}

	criteria := &chuper.ResponseCriteria{
		//Method:      "GET",
		ContentType: "text/html",
		Status:      200,
		Host:        parsedURL.Host,
	}

	crawler := chuper.New()
	crawler.CrawlDelay = rule.Delay * 1e6
	crawler.UserAgent = rule.UserAgent
	crawler.Cache = sm.urlPool
	crawler.CrawlPoliteness = rule.Politeness
	crawler.Logger = sm.logger

	crawler.Register(criteria, chuper.ProcessorFunc(h))

	return &Spider{
		Crawler:  crawler,
		Enqueuer: crawler.Start(),
	}, nil

}

//"GET", u, "www.google.com", depth
func (sm *SpiderMgr) PowerOn() error {

	for domain, spider := range sm.spiders {
		spider.Enqueuer.Enqueue(HTTP_GET, sm.rules[domain].Seed, "", sm.rules[domain].MaxDepth)

	}

}

//Register 注册一个蜘蛛生成器
func Register(domain string, h HTMLHandler) error {
	return SM.register(domain, h)

}

func SharedInstance() *SpiderMgr {
	if SM == nil {
		panic("never happen")
	}
	return SM
}
