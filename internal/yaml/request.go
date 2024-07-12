package yaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"stuber/internal/server/stub"
)

type StubCollection struct {
	Items map[string]*Stub `yaml:"collection"`
}

type Stub struct {
	HttpMethod    string         `yaml:"http_method"`
	Path          string         `yaml:"path"`
	Body          string         `yaml:"body"`
	BodyPath      string         `yaml:"body_path"`
	Status        int            `yaml:"status"`
	DynamicBody   bool           `yaml:"dynamic_body"`
	CollectParams *CollectParams `yaml:"collect_params"`
}

type CollectParams struct {
	JsonPath   string `yaml:"json_path"`
	QueryParam string `yaml:"query_param"`
	PathParam  string `yaml:"path_param"`
}

func NewStubCollection(stubFilePath string) (*StubCollection, error) {
	f, err := os.ReadFile(stubFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read stub file: %w", err)
	}

	var c StubCollection
	if err = yaml.Unmarshal(f, &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stub file: %w", err)
	}

	return &c, nil
}

func (c *StubCollection) MapToStubs() ([]*stub.Stub, error) {
	if len(c.Items) == 0 {
		return nil, nil
	}

	var stubs []*stub.Stub
	for _, i := range c.Items {
		cp := i.mapStubCollectParam()

		body, err := i.getBodyContent()
		if err != nil {
			return nil, err
		}

		stubs = append(stubs, &stub.Stub{
			HttpMethod:    i.HttpMethod,
			Path:          i.Path,
			Body:          body,
			Status:        i.Status,
			DynamicBody:   i.DynamicBody,
			CollectParams: cp,
		})
	}

	return stubs, nil
}

func (s *Stub) mapStubCollectParam() *stub.CollectParams {
	if s.CollectParams == nil {
		return nil
	}

	switch {
	case s.CollectParams.QueryParam != "":
		return &stub.CollectParams{
			Type:  stub.CollectTypeQueryParam,
			Value: s.CollectParams.QueryParam,
		}
	case s.CollectParams.PathParam != "":
		return &stub.CollectParams{
			Type:  stub.CollectTypePathParam,
			Value: s.CollectParams.PathParam,
		}
	case s.CollectParams.JsonPath != "":
		return &stub.CollectParams{
			Type:  stub.CollectTypeJsonPath,
			Value: s.CollectParams.JsonPath,
		}
	default:
		return nil
	}
}

func (s *Stub) getBodyContent() (string, error) {
	if s.Body != "" {
		return s.Body, nil
	}

	if s.BodyPath != "" {
		fc, err := os.ReadFile(s.BodyPath)
		if err != nil {
			return "", fmt.Errorf("failed to read body file: %w", err)
		}
		return string(fc), nil
	}

	return "", nil
}
