package http

import "net/http"

type Context struct {
	Res http.ResponseWriter
	Req *http.Request
}

type Method func(*Context)

type Router struct {
	Get     Method
	Head    Method
	Post    Method
	Put     Method
	Patch   Method
	Delete  Method
	Connect Method
	Options Method
	Trace   Method
}

func (h *Router) match(method string) Method {
	switch method {
	case http.MethodHead:
		return h.Head
	case http.MethodGet:
		return h.Get
	case http.MethodPost:
		return h.Post
	case http.MethodPut:
		return h.Put
	case http.MethodPatch:
		return h.Patch
	case http.MethodDelete:
		return h.Delete
	case http.MethodConnect:
		return h.Connect
	case http.MethodOptions:
		return h.Options
	case http.MethodTrace:
		return h.Trace
	default:
		return nil
	}
}
