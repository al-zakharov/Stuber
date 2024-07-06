package yaml

import (
	"gopkg.in/yaml.v3"
	"os"
	"stuber/internal/router"
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
	f, err := os.ReadFile(stubFilePath)
	if err != nil {
		return nil, err
	}

	var c StubCollection
	if err := yaml.Unmarshal(f, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *StubCollection) MapToStubs() []*router.Stub {
	s := make([]*router.Stub, 0)
	for _, i := range c.Items {
		s = append(s, &router.Stub{
			HttpMethod: i.HttpMethod,
			Path:       i.Path,
			Body:       i.Body,
			Status:     i.Status,
		})
	}

	return s
}
