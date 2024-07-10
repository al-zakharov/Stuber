package dynamic

import (
	"io"
	"log"
	"net/http"
	"stuber/internal/server/route"
)

func MakeDynamicBodyHandler(routes []*route.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, i := range routes {
			if i.Pattern.MatchString(r.URL.Path) {
				b, err := io.ReadAll(r.Body)
				if err != nil {
					log.Printf("Error reading request body: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				i.Stub.Body = string(b)
			}
		}
	}
}
