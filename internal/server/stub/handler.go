package stub

import (
	"fmt"
	"net/http"
)

func (s *Stub) MakeStubHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != s.HttpMethod {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(s.Status)
		fmt.Fprint(w, s.Body)
	}
}
