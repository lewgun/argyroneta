package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lewgun/argyroneta/cmd/websrvd/pkg/controller/api"
)

//SetupRouters setup all controllers
func SetupRouters(r *gin.Engine) {
	r.POST("/api/entry/top", api.EntryTop)
}
