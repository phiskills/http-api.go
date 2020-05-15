package http

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	httpPort   = 80
	bufferSize = 10
)

type api struct {
	name   string
	port   int
	router *http.ServeMux
	server *http.Server
	status status
	change chan status
}

func New(name string) (a *api) {
	a = &api{name: name, change: make(chan status, bufferSize)}
	a.port = httpPort
	a.router = http.NewServeMux()
	a.Register("/check", &Router{
		Head:    a.healthCheck,
		Get:     a.healthCheck,
		Options: a.healthCheck,
	})
	return
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
	a.router.HandleFunc(pattern, handler)
}

func (a *api) UseParams() {
	a.parseEnv()
	a.parseFlag()
}

func (a *api) UsePort(port int) {
	a.port = port
}

func (a *api) Start() {
	addr := fmt.Sprintf(":%d", a.port)
	a.server = &http.Server{Addr: addr, Handler: a.router}
	a.Activate()
	if err := a.server.ListenAndServe(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func (a *api) Stop() {
	a.server.Shutdown(context.Background())
	a.Deactivate()
}

func (a *api) Activate() {
	a.updateStatus(Serving)
}

func (a *api) Deactivate() {
	a.updateStatus(NotServing)
}

func (a *api) Server() *http.Server {
	return a.server
}

func (a *api) Name() string {
	return a.name
}

func (a *api) Port() int {
	return a.port
}

func (a *api) Status() status {
	return a.status
}

func (a *api) Change() status {
	return <-a.change
}

func (a *api) healthCheck(ctx *Context) {
	ctx.Res.WriteHeader(a.status.Value())
	ctx.Res.Header().Set("Content-Type", "application/json")
	data := struct {
		Status string `json:"status"`
	}{Status: a.status.String()}
	body, err := json.Marshal(data)
	if err != nil {
		ctx.Res.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx.Res.Write(body)
}

func (a *api) updateStatus(status status) {
	a.status = status
	a.change <- a.status
}

func (a *api) parseFlag() {
	port := *flag.Int("port", 0, "the server port")
	flag.Parse()
	if port != 0 {
		a.port = port
	}
}

func (a *api) parseEnv() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err == nil && port != 0 {
		a.port = port
	}
}
