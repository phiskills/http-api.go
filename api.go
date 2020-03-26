package http

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type api struct {
	Port   *int
	Server *http.ServeMux
}

func New() (a *api) {
	a = &api{}
	a.Port = flag.Int("port", 80, "the server port")
	flag.Parse()
	a.Server = http.NewServeMux()
	return
}

func (a *api) Start() {
	port := fmt.Sprintf(":%d", *a.Port)
	if err := http.ListenAndServe(port, a.Server); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func (a *api) Register(pattern string, router *Router) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != pattern {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		method := router.match(r.Method)
		if method == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		method(&Context{Res: w, Req: r})
	}
	a.Server.HandleFunc(pattern, handler)
}
