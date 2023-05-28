package types

import (
	"net/http"
)

type RegisterRouterType struct {
	Url         string
	HandlerFunc func(w http.ResponseWriter, r *http.Request)
	Methods     []string
}
