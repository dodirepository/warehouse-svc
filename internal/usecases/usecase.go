package usecases

import (
	repo "github.com/dodirepository/warehouse-svc/internal/domain/repository"
	usecases "github.com/dodirepository/warehouse-svc/internal/domain/usecases"
)

type WarehouseUsecase struct {
	warehouseRepo repo.WarehouseRepositoryInterface
}

func WarehouseUsecaseHandler(wRepo repo.WarehouseRepositoryInterface) usecases.WarehouseUseCaseInterface {
	return &WarehouseUsecase{
		warehouseRepo: wRepo,
	}
}
