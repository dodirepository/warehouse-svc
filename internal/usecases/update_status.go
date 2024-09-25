package usecases

import (
	"fmt"
	"net/http"

	domain "github.com/dodirepository/warehouse-svc/internal/domain/usecases"
	"github.com/sirupsen/logrus"
)

func (w WarehouseUsecase) UpdateStatusWarehouseByID(ID int64, isActive bool) *domain.ErrorResponse {
	wh, err := w.warehouseRepo.GetWarehouseByID(ID)
	if err != nil {
		logrus.WithError(err).Error("failed to get warehouse")
		return &domain.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	if wh == nil {
		return &domain.ErrorResponse{
			Message: http.StatusText(http.StatusNotFound),
			Status:  http.StatusNotFound,
		}
	}

	if wh.IsActive == isActive {
		return &domain.ErrorResponse{
			Message: fmt.Sprintf("%s already %t", wh.Name, wh.IsActive),
			Status:  http.StatusUnprocessableEntity,
		}
	}

	if err := w.warehouseRepo.UpdateStatusWarehouseByID(ID, isActive); err != nil {
		logrus.WithError(err).Error("failed to update warehouse")
		return &domain.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return nil
}
