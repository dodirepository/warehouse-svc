package http

import (
	"context"
)

type Server interface {
	Run(context.Context) error
	Done()
}
