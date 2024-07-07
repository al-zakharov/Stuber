package request_collector

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func MakeCollectorHandler(c map[string][]*RequestRecord, k string, n http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if k != "" {
			//TODO handle error
			b, _ := io.ReadAll(r.Body)
			ub, _ := unmarshalBody(b)

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

func MakeGetCollectionHandler(c map[string][]*RequestRecord, q string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if e, ok := c[q]; ok {
			fmt.Fprint(w, e)
		}
	}
}
