package usecases

import (
	"net/http"

	domainRepo "github.com/dodirepository/warehouse-svc/internal/domain/repository"
	domain "github.com/dodirepository/warehouse-svc/internal/domain/usecases"
	usecases "github.com/dodirepository/warehouse-svc/internal/domain/usecases"
)

func (w WarehouseUsecase) CreateHouse(house usecases.CreateWarehouse) *usecases.ErrorResponse {
	data := domainRepo.CreateWarehouse{
		Warehouse:       domainRepo.Warehouse{},
		WarehouseDetail: []domainRepo.WarehouseDetail{},
	}

	data.Warehouse.Name = house.Name
	for _, v := range house.WarehouseDetail {
		data.WarehouseDetail = append(data.WarehouseDetail, domainRepo.WarehouseDetail{
			ItemID: v.ID,
			Qty:    v.Qty,
		})

	}
	if err := w.warehouseRepo.CreateWarehouse(data); err != nil {
		return &domain.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return nil
}
