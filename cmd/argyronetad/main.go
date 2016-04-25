package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/lewgun/argyroneta/cmd/argyronetad/pkg/config"

	"github.com/lewgun/argyroneta/pkg/store/bolt"
	"github.com/lewgun/argyroneta/pkg/store/mysql"

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

	//bolt
	err := bolt.SharedInstInit(c.Store.BoltConf, logger)
	if err != nil {
		logger.Fatalln("boot up bolt failed with error: %v", err)
	}

	//mysql
	err = mysql.SharedInstInit(c.Store.MySQLConf, logger)
	if err != nil {
		logger.Fatalln("boot up mysql failed with invalid parameter")
	}

	rules, err := mysql.M.Rules()
	if err != nil {
		logger.Fatalf("can't get rules with error: %v", err)
	}

	//spider manager
	err = spidermgr.SharedInstInit(rules, logger)
	if err != nil {
		logger.Fatalf("boot up spider manager failed with error: %v", err)
	}

	spidermgr.SM.PowerOn()

	return func(sig <-chan os.Signal) {
		<-sig

		spidermgr.SM.PowerOff()

		bolt.B.PowerOff()
		mysql.M.PowerOff()

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
