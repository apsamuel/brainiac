package common

import (
	"net/http"
)

type Route struct {
	Endpoint string
	Methods  []string
	Handler  http.HandlerFunc
	Auth     string
}
