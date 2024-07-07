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
	CollectParams *CollectParams `yaml:"collect_params"`
}

type CollectParams struct {
	JsonPath   string `yaml:"json_path"`
	QueryParam string `yaml:"query_param"`
}

func NewStubCollection(stubFilePath string) (*StubCollection, error) {
	f, err := os.ReadFile(stubFilePath)
	if err != nil {
		return nil, err
	}

	var c StubCollection
	if err = yaml.Unmarshal(f, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *StubCollection) MapToStubs() []*stub.Stub {
	s := make([]*stub.Stub, 0)
	for _, i := range c.Items {
		var cp *stub.CollectParams
		if i.CollectParams != nil {
			if i.CollectParams.QueryParam != "" {
				cp = &stub.CollectParams{
					Type:  stub.CollectTypeQueryParam,
					Value: i.CollectParams.QueryParam,
				}
			} else if i.CollectParams.JsonPath != "" {
				cp = &stub.CollectParams{
					Type:  stub.CollectTypeJsonPath,
					Value: i.CollectParams.JsonPath,
				}
			}
		}

		body := ""
		if i.Body != "" {
			body = i.Body
		} else if i.BodyPath != "" {
			fc, err := os.ReadFile(i.BodyPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			body = string(fc)
		}

		s = append(s, &stub.Stub{
			HttpMethod:    i.HttpMethod,
			Path:          i.Path,
			Body:          body,
			Status:        i.Status,
			CollectParams: cp,
		})
	}

	return s
}
