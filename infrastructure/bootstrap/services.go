package bootstrap

import (
	"context"

	"github.com/dodirepository/warehouse-svc/infrastructure/database"
	"github.com/dodirepository/warehouse-svc/internal/adapters/repository"
	domain "github.com/dodirepository/warehouse-svc/internal/domain/usecases"
	"github.com/dodirepository/warehouse-svc/internal/usecases"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

type Service struct {
	Connection Connection
	Warehouse  domain.WarehouseUseCaseInterface
}

type Connection struct {
	MySQLCon *gorm.DB
}

func NewService(ctx context.Context) *Service {
	mysqlCon := database.GetConnection()
	repo := repository.WarehouseRepositoryHandler(mysqlCon)
	warehouseUsecase := usecases.WarehouseUsecaseHandler(repo)

	return &Service{
		Connection: Connection{
			MySQLCon: mysqlCon,
		},
		Warehouse: warehouseUsecase,
	}
}
