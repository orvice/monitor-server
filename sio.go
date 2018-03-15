package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"github.com/orvice/monitor-server/mod"
)

var (
	SocketioServer *socketio.Server
)

func InitSocketIoSever() error {
	var err error
	SocketioServer, err = NewSocketIOServer()
	return err
}

func NewSocketIOServer() (*socketio.Server, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}
	server.On("connection", func(so socketio.Socket) {
		soc := mod.NewIoClient(so)

		ioManager.AddIOClient(soc)

		logger.Info("on connection")
		so.Join("chat")
		so.On("chat message", func(msg string) {
			logger.Info("emit:", so.Emit("chat message", msg))
			so.BroadcastTo("chat", "chat message", msg)
		})
		so.On("disconnection", func() {
			logger.Info("on disconnect")
			ioManager.RemoveIOClient(soc)
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		logger.Error("error:", err)
	})

	return server, nil
}

func SocketIOHandler(w http.ResponseWriter, r *http.Request) {
	SocketioServer.ServeHTTP(w, r)
}

func SocketIOGinHandler(c *gin.Context) {
	SocketioServer.ServeHTTP(c.Writer, c.Request)
}
