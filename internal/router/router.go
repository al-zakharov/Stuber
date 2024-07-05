package router

import (
	"log"
	"net/http"
	"stuber/internal/router/request_collector"
)

type Router struct {
	requests []*Stub
}

func Run(stubCollection []*Stub) {
	//h := make([]*request_collector.RequestRecord, 0)
	//sh := make(map[string][]*request_collector.RequestRecord)

	for _, s := range stubCollection {
		http.HandleFunc(s.Path, request_collector.MakeCollectorHandler(s.makeStubHandler()))
	}

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
