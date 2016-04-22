//Package spidermgr
package spidermgr

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/lewgun/argyroneta/pkg/cache"
	"github.com/lewgun/argyroneta/pkg/cache/memory"
	"github.com/lewgun/argyroneta/pkg/types"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
	"github.com/lewgun/argyroneta/pkg/errutil"
	"github.com/oli-g/chuper"
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
	SM = &SpiderMgr{
		spiders:  make(map[types.Domain]*Spider),
		handlers: make(map[types.Domain]HTMLHandler),
	}
}

type SpiderMgr struct {
	//about spiders
	spiders  map[types.Domain]*Spider
	handlers map[types.Domain]HTMLHandler
	rules    map[types.Domain]*types.Site

	chanExit <-chan struct{}

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

	for domain, spider := range sm.spiders {
		sm.wg.Add(1)
		go func() {
			defer sm.wg.Done()
			sm.logger.Infof("wait for spider for %s finished\n", domain)
			spider.Crawler.Finish()
			spider.Crawler.Block()
		}()

	}

	sm.wg.Wait()

	sm.logger.Infof("spider manager is power off now")
	return nil

}

func (sm *SpiderMgr) Init(rules map[types.Domain]*types.Site, logger *logrus.Logger) []error {
	if rules == nil {
		return []error{errutil.ErrInvalidParameter}
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
			err := fmt.Errorf("no html handler for '%s' ", domain)
			sm.logger.Errorln(err)

			errs = append(errs, err)
			continue
		}

		spider, err := sm.spiderMaker(domain, h)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		sm.spiders[domain] = spider
	}

	sm.logger.Info("spider manager's init is finished")
	return errs

}

func (sm *SpiderMgr) spiderMaker(domain types.Domain, h HTMLHandler) (*Spider, error) {

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
	crawler.CrawlDelay = rule.Delay * time.Second
	crawler.UserAgent = rule.UserAgent
	crawler.Cache = sm.urlPool
	crawler.CrawlPoliteness = rule.Politeness
	crawler.Logger = sm.logger

	crawler.Register(criteria, chuper.ProcessorFunc(h))

	sm.logger.Infof("spider for %s is running", domain)

	return &Spider{
		Crawler:  crawler,
		Enqueuer: crawler.Start(),
	}, nil

}

//"GET", u, "www.google.com", depth
func (sm *SpiderMgr) PowerOn() error {

	for domain, spider := range sm.spiders {
		spider.Enqueuer.Enqueue(HTTP_GET, sm.rules[domain].Seed, "", sm.rules[domain].MaxDepth)
		sm.logger.Infof("seed for %s is putted with max depth: %d", domain, sm.rules[domain].MaxDepth)

	}
	return nil

}

//Register 注册一个蜘蛛生成器
func Register(domain types.Domain, h HTMLHandler) error {
	return SM.register(domain, h)

}

func SharedInstance() *SpiderMgr {
	if SM == nil {
		panic("never happen")
	}
	return SM
}
