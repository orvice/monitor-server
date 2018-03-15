package mod

import (
	"github.com/googollee/go-socket.io"
	cm "github.com/orvice/monitor-client/mod"
)

type Packet struct {
	NodeID  int32
	Message []byte
}

type NodeStat struct {
	NodeID int32 `json:"node_id"`
	Stat   cm.SystemInfo
}

type IoClient struct {
	So socketio.Socket
}

func (i *IoClient) ID() string {
	return i.So.Id()
}

func NewIoClient(so socketio.Socket) *IoClient {
	c := &IoClient{
		So: so,
	}
	return c
}

type Node struct {
	ID   int32  `json:"id"`
	Addr string `json:"addr"`
}

type NodeLoader interface {
	GetNodes() ([]Node, error)
}
