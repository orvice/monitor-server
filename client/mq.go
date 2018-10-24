package client

import (
	"encoding/json"
	"github.com/orvice/kit/log"
	"github.com/orvice/kit/utils/mq"
	cm "github.com/orvice/monitor-client/mod"
	"github.com/orvice/monitor-server/mod"
)

type MqClient struct {
	messageRecv chan mod.Packet
	logger      log.Logger
	consumer    *mq.Consumer
	url, queue  string
}

func NewMqClient(url, queue string, ch chan mod.Packet, log log.Logger) *MqClient {
	return &MqClient{
		messageRecv: ch,
		logger:      log,
		url:         url,
		queue:       queue,
	}
}

func (m *MqClient) Init() {
	var err error
	m.consumer, err = mq.NewConsumer(m.url, m.queue, "", m.handle)
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

	m.messageRecv <- mod.Packet{
		NodeID:  ns.NodeID,
		Message: bs,
	}

}