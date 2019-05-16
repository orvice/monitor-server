package web

import (
	"fmt"
	"github.com/orvice/monitor-server/internal/sio"
	"github.com/weeon/contract"
	"github.com/weeon/utils/ginutil"
	"net/http"

	"github.com/DeanThompson/ginpprof"
	"github.com/catpie/cors"
	"github.com/gin-gonic/gin"
	"github.com/orvice/monitor-server/internal/config"
	"github.com/orvice/monitor-server/internal/hub"
	"github.com/orvice/monitor-server/internal/mod"
)

var (
	nodes []mod.Node
)

func getNodesMap() map[string]mod.Node {
	m := make(map[string]mod.Node)
	for _, n := range nodes {
		m[fmt.Sprintf("%d", n.ID)] = n
	}
	return m
}

func InitWeb(l contract.Logger) {
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(ginutil.WrapRequestID)
	r.GET("/nodes", func(c *gin.Context) {
		var err error
		if nodes != nil {
			c.JSON(http.StatusOK, nodes)
			return
		}
		nodes, err = hub.NodeLoader.GetNodes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			return
		}
		c.JSON(http.StatusOK, nodes)
		return
	})

	if config.Debug {
		ginpprof.Wrapper(r)
	}

	r.GET("/socket.io/", sio.SocketIOGinHandler)
	r.POST("/socket.io/", sio.SocketIOGinHandler)
	r.Handle("WS", "/socket.io/", sio.SocketIOGinHandler)
	r.Handle("WSS", "/socket.io/", sio.SocketIOGinHandler)
	r.POST("/nodes/:node_id/data", DataCollect)
	r.Run(config.ListenAddr) // listen and serve on 0.0.0.0:8080
}
