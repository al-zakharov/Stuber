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
		rp := r.URL.Query().Get(SaveRequestParam)
		if rp != "" {
			b, _ := io.ReadAll(r.Body)
			ub, _ := unmarshalBody(b)

			if _, ok := c[rp]; !ok {
				c[rp] = make([]*RequestRecord, 0)
			}

			c[rp] = append(c[rp], &RequestRecord{
				HTTPMethod: r.Method,
				URL:        r.URL.String(),
				Body:       ub,
			})
		} else if cp != nil && cp.Value != "" {
			b, _ := io.ReadAll(r.Body)
			ub, _ := unmarshalBody(b)

			k := gjson.GetBytes(b, cp.Value).String()

			if _, ok := c[k]; !ok {
				c[k] = make([]*RequestRecord, 0)
			}

			c[k] = append(c[k], &RequestRecord{
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
		gp := r.URL.Query().Get(SaveRequestParam)
		w.WriteHeader(http.StatusOK)
		if e, ok := c[q]; ok {
			json.NewEncoder(w).Encode(e)
		}
	}
}
