package main

import (
	"os"

	"github.com/orvice/kit/log"
	"github.com/orvice/monitor-server/mod"
)

var (
	logger     log.Logger
	ioManager  *IoManager
	manager    *Manager
	nodeLoader mod.NodeLoader
)

func main() {
	var err error
	InitEnv()
	logger = log.NewDefaultLogger()
	ioManager = NewIoManager(logger)
	nodeLoader = NewFileNodeLoader(NodeConfigPath)
	manager = NewManager(nodeLoader, logger)
	err = InitSocketIoSever()
	errExit(err)
	err = manager.Run()
	errExit(err)
	web()
}

func errExit(err error) {
	if err != nil {
		logger.Error(err)
		os.Exit(0)
	}
}
