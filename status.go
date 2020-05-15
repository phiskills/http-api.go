package http

import (
	"net/http"
)

type status int

const (
	Unknown status = iota
	Serving
	NotServing
)

func (s status) String() string {
	return map[status]string{
		Unknown:    "UNKNOWN",
		Serving:    "SERVING",
		NotServing: "NOT_SERVING",
	}[s]
}

func (s status) Value() int {
	return map[status]int{
		Unknown:    http.StatusServiceUnavailable,
		Serving:    http.StatusOK,
		NotServing: http.StatusServiceUnavailable,
	}[s]
}
