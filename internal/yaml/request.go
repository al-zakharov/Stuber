package yaml

import (
	"gopkg.in/yaml.v3"
	"os"
)

type StubCollection struct {
	Items map[string]*Stub `yaml:"collection"`
}

type Stub struct {
	HttpMethod string `yaml:"http_method"`
	Path       string `yaml:"path"`
	Body       string `yaml:"body"`
	BodyPath   string `yaml:"body_path"`
	Status     int    `yaml:"status"`
}

func NewStubCollection(stubFilePath string) (*StubCollection, error) {
	var sc StubCollection

	f, err := os.ReadFile(stubFilePath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(f, &sc); err != nil {
		return nil, err
	}

	return &sc, nil
}
