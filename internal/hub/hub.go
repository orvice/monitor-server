package hub

import (
	"github.com/orvice/monitor-server/internal/mod"
	"github.com/orvice/monitor-server/internal/node"
	"github.com/weeon/contract"
)

var (
	Logger    contract.Logger
	Manager    *node.Manager
	NodeLoader mod.NodeLoader
)
