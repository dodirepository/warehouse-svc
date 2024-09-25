package domain

import "net/http"

type HttpApi interface {
	Handlers(w http.ResponseWriter, r *http.Request)
}
