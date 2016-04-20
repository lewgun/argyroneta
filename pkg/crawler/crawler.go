package crawler

import (
    
    "os"
    "net/http"
    
	"github.com/lewgun/argyroneta/pkg/cache"
	"github.com/lewgun/argyroneta/pkg/cache/memory"   
    
	"github.com/PuerkitoBio/fetchbot"
    "github.com/Sirupsen/logrus"
)

const (
    DefaultCrawlPoliteness = false
	DefaultLogFormat       = "text"
	DefaultLogLevel        = "info"
    
    
)

var (
	DefaultHTTPClient = http.DefaultClient
    
)

type Crawler struct {
	// CrawlDelay      time.Duration
	// CrawlDuration   time.Duration
	 CrawlPoliteness bool
	 LogFormat       string
	 LogLevel        string
	 Logger          *logrus.Logger
	 UserAgent       string
	 HTTPClient      fetchbot.Doer
	 URLPool           cache.Cache

	 mux *fetchbot.Mux
	 f   *fetchbot.Fetcher
	 q   *fetchbot.Queue
}

func New() *Crawler {
    
    return &Crawler{
        URLPool: memory.New(),
        mux:     fetchbot.NewMux(),
        LogFormat:       DefaultLogFormat,
		LogLevel:        DefaultLogLevel,
    }
}

func (c *Crawler) Start() {
    if c.Logger == nil {
		c.Logger = newLogger(c.LogFormat, c.LogLevel)
	}

	c.mux.HandleErrors(c.newErrorHandler())
    
    f := fetchbot.New(h)
	//f.CrawlDelay = c.CrawlDelay
	f.DisablePoliteness = !c.CrawlPoliteness
	f.HttpClient = c.HTTPClient
	f.UserAgent = c.UserAgent

	c.f = f
	c.q = c.f.Start()
    
}

func (c *Crawler) Block() {
	c.q.Block()
}

func (c *Crawler) Finish() {
	c.q.Close()
}


func (c *Crawler) newErrorHandler() fetchbot.Handler {
	return fetchbot.HandlerFunc(func(ctx *fetchbot.Context, res *http.Response, err error) {
		// c.Logger.WithFields(logrus.Fields{
		// 	"url":    ctx.Cmd.URL(),
		// 	"method": ctx.Cmd.Method(),
		// }).Error(err)
	})
}

func newLogger(format, level string) *logrus.Logger {
	log := logrus.New()
	log.Out = os.Stdout
	log.Formatter = newFormatter(format)
	log.Level = parseLogLevel(level)
	return log
}

func newFormatter(format string) logrus.Formatter {
	switch format {
	case "text", "":
		return &logrus.TextFormatter{}
	case "json":
		return &logrus.JSONFormatter{}
	default:
		return &logrus.TextFormatter{}
	}
}

func parseLogLevel(level string) logrus.Level {
	switch level {
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	case "error":
		return logrus.ErrorLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	default:
		return logrus.InfoLevel
	}
}