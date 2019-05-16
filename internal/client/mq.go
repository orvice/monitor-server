package client

import (
	"encoding/json"
	cm "github.com/orvice/monitor-client/mod"
	"github.com/orvice/monitor-server/internal/mod"
	"github.com/weeon/contract"
	"github.com/weeon/utils/mq"
)

type MqClient struct {
	messageRecv chan mod.Packet
	logger      contract.Logger
	consumer    *mq.Consumer
	url, queue  string
}

func NewMqClient(url, queue string, ch chan mod.Packet, log contract.Logger) *MqClient {
	return &MqClient{
		messageRecv: ch,
		logger:      log,
		url:         url,
		queue:       queue,
	}
}

func (m *MqClient) Init(l contract.Logger) {
	var err error
	m.consumer, err = mq.NewConsumer(m.url, m.queue, "", m.handle, l, false)
	if err != nil {
		m.logger.Error(err)
		return
	}

}

func (m MqClient) handle(b []byte) {
	var ns cm.NodeStat
	err := json.Unmarshal(b, &ns)
	if err != nil {
		return
	}

	bs, err := json.Marshal(ns.Stat)
	if err != nil {
		return
	}

	m.logger.Debugw("recv message from mq",
		"mod", ns,
	)

	m.messageRecv <- mod.Packet{
		NodeID:  ns.NodeID,
		Message: bs,
	}

}
