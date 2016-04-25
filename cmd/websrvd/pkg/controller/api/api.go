package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lewgun/argyroneta/pkg/types"

	"github.com/lewgun/argyroneta/pkg/misc"
	"github.com/lewgun/argyroneta/pkg/store/mysql"
)

func EntryTop(ctx *gin.Context) {
	req := &types.TopReq{}
	err := ctx.BindJSON(req)
	if err != nil {
		misc.FastResponse(ctx, err)
		return
	}

	entries, err := mysql.M.EntryTopN(req)
	if err != nil {
		misc.FastResponse(ctx, err)
		return
	}

	m := gin.H{
		"entries": entries,
	}

	misc.FastResponse(ctx, m)

}
