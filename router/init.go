package router

import (
	"tinyurl/global"
	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine

// InitRouter add route in this function
func InitRouter() {
	gin.SetMode(gin.ReleaseMode)

	Engine = gin.Default()

	//Engine.Use(global.AuthMiddleware())
	//Engine.Use(global.Recovery(RecoveryHandler))
	//Wrapper(Engine)

	v1 := Engine.Group("/v1")
	{
        v1.GET("/:index", ResolveHandler)
		v1.POST("/short", ShortenHandler)
	}
}
