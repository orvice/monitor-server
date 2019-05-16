package sio

import (
	"github.com/weeon/contract"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"github.com/orvice/monitor-server/internal/mod"
)

var (
	SocketioServer *socketio.Server
	IOM  *IoManager
)

func InitSocketIoSever(l contract.Logger) error {
	var err error
	SocketioServer, err = NewSocketIOServer(l)
	return err
}

func NewSocketIOServer(l contract.Logger) (*socketio.Server, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}
	server.On("connection", func(so socketio.Socket) {
		soc := mod.NewIoClient(so)

		IOM.AddIOClient(soc)

		l.Info("on connection")
		so.Join("chat")
		so.On("chat message", func(msg string) {
			l.Info("emit:", so.Emit("chat message", msg))
			so.BroadcastTo("chat", "chat message", msg)
		})
		so.On("disconnection", func() {
			l.Info("on disconnect")
			IOM.RemoveIOClient(soc)
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		l.Error("error:", err)
	})

	return server, nil
}

func SocketIOHandler(w http.ResponseWriter, r *http.Request) {
	SocketioServer.ServeHTTP(w, r)
}

func SocketIOGinHandler(c *gin.Context) {
	SocketioServer.ServeHTTP(c.Writer, c.Request)
}
