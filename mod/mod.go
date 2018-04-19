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
	Name string `json:"name"`
	Addr string `json:"addr"`

	// Stat
	BytesTotal        int64   `json:"bytesTotal"`
	MemoryTotal       int64   `json:"memoryTotal"`
	MemoryUsed        int64   `json:"memoryUsed"`
	MemoryUsedPercent float64 `json:"memoryUsedPercent"`
	DiskUsagePercent  float64 `json:"disk_usage_percent"`
	CpuCount          int     `json:"cpu_count"`
	Load1             float64 `json:"load1"`
	Load5             float64 `json:"load5"`
	Load15            float64 `json:"load15"`
}

type NodeLoader interface {
	GetNodes() ([]Node, error)
}
