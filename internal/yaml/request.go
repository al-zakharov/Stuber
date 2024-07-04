package yaml

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Stub struct {
	HttpMethod string `yaml:"http_method"`
	Path       string `yaml:"path"`
	Body       string `yaml:"body"`
	BodyPath   string `yaml:"body_path"`
	Status     int    `yaml:"status"`
}

func NewStub(stubFilePath string) (*Stub, error) {
	var s *Stub

	c, err := os.ReadFile(stubFilePath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(c, &s); err != nil {
		return nil, err
	}

	return s, nil
}
