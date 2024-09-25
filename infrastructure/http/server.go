package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dodirepository/warehouse-svc/infrastructure/bootstrap"
	"github.com/dodirepository/warehouse-svc/internal/adapters/controllers/router"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	godotenv.Load()
}

type httpServer struct {
	router router.Router
}

// NewHTTPServer :nodoc:
func NewHTTPServer(services *bootstrap.Service) Server {
	return &httpServer{
		router: router.NewRouter(services),
	}
}

// Run :nodoc:
func (h *httpServer) Run(ctx context.Context) error {
	var err error
	readTimeOut, _ := strconv.Atoi(os.Getenv("APP_READ_TIMEOUT_SECOND"))
	writeTimeout, _ := strconv.Atoi(os.Getenv("APP_WRITE_TIMEOUT_SECOND"))
	server := http.Server{
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("APP_ADDRESS"), os.Getenv("APP_PORT")),
		Handler:      h.router.Route(),
		ReadTimeout:  time.Duration(readTimeOut) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
	}

	go func() {
		err = server.ListenAndServe()
		if err != http.ErrServerClosed {
			logrus.Info(fmt.Sprintf("http server error : %s", err.Error()))
		}
		if err != nil {
			logrus.Info(err.Error())
		}
	}()

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer func() {
		cancel()
	}()

	if err = server.Shutdown(ctxShutDown); err != nil {
		logrus.Info(fmt.Sprintf("http server shutdown failed : %s", err.Error()))
	}

	logrus.Info("http server exited successfully")
	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}

// Done :nodoc:
func (h *httpServer) Done() {
	log.Print("service http stopped")
}
