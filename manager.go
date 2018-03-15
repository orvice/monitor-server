package main

import (
	"context"
	"sync"

	"github.com/orvice/kit/log"
	"github.com/orvice/monitor-server/mod"
)

type IoManager struct {
	ctx     context.Context
	clients *sync.Map
	logger  log.Logger
}

func NewIoManager(logger log.Logger) *IoManager {
	iom := &IoManager{
		clients: new(sync.Map),
		logger:  logger,
	}

	return iom
}

func (i *IoManager) AddIOClient(c *mod.IoClient) {
	i.logger.Infof("add client %s", c.ID())
	i.clients.Store(c.ID(), c)
}

func (i *IoManager) RemoveIOClient(c *mod.IoClient) {
	i.logger.Infof("remove client %s", c.ID())
	i.clients.Delete(c.ID())
}

func (i *IoManager) Broadcast(event string, args ...interface{}) {
	i.clients.Range(func(k, v interface{}) bool {
		c, ok := v.(*mod.IoClient)
		if !ok {
			return true
		}
		c.So.Emit(event, args...)
		return true
	})
}
