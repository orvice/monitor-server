package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/orvice/monitor-server/mod"
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
