package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/lewgun/argyroneta/cmd/argyronetad/pkg/config"

	"github.com/lewgun/argyroneta/pkg/store"
	_ "github.com/lewgun/argyroneta/pkg/store/bolt"

	"github.com/lewgun/argyroneta/pkg/spidermgr"
	_ "github.com/lewgun/argyroneta/pkg/spidermgr/site/netease"

	"github.com/lewgun/argyroneta/pkg/logger"
)

var (
	confPath = flag.String("conf", "./argyronetad.json", "the path to the config file")
)

func powerOn(c *config.Config) func(<-chan os.Signal) {

	//logger
	logger := logger.New(c.Format, c.Level)

	//store
	s, err := store.PowerOn(store.Bolt, c.Store.Path, logger)
	if err != nil {
		panic(err)
	}

	//spider manager
	sm := spidermgr.SharedInstance()
	errs := sm.Init(c.Sites, logger)
	if err != nil {
		panic(errs[0])
	}

	sm.PowerOn()

	return func(sig <-chan os.Signal) {
		<-sig

		sm.PowerOff()
		s.PowerOff()

	}
}

func mustPrepare(path string) *config.Config {

	c := config.New()
	err := c.Init(path)
	if err != nil {
		panic(err)
	}
	return c

}
func main() {

	flag.Parse()
	sig := make(chan os.Signal, 1)

	go signal.Notify(sig, os.Interrupt, os.Kill)

	powerOn(mustPrepare(*confPath))(sig)

}
