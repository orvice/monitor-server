package main

import (
	"net/http"

	"github.com/DeanThompson/ginpprof"
	"github.com/catpie/cors"
	"github.com/gin-gonic/gin"
	"github.com/orvice/monitor-server/mod"
)

var (
	nodes []mod.Node
)

func web() {
	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/nodes", func(c *gin.Context) {
		var err error
		if nodes != nil {
			c.JSON(http.StatusOK, nodes)
			return
		}
		nodes, err = nodeLoader.GetNodes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		c.JSON(http.StatusOK, nodes)
		return
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
