package request_collector

import "encoding/json"

type RequestRecord struct {
	HTTPMethod string          `json:"httpMethod"`
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

func (r *RequestRecord) MapIncomeRequestBody(requestBody []byte) error {
	var jsonBody json.RawMessage
	if err := json.Unmarshal(requestBody, &jsonBody); err != nil {
		return err
	}

	r.Body = jsonBody
	return nil
}
