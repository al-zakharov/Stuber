package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"stuber/internal/server/stub"
	"sync"
)

func MakeHistoryHandler(h *[]*RequestRecord, next http.Handler, m *sync.Mutex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ub := unmarshalRequestBody(r)

		m.Lock()
		*h = append(*h, NewRequestRecord(r.Method, r.URL.String(), &ub))
		m.Unlock()

		next.ServeHTTP(w, r)
	}
}

func MakeCollectorHandler(c map[string][]*RequestRecord, s *stub.Stub, next http.Handler, m *sync.Mutex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.CollectParams != nil {
			b, ub := unmarshalRequestBody(r)
			key, err := getCollectKeyFromRequest(s, r, b)
			if err == nil {
				m.Lock()
				if _, ok := c[key]; !ok {
					c[key] = make([]*RequestRecord, 0)
				}

				c[key] = append(c[key], &RequestRecord{
					HTTPMethod: r.Method,
					URL:        r.URL.Path,
					Body:       &ub,
				})
				m.Unlock()
			}
		}

		next.ServeHTTP(w, r)
	}
}

func MakeAllRequestsHandler(h *[]*RequestRecord) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(h); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	}
}

func MakeLastRequestHandler(h *[]*RequestRecord) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if len(*h) > 0 {
			if err := json.NewEncoder(w).Encode((*h)[len(*h)-1]); err != nil {
				log.Printf("Error encoding response: %v", err)
			}
		}
	}
}

func MakeGetCollectionHandler(c map[string][]*RequestRecord) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srp := r.URL.Query().Get(SearchRequestParam)
		w.WriteHeader(http.StatusOK)
		if e, ok := c[srp]; ok {
			if err := json.NewEncoder(w).Encode(e); err != nil {
				log.Printf("Error encoding response: %v", err)
			}
		}
	}
}

func unmarshalRequestBody(r *http.Request) ([]byte, json.RawMessage) {
	b, err := readRequestBody(r)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		return nil, nil
	}

	var ub json.RawMessage
	if len(b) > 0 {
		ub, err = unmarshalBody(b)
		if err != nil {
			log.Printf("Error unmarshaling request body: %v", err)
			return nil, nil
		}
	}

	return b, ub
}

func readRequestBody(r *http.Request) ([]byte, error) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	r.Body = io.NopCloser(bytes.NewBuffer(b))
	return b, nil
}

func getCollectKeyFromRequest(s *stub.Stub, r *http.Request, body []byte) (string, error) {
	switch s.CollectParams.Type {
	case stub.CollectTypeJsonPath:
		return gjson.GetBytes(body, s.CollectParams.Value).String(), nil
	case stub.CollectTypeQueryParam:
		return r.URL.Query().Get(s.CollectParams.Value), nil
	case stub.CollectTypePathParam:
		return extractPathParam(s.Path, r.URL.String(), s.CollectParams.Value), nil
	default:
		return "", fmt.Errorf("unsupported collect type")
	}
}
