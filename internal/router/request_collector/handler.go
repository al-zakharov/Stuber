package request_collector

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"stuber/internal/router/stub"
)

func MakeHistoryHandler(h *[]*RequestRecord, n http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO handle error
		b, _ := io.ReadAll(r.Body)
		ub, _ := unmarshalBody(b)

		*h = append(*h, NewRequestRecord(r.Method, r.URL.String(), ub))

		n.ServeHTTP(w, r)
	}
}

func MakeCollectorHandler(c map[string][]*RequestRecord, cp *stub.CollectParams, n http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cp != nil {
			b, _ := io.ReadAll(r.Body)
			ub, _ := unmarshalBody(b)
			key := ""

			switch cp.Type {
			case stub.CollectTypeJsonPath:
				key = gjson.GetBytes(b, cp.Value).String()
			case stub.CollectTypeQueryParam:
				key = r.URL.Query().Get(cp.Value)
			}

			if _, ok := c[key]; !ok {
				c[key] = make([]*RequestRecord, 0)
			}

			c[key] = append(c[key], &RequestRecord{
				HTTPMethod: r.Method,
				URL:        r.URL.String(),
				Body:       ub,
			})
		}

		n.ServeHTTP(w, r)
	}
}

func MakeAllRequestsHandler(h *[]*RequestRecord) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(h)
	}
}

func MakeLastRequestHandler(h *[]*RequestRecord) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if len(*h) > 0 {
			json.NewEncoder(w).Encode((*h)[len(*h)-1])
		}
	}
}

func MakeGetCollectionHandler(c map[string][]*RequestRecord) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srp := r.URL.Query().Get(SearchRequestParam)
		w.WriteHeader(http.StatusOK)
		if e, ok := c[srp]; ok {
			json.NewEncoder(w).Encode(e)
		}
	}
}
