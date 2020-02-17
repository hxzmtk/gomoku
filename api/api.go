package api

import (
	v1 "github.com/bzyy/gomoku/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// static file api
	LoadStatic(r)

	// load html
	LoadHtml(r)

	v1.LoadV1(r)
	return r
}
