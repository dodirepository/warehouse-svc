package usecases

import (
	"net/http"

	domainRepo "github.com/dodirepository/warehouse-svc/internal/domain/repository"
	usecases "github.com/dodirepository/warehouse-svc/internal/domain/usecases"
	"github.com/sirupsen/logrus"
)

func (w WarehouseUsecase) AddProduct(req usecases.AddingProductToWarehouse) *usecases.ErrorResponse {
	wh, err := w.warehouseRepo.GetWarehouseByID(req.ID)
	if err != nil {
		logrus.WithError(err).Error("failed to get warehouse")
		return &usecases.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	if wh == nil {
		return &usecases.ErrorResponse{
			Message: http.StatusText(http.StatusNotFound),
			Status:  http.StatusNotFound,
		}
	}
	dt := []domainRepo.WarehouseDetail{}
	updAteDt := []domainRepo.WarehouseDetail{}
	for _, v := range req.Items {

		detail, err := w.warehouseRepo.GetWarehousedetailByItemID(wh.ID, v.ID)
		if err != nil {
			return &usecases.ErrorResponse{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		if detail == nil {
			dt = append(dt, domainRepo.WarehouseDetail{
				WarehouseID: wh.ID,
				ItemID:      v.ID,
				Qty:         v.Qty,
			})
		}

		if detail != nil {
			updAteDt = append(updAteDt, domainRepo.WarehouseDetail{
				ID:          detail.ID,
				WarehouseID: wh.ID,
				ItemID:      v.ID,
				Qty:         detail.Qty + v.Qty,
			})
		}

	}
	if len(dt) != 0 {
		if err := w.warehouseRepo.CreateWarehouseDetail(dt); err != nil {
			logrus.WithError(err).Error("failed to create warehouse detail")
			return &usecases.ErrorResponse{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if len(updAteDt) != 0 {
		if err := w.warehouseRepo.AddingQtyToWarehouse(updAteDt); err != nil {
			logrus.WithError(err).Error("failed to update warehouse detail")
			return &usecases.ErrorResponse{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	return nil
}
