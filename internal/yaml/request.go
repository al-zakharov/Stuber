package yaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"stuber/internal/router/stub"
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

	s := make([]*stub.Stub, 0, len(c.Items))
	for _, i := range c.Items {
		var cp *stub.CollectParams
		if i.CollectParams != nil {
			cp = i.mapStubCollectParam()
		}

		body, err := i.getBodyContent()
		if err != nil {
			return nil, err
		}

		s = append(s, &stub.Stub{
			HttpMethod:    i.HttpMethod,
			Path:          i.Path,
			Body:          body,
			Status:        i.Status,
			DynamicBody:   i.DynamicBody,
			CollectParams: cp,
		})
	}

	return s, nil
}

func (s *Stub) mapStubCollectParam() *stub.CollectParams {
	var scp stub.CollectParams
	if s.CollectParams.QueryParam != "" {
		scp = stub.CollectParams{
			Type:  stub.CollectTypeQueryParam,
			Value: s.CollectParams.QueryParam,
		}
	} else if s.CollectParams.JsonPath != "" {
		scp = stub.CollectParams{
			Type:  stub.CollectTypeJsonPath,
			Value: s.CollectParams.JsonPath,
		}
	} else if s.CollectParams.PathParam != "" {
		scp = stub.CollectParams{
			Type:  stub.CollectTypePathParam,
			Value: s.CollectParams.PathParam,
		}
	}

	return &scp
}

func (s *Stub) getBodyContent() (string, error) {
	body := ""
	if s.Body != "" {
		body = s.Body
	} else if s.BodyPath != "" {
		fc, err := os.ReadFile(s.BodyPath)
		if err != nil {
			return "", fmt.Errorf("failed to read body file: %w", err)
		}
		body = string(fc)
	}

	return body, nil
}
