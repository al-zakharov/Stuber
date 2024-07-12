package server

import (
	"fmt"
	"net/http"
	"regexp"
	"stuber/internal/server/collector"
	"stuber/internal/server/dynamic"
	"stuber/internal/server/route"
	"stuber/internal/server/stub"
	"sync"
)

type Router struct {
	routes []*route.Route
	h      []*collector.RequestRecord
	sh     map[string][]*collector.RequestRecord
	m      sync.RWMutex
}

func NewRouter(stubCollection []*stub.Stub) *Router {
	r := &Router{
		routes: make([]*route.Route, 0),
		h:      make([]*collector.RequestRecord, 0),
		sh:     make(map[string][]*collector.RequestRecord),
	}

	for _, s := range stubCollection {
		rePattern := regexp.MustCompile(`:([^/]+)`).ReplaceAllStringFunc(s.Path, func(m string) string {
			return fmt.Sprintf("(?P<%s>[^/]+)", m[1:])
		})
		re := regexp.MustCompile("^" + rePattern + "$")
		r.routes = append(r.routes, &route.Route{
			Pattern: re,
			Stub:    s,
		})
	}

	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, i := range r.routes {
		if i.Pattern.MatchString(req.URL.Path) {
			handler := collector.MakeHistoryHandler(&r.h, collector.MakeCollectorHandler(r.sh, i.Stub, i.Stub.MakeStubHandler(), &r.m), &r.m)
			handler.ServeHTTP(w, req)
			return
		}
	}
	http.NotFound(w, req)
}

func Run(stubCollection []*stub.Stub) error {
	router := NewRouter(stubCollection)

	http.HandleFunc("/income_request/last", collector.MakeLastRequestHandler(&router.h))
	http.HandleFunc("/income_request/all", collector.MakeAllRequestsHandler(&router.h))
	http.HandleFunc("/income_request", collector.MakeGetCollectionHandler(router.sh))
	http.HandleFunc("/dynamic_body", dynamic.MakeDynamicBodyHandler(router.routes))
	http.Handle("/", router)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}

	return nil
}
