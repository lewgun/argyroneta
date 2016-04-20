package site

import (
    "net/url"
    
	//"github.com/lewgun/argyroneta/pkg/command"
	"github.com/lewgun/argyroneta/pkg/rule"

	"github.com/PuerkitoBio/goquery"
)


type NetEaseCommand struct {
     *url.URL
	Method string
    Depth int 
}


func (ne *NetEaseCommand)Handle(ctx *Context, rspn *http.Response, error) {
    
    if err != nil {
        //todo error 
    }
    
    doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
			// c.Logger.WithFields(logrus.Fields{
			// 	"url":    context.URL(),
			// 	"method": context.Method(),
			// }).Error(err)
			// return
	}
    
    ne.depistcher(doc)
        
}

func (ne *NetEaseCommand) depistcher( doc *goquery.Document ) error 


func init() {

	command.Register(command.NetEase, neteaseCMDMaker)
}

func neteaseCMDMaker( ctx rawURL, method string , depth uint32) (command.Command, error )  {
    
    parsed, err := url.Parse(rawURL)
    if err != nil {
        return nil, errutil.ErrInvalidParameter
        
    }
    nec := &NetEaseCommand {
        URL : parsed,
        Method: method,
        Depth: depth,
        
    }
    
    return nec, nil 

}

func (c *Crawler) newHTMLHandler(procs ...Processor) fetchbot.Handler {
	return fetchbot.HandlerFunc(func(ctx *fetchbot.Context, res *http.Response, err error) {
		context := &Ctx{ctx, c.Cache, c.Logger}
		doc, err := goquery.NewDocumentFromResponse(res)
		if err != nil {
			c.Logger.WithFields(logrus.Fields{
				"url":    context.URL(),
				"method": context.Method(),
			}).Error(err)
			return
		}

		for _, p := range procs {
			ok := p.Process(context, doc)
			if !ok {
				return
			}
		}
	})
}
