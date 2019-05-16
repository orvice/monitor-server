package node

import (
	"encoding/json"
	"github.com/catpie/musdk-go"
	"github.com/orvice/monitor-server/internal/config"
	"github.com/orvice/monitor-server/internal/mod"
	"io/ioutil"
)

type FileNodeLoader struct {
	path string
}

func NewFileNodeLoader(path string) *FileNodeLoader {
	return &FileNodeLoader{
		path: path,
	}
}

func (f *FileNodeLoader) GetNodes() ([]mod.Node, error) {
	s, err := ioutil.ReadFile(f.path)
	if err != nil {
		return nil, err
	}
	var nodes []mod.Node
	err = json.Unmarshal(s, &nodes)
	return nodes, err
}

type MuLoader struct {
	client *musdk.Client
}

func NewMuLoader(muUri, token string) *MuLoader {
	cli := musdk.NewClient(muUri, token, 0, 0)
	return &MuLoader{
		client: cli,
	}
}

func (m *MuLoader) GetNodes() ([]mod.Node, error) {
	nodes, err := m.client.GetNodes()
	if err != nil {
		return nil, err
	}
	out := make([]mod.Node, 0)
	for _, v := range nodes {
		if len(v.ServerMonitorAddr) == 0 {
			continue
		}
		out = append(out, mod.Node{
			ID:       v.ID,
			Name:     v.Name,
			Addr:     v.GetMonitorWsAddr(),
			GrpcAddr: v.GetMonitorGrpcAddr(),
		})
	}
	return out, nil
}

func InitNodeLoader() mod.NodeLoader {
	if config.LoaderMode == config.MuMode {
		return  NewMuLoader(config.MuUrl, config.MuToken)
	}
	return NewFileNodeLoader(config.NodeConfigPath)
}
