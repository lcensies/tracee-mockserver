package routes

import (
	"github.com/gin-gonic/gin"
)

func RequestCancelRecover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Request.Context().Done()
			}
		}()
		c.Next()
	}
}

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(RequestCancelRecover(), gin.Recovery())

	router.POST("/", HandleEventsSink)
	router.GET("/", HandleEventsCount)
	router.GET("/fileevents", HandleFileEventsCount)
	router.POST("/reset", HandleEventsCountReset)

	return router
}
