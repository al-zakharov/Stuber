package request_collector

import (
	"fmt"
	"net/http"
)

func MakeCollectorHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		fmt.Println(r.Body)
		h.ServeHTTP(w, r)
	}
}
