package dynamic

import "encoding/json"

type Body struct {
	Path string          `json:"path"`
	Body json.RawMessage `json:"body"`
}
