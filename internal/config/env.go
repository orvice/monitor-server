package config

import "github.com/orvice/utils/env"

var (
	ListenAddr     string
	NodeConfigPath string
	Debug          bool

	LoaderMode string
	SkipStream bool

	MuUrl, MuToken string

	MonitorMethod string

	MqUrl, MqQueue string

	PostKey string
)

const (
	MuMode = "mu"

	Grpc = "grpc"
)

var (
	DefaultNodeConfigPath = "/etc/monitor-server/config.json"
)

func InitEnv() {
	LoaderMode = env.Get("LOADER_MODE")
	ListenAddr = env.Get("LISTEN_ADDR", ":80")
	MonitorMethod = env.Get("MONITOR_METHOD", "ws")
	NodeConfigPath = env.Get("NODE_CONFIG_PATH", "config.json")
	if env.Get("DEBUG") == "true" {
		Debug = true
	}

	if env.Get("SKIP_STREAM") == "true" {
		SkipStream = true
	}

	MuUrl = env.Get("MU_URL")
	MuToken = env.Get("MU_TOKEN")

	MqQueue = env.Get("MQ_QUEUE")
	MqUrl = env.Get("MQ_URL")

	PostKey = env.Get("POST_KEY")
}
