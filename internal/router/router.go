package router

import (
	"log"
	"net/http"
	rc "stuber/internal/router/request_collector"
	"stuber/internal/router/stub"
)

type Router struct {
	requests []*stub.Stub
}

func Run(stubCollection []*stub.Stub) {
	h := make([]*rc.RequestRecord, 0)
	sh := make(map[string][]*rc.RequestRecord)

	for _, s := range stubCollection {
		http.HandleFunc(s.Path, rc.MakeHistoryHandler(&h, rc.MakeCollectorHandler(sh, s.CollectParams, s.MakeStubHandler())))
	}

	http.HandleFunc("/income_request/last", rc.MakeLastRequestHandler(&h))
	http.HandleFunc("/income_request/all", rc.MakeAllRequestsHandler(&h))
	http.HandleFunc("/income_request", rc.MakeGetCollectionHandler(sh))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
