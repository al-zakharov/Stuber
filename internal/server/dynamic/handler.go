package dynamic

import (
	"encoding/json"
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

		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var uib Body
		err = json.Unmarshal(b, &uib)
		if err != nil {
			http.Error(w, "wrong json", http.StatusInternalServerError)
			return
		}

		for _, i := range routes {
			if i.Pattern.MatchString(uib.Path) {
				if i.Stub.DynamicBody {
					i.Stub.Body = string(uib.Body)
				}
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		http.Error(w, "path not found", http.StatusNotFound)
		return
	}
}
