package dynamic

import (
	"io"
	"log"
	"net/http"
	"stuber/internal/server/route"
)

func MakeDynamicBodyHandler(routes []*route.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		for _, i := range routes {
			if i.Pattern.MatchString(r.URL.Path) {
				b, err := io.ReadAll(r.Body)
				if err != nil {
					log.Printf("Error reading request body: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				if i.Stub.DynamicBody {
					i.Stub.Body = string(b)
				}
			}
		}
	}
}
