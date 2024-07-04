package yaml

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Stub struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	HttpMethod  string `yaml:"http_method"`
	Endpoint    string `yaml:"endpoint"`
	Body        string `yaml:"body"`
	BodyPath    string `yaml:"body_path"`
	Status      string `yaml:"status"`
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
