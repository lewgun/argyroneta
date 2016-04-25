package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"time"

	//config
	"github.com/lewgun/argyroneta/cmd/websrvd/pkg/config"

	//controller
	"github.com/lewgun/argyroneta/cmd/websrvd/pkg/controller"

	//logger
	"github.com/lewgun/argyroneta/pkg/logger"

	//webserver
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/lewgun/argyroneta/pkg/constants"
	"github.com/lewgun/argyroneta/pkg/store/bolt"
	"github.com/lewgun/argyroneta/pkg/store/mysql"
)

var (
	confPath = flag.String("conf", "./websrvd.json", "the path to the config file")
)

func main() {

	flag.Parse()

	runtime.GOMAXPROCS(runtime.NumCPU())
	c := mustPrepare(*confPath)
	powerOn(c)
	powerOff()
}

func setupRouter(r *gin.Engine) {

	r.Use(static.Serve("/", static.LocalFile("web", false)))

	r.StaticFile("/", "web/index.html")

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	controller.SetupRouters(r)

}

func powerOn(c *config.Config) {

	log := logger.New(c.Format, c.Level)

	//bolt
	err := bolt.SharedInstInit(c.Store.BoltConf, log)
	if err != nil {
		panic("boot up bolt failed")
	}

	//mysql
	err = mysql.SharedInstInit(c.Store.MySQLConf, log)
	if err != nil {
		panic("boot up mysql failed")
	}

	if c.RunMode == constants.ModeRelease {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(ginrus.Ginrus(log, time.RFC3339, false))

	r.Use(cors.Default())

	setupRouter(r)

	fmt.Println("websrvd is running at: ", c.HTTPPort)

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", c.HTTPPort),
		Handler:        r,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	srv.ListenAndServe()
}

func powerOff() {
}

func mustPrepare(path string) *config.Config {

	c := config.New()
	err := c.Init(path)
	if err != nil {
		panic(err)
	}

	return c
}
