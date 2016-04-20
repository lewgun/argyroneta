//Package cmdmgr 
package cmdmgr

import (
"sync"
"net/url"

    "github.com/lewgun/argyroneta/pkg/types"
    
	"github.com/PuerkitoBio/fetchbot"
	
)



const (
	NetEase types.Site = "netease"
)


type Command interface {
	fetchbot.Command
    Depth() uint32
}


type CMDMaker func (*types.Rule)

var (
	makers   = map[string]CMDMaker{}
	mu sync.Mutex
)


//Register 注册一个存储设备
func Register(domain string, maker CMDMaker) {
	mu.Lock()
	defer mu.Unlock()
	makers[domain] = maker
}

func Maker(domain string) CMDMaker {
	if maker, ok := makers[domain]; !ok {
		return nil
        
	} else {
		return maker
	}

}

