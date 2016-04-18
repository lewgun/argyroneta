package command

import (
"sync"
	"github.com/lewgun/argyroneta/vendor/github.com/PuerkitoBio/fetchbot"
	"github.com/MieYua/Aliyun-OSS-Go-SDK/oss/types"
)


const (
	NetEase = "netease"
)


type Command interface {
	fetchbot.Command
}

type Cmd struct {
	*fetchbot.Cmd
	S *url.URL
	D int
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
