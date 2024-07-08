package router

import (
	"fmt"
	"net/http"
	"regexp"
	rc "stuber/internal/router/request_collector"
	"stuber/internal/router/stub"
)

type Router struct {
	routes []*Route
	h      []*rc.RequestRecord
	sh     map[string][]*rc.RequestRecord
}

type Route struct {
	pattern *regexp.Regexp
	stub    *stub.Stub
}

func NewRouter(stubCollection []*stub.Stub) *Router {
	r := &Router{
		routes: make([]*Route, 0),
		h:      make([]*rc.RequestRecord, 0),
		sh:     make(map[string][]*rc.RequestRecord),
	}

	for _, s := range stubCollection {
		rePattern := regexp.MustCompile(`:([^/]+)`).ReplaceAllStringFunc(s.Path, func(m string) string {
			return fmt.Sprintf("(?P<%s>[^/]+)", m[1:])
		})
		re := regexp.MustCompile("^" + rePattern + "$")
		r.routes = append(r.routes, &Route{
			pattern: re,
			stub:    s,
		})
	}

	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if route.pattern.MatchString(req.URL.Path) {
			rc.MakeHistoryHandler(&r.h, rc.MakeCollectorHandler(r.sh, route.stub, route.stub.MakeStubHandler())).ServeHTTP(w, req)
			return
		}
	}
	http.NotFound(w, req)
}

func Run(stubCollection []*stub.Stub) error {
	router := NewRouter(stubCollection)

	http.HandleFunc("/income_request/last", rc.MakeLastRequestHandler(&router.h))
	http.HandleFunc("/income_request/all", rc.MakeAllRequestsHandler(&router.h))
	http.HandleFunc("/income_request", rc.MakeGetCollectionHandler(router.sh))
	http.Handle("/", router)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}

	return nil
}
