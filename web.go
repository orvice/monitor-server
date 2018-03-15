package main

import (
	"github.com/DeanThompson/ginpprof"
	"github.com/catpie/cors"
	"github.com/gin-gonic/gin"
)

func web() {
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/g", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	if Debug {
		ginpprof.Wrapper(r)
	}

	r.GET("/socket.io/", SocketIOGinHandler)
	r.POST("/socket.io/", SocketIOGinHandler)
	r.Handle("WS", "/socket.io/", SocketIOGinHandler)
	r.Handle("WSS", "/socket.io/", SocketIOGinHandler)
	r.Run(ListenAddr) // listen and serve on 0.0.0.0:8080
}
