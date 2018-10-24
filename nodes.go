package main

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/orvice/kit/log"
	cm "github.com/orvice/monitor-client/mod"
	"github.com/orvice/monitor-server/client"
	"github.com/orvice/monitor-server/enum"
	"github.com/orvice/monitor-server/mod"
)

type Manager struct {
	packetCh chan mod.Packet
	ctx      context.Context
	cancel   func()
	clients  *sync.Map

	nodeLoader mod.NodeLoader
	logger     log.Logger
}

func NewManager(nl mod.NodeLoader, l log.Logger) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	m := &Manager{
		packetCh:   make(chan mod.Packet, 10),
		ctx:        ctx,
		cancel:     cancel,
		clients:    new(sync.Map),
		nodeLoader: nl,
		logger:     l,
	}
	return m
}

func (m *Manager) Run() error {
	nodes, err := m.nodeLoader.GetNodes()
	if err != nil {
		return err
	}
	for _, n := range nodes {

		switch MonitorMethod {
		case Grpc:
			c := client.NewGrpcClient(n.ID, n.GrpcAddr, m.packetCh, m.logger)
			m.logger.Infof("run node %d addr: %s in grpc mod", n.ID, n.GrpcAddr)
			go c.Run()
			m.clients.Store(n.ID, c)
		default:
			c := client.NewClient(n.ID, n.Addr, m.packetCh, m.logger)
			m.logger.Infof("run node %d addr: %s", n.ID, n.Addr)
			go c.Run()
			m.clients.Store(n.ID, c)
		}

		if len(MqUrl) != 0 {
			mqc := client.NewMqClient(MqUrl, MqQueue, m.packetCh, m.logger)
			go mqc.Init()
		}

	}

	go m.packetHandle()
	return nil
}

func (m *Manager) packetHandle() {
	for {
		select {
		case p := <-m.packetCh:
			var stat cm.SystemInfo
			err := json.Unmarshal(p.Message, &stat)
			if err != nil {
				m.logger.Error(err)
				continue
			}
			ioManager.Broadcast(enum.EventServerStat, mod.NodeStat{
				NodeID: p.NodeID,
				Stat:   stat,
			})
		case <-m.ctx.Done():
			return
		}
	}
}
