package http

import (
	"context"
	"os"
	"sync"

	"github.com/dodirepository/warehouse-svc/infrastructure/bootstrap"
	"github.com/dodirepository/warehouse-svc/infrastructure/http"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	godotenv.Load()
}

func Start(ctx context.Context) {
	services := bootstrap.NewService(ctx)
	// defer services.DisconnectAllConnection(ctx)
	run(ctx, services)
}
func StartThread(ctx context.Context, wg *sync.WaitGroup, service *bootstrap.Service) {
	defer wg.Done()
	run(ctx, service)
}

// Start :nodoc:
func run(ctx context.Context, service *bootstrap.Service) {
	httpServer := http.NewHTTPServer(service)
	defer httpServer.Done()

	// _, err := database.DBInit()
	// if err != nil {
	// 	logrus.Fatal("Failed connected to database")
	// }

	logrus.Infof("http server start product svc on port %s", os.Getenv("APP_PORT"))
	if err := httpServer.Run(ctx); err != nil {
		logrus.Info("http server stopped")
	}
}
