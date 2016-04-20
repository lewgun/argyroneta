package context

import (
	"net/url"


	"github.com/lewgun/argyroneta/pkg/cache"
  
	"github.com/lewgun/argyroneta/pkg/cmdmgr"
      
	"github.com/PuerkitoBio/fetchbot"
    
	"github.com/Sirupsen/logrus"
)

type Context interface {
	URLPool() cache.Cache
    
	//Queue() Enqueuer
	Log(fields map[string]interface{}) *logrus.Entry
	URL() *url.URL
	Method() string
	Depth() uint32
}

type Ctx struct {
	*fetchbot.Context
	urlPool cache.Cache
	L *logrus.Logger
}

func (c *Ctx) URLPool() cache.Cache {
	return c.urlPool
}

// func (c *Ctx) Queue() Enqueuer {
// 	return &Queue{c.Q}
// }

func (c *Ctx) Log(fields map[string]interface{}) *logrus.Entry {
	data := logrus.Fields{}
	for k, v := range fields {
		data[k] = v
	}
	return c.L.WithFields(data)
}

func (c *Ctx) URL() *url.URL {
	return c.Cmd.URL()
}

func (c *Ctx) Method() string {
	return c.Cmd.Method()
}


func (c *Ctx) Depth() uint32 {
    
    myCmd := c.Cmd.(cmdmgr.Command)
    return myCmd.Depth()
}
