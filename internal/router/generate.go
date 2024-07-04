package router

import (
	"fmt"
	"log"
	"net/http"
)

type Router struct {
	requests []*Stub
}

func Register(stubCollection []*Stub) {
	for _, stub := range stubCollection {
		http.HandleFunc(stub.Path, stub.makeHandler())
	}

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func (s Stub) makeHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != s.HttpMethod {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(s.Status)
		fmt.Fprint(w, s.Body)
	}
}
