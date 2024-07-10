package router

import (
	"fmt"
	"io"
	"log"
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
	Pattern *regexp.Regexp
	Stub    *stub.Stub
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
			Pattern: re,
			Stub:    s,
		})
	}

	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if route.Pattern.MatchString(req.URL.Path) {
			rc.MakeHistoryHandler(&r.h, rc.MakeCollectorHandler(r.sh, route.Stub, route.Stub.MakeStubHandler())).ServeHTTP(w, req)
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
	http.HandleFunc("/dynamic_body", MakeDynamicBodyHandler(router.routes))
	http.Handle("/", router)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		return err
	}

	return nil
}

func MakeDynamicBodyHandler(routes []*Route) http.HandlerFunc {
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
