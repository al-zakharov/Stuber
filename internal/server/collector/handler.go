package collector

import (
	"bytes"
	"encoding/json"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"stuber/internal/server/stub"
)

func MakeHistoryHandler(h *[]*RequestRecord, n http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var ub json.RawMessage
		if len(b) > 0 {
			ub, err = unmarshalBody(b)
			if err != nil {
				log.Printf("Error unmarshaling request body: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		} else {
			ub = nil
		}

		r.Body = io.NopCloser(bytes.NewBuffer(b))
		*h = append(*h, NewRequestRecord(r.Method, r.URL.String(), &ub))

		n.ServeHTTP(w, r)
	}
}

func MakeCollectorHandler(c map[string][]*RequestRecord, s *stub.Stub, n http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.CollectParams != nil {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading request body: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			var ub json.RawMessage
			if len(b) > 0 {
				ub, err = unmarshalBody(b)
				if err != nil {
					log.Printf("Error unmarshaling request body: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			} else {
				ub = nil
			}

			key := ""

			switch s.CollectParams.Type {
			case stub.CollectTypeJsonPath:
				key = gjson.GetBytes(b, s.CollectParams.Value).String()
			case stub.CollectTypeQueryParam:
				key = r.URL.Query().Get(s.CollectParams.Value)
			case stub.CollectTypePathParam:
				key = extractPathParam(s.Path, r.URL.String(), s.CollectParams.Value)
			}

			if _, ok := c[key]; !ok {
				c[key] = make([]*RequestRecord, 0)
			}

			c[key] = append(c[key], &RequestRecord{
				HTTPMethod: r.Method,
				URL:        r.URL.Path,
				Body:       &ub,
			})
		}

		n.ServeHTTP(w, r)
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
