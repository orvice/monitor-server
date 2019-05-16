package main

import (
	"github.com/orvice/monitor-server/internal/sio"
	"os"

	"github.com/orvice/monitor-server/internal/config"
	"github.com/orvice/monitor-server/internal/hub"
	"github.com/orvice/monitor-server/internal/node"
	"github.com/orvice/monitor-server/internal/web"
	"github.com/weeon/log"
	"go.uber.org/zap/zapcore"
)

func main() {
	var err error
	config.InitEnv()
	hub.Logger, _ = log.NewLogger("/app/log/app.log", zapcore.DebugLevel)
	sio.IOM = sio.NewIoManager(hub.Logger)
	hub.NodeLoader = node.InitNodeLoader()
	hub.Manager = node.NewManager(hub.NodeLoader, hub.Logger)
	err = sio.InitSocketIoSever(hub.Logger)
	errExit(err)
	err = hub.Manager.Run()
	errExit(err)
	web.InitWeb()
}

func errExit(err error) {
	if err != nil {
		hub.Logger.Error(err)
		os.Exit(0)
	}
}
