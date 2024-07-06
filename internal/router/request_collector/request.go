package request_collector

import "encoding/json"

type RequestRecord struct {
	HTTPMethod string          `json:"http_method"`
	URL        string          `json:"url"`
	Body       json.RawMessage `json:"body"`
}

func NewRequestRecord(HTTPMethod string, URL string, Body json.RawMessage) *RequestRecord {
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
