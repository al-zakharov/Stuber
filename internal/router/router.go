package router

import (
	"log"
	"net/http"
	rc "stuber/internal/router/request_collector"
)

type Router struct {
	requests []*Stub
}

func Run(stubCollection []*Stub) {
	h := make([]*rc.RequestRecord, 0)
	sh := make(map[string][]*rc.RequestRecord)

	for _, s := range stubCollection {
		http.HandleFunc(s.Path, rc.MakeHistoryHandler(&h, rc.MakeCollectorHandler(sh, s.RequestIdKey, s.makeStubHandler())))
	}

	http.HandleFunc("/income_request/last", rc.MakeLastRequestHandler(&h))
	http.HandleFunc("/income_request/all", rc.MakeAllRequestsHandler(&h))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
