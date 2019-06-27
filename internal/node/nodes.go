package node

import (
	"context"
	"encoding/json"
	"github.com/orvice/monitor-server/internal/sio"
	"github.com/weeon/contract"
	"sync"
	cm "github.com/orvice/monitor-client/mod"
	"github.com/orvice/monitor-server/enum"
	"github.com/orvice/monitor-server/internal/client"
	"github.com/orvice/monitor-server/internal/config"
	"github.com/orvice/monitor-server/internal/mod"
	"time"
)

type Manager struct {
	packetCh chan mod.Packet
	ctx      context.Context
	cancel   func()
	clients  *sync.Map

	nodeLoader mod.NodeLoader
	logger     contract.Logger
	lastTime   time.Time
}

func NewManager(nl mod.NodeLoader, l contract.Logger) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	m := &Manager{
		packetCh:   make(chan mod.Packet, 10),
		ctx:        ctx,
		cancel:     cancel,
		clients:    new(sync.Map),
		nodeLoader: nl,
		logger:     l,
		lastTime: time.Now(),
	}
	return m
}

func (m *Manager) Run() error {
	nodes, err := m.nodeLoader.GetNodes()
	if err != nil {
		return err
	}
	for _, n := range nodes {

		if config.SkipStream {
			continue
		}

		switch config.MonitorMethod {
		default:
			c := client.NewGrpcClient(n.ID, n.GrpcAddr, m.packetCh, m.logger)
			m.logger.Infof("run node %d addr: %s in grpc mod", n.ID, n.GrpcAddr)
			go c.Run()
			m.clients.Store(n.ID, c)
		}

	}

	if len(config.MqUrl) != 0 {
		mqc := client.NewMqClient(config.MqUrl, config.MqQueue, m.packetCh, m.logger)
		go mqc.Init(m.logger)
	}

	go m.packetHandle()
	return nil
}

func(m *Manager) GetLastTime() time.Time{
	return m.lastTime
}

func (m *Manager) packetHandle() {
	for {
		select {
		case p := <-m.packetCh:
			var stat cm.SystemInfo
			err := json.Unmarshal(p.Message, &stat)
			m.lastTime = time.Now()
			if err != nil {
				m.logger.Error(err)
				continue
			}
			sio.IOM.Broadcast(enum.EventServerStat, mod.NodeStat{
				NodeID: p.NodeID,
				Stat:   stat,
			})
		case <-m.ctx.Done():
			return
		}
	}
}

func (m *Manager) SendPacket(p mod.Packet){
	m.packetCh <- p
}