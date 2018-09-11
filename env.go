package main

import "github.com/orvice/utils/env"

var (
	ListenAddr     string
	NodeConfigPath string
	Debug          bool

	LoaderMode string

	MuUrl, MuToken string
)

const (
	MuMode = "mu"
)

var (
	DefaultNodeConfigPath = "/etc/monitor-server/config.json"
)

func InitEnv() {
	LoaderMode = env.Get("LOADER_MODE")
	ListenAddr = env.Get("LISTEN_ADDR", ":80")
	NodeConfigPath = env.Get("NODE_CONFIG_PATH", "config.json")
	if env.Get("DEBUG") == "true" {
		Debug = true
	}

	MuUrl = env.Get("MU_URL")
	MuToken = env.Get("MU_TOKEN")
}
