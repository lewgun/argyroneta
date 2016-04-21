package main

import (
	"flag"

	"github.com/lewgun/argyroneta/cmd/argyronetad/pkg/config"

	"github.com/lewgun/argyroneta/pkg/store"
	_ "github.com/lewgun/argyroneta/pkg/store/bolt"

	"github.com/lewgun/argyroneta/pkg/spidermgr"
	"github.com/lewgun/argyroneta/pkg/logger"
)

var (
	confPath = flag.String("conf", "./argyronetad.json", "the path to the config file")
)

func powerOn(c *config.Config) func() {

	s, err := store.PowerOn(store.Bolt, c.Store.Path)
	if err != nil {
		panic(err)
	}

	logger := logger.New(c.Format, c.Level)

	sm := spidermgr.SharedInstance()
	errs := sm.Init(c.Sites, logger)
	if err != nil {
		panic(errs[0])
	}

	sm.PowerOn()

	return func() {
		s.PowerOff()
		sm.PowerOff()
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

	defer powerOn(mustPrepare(*confPath))()

}
