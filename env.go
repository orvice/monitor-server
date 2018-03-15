package main

import "github.com/orvice/utils/env"

var (
	ListenAddr     string
	NodeConfigPath string
	Debug          bool
)

var (
	DefaultNodeConfigPath = "/etc/monitor-server/config.json"
)

func InitEnv() {
	ListenAddr = env.Get("LISTEN_ADDR", ":80")
	NodeConfigPath = env.Get("NODE_CONFIG_PATH", "config.json")
	if env.Get("DEBUG") == "true" {
		Debug = true
	}
}
