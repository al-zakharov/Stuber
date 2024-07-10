package collector

import (
	"encoding/json"
	"fmt"
	"regexp"
)

const SearchRequestParam = "searchRequestParam"

type RequestRecord struct {
	HTTPMethod string           `json:"http_method"`
	URL        string           `json:"url"`
	Body       *json.RawMessage `json:"body"`
}

func NewRequestRecord(HTTPMethod string, URL string, Body *json.RawMessage) *RequestRecord {
	return &RequestRecord{
		HTTPMethod: HTTPMethod,
		URL:        URL,
		Body:       Body,
	}
}

func unmarshalBody(b []byte) (json.RawMessage, error) {
	var j json.RawMessage
	if err := json.Unmarshal(b, &j); err != nil {
		return nil, err
	}

	return j, nil
}

func extractPathParam(pattern, path, param string) string {
	rePattern := regexp.MustCompile(`:([^/]+)`).ReplaceAllStringFunc(pattern, func(m string) string {
		return fmt.Sprintf("(?P<%s>[^/]+)", m[1:])
	})

	re := regexp.MustCompile(rePattern)

	match := re.FindStringSubmatch(path)
	if match == nil {
		return ""
	}

	ns := re.SubexpNames()
	for i, n := range ns {
		if n == param {
			return match[i]
		}
	}

	return ""
}
