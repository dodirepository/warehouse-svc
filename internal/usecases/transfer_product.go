package usecases

import (
	"context"
	"fmt"
	"net/http"

	domain "github.com/dodirepository/warehouse-svc/internal/domain/repository"
	usecases "github.com/dodirepository/warehouse-svc/internal/domain/usecases"
	"github.com/sirupsen/logrus"
)

func (w WarehouseUsecase) TransferProduct(ctx context.Context, req usecases.ProductTransfer) *usecases.ErrorResponse {
	errs := w.valiadateWarehouse(ctx, req)
	if errs != nil {
		return &usecases.ErrorResponse{
			Message: errs.Message,
			Status:  errs.Status,
		}
	}

	err := w.warehouseRepo.WithTransaction(ctx, func(ctx context.Context) error {
		if err := w.warehouseRepo.UpdateDeductQty(ctx, domain.WarehouseDetail{
			WarehouseID: req.ID,
			ItemID:      req.ItemID,
			Qty:         req.Qty,
		}); err != nil {
			return err
		}
		if err := w.warehouseRepo.UpdateAddQty(ctx, domain.WarehouseDetail{
			WarehouseID: req.ToWarehouseID,
			ItemID:      req.ItemID,
			Qty:         req.Qty,
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return &usecases.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return nil

}

func (w WarehouseUsecase) valiadateWarehouse(ctx context.Context, req usecases.ProductTransfer) *usecases.ErrorResponse {
	ids := []int64{req.ID, req.ToWarehouseID}
	wh, err := w.warehouseRepo.GetWarehouseByIDs(ids, req.ItemID)
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

	if len(wh) != len(ids) {
		return &usecases.ErrorResponse{
			Message: "target warehouses not found or item on target warehouse not found",
			Status:  http.StatusNotFound,
		}
	}

	whTo, err := w.warehouseRepo.GetWarehouseByID(req.ToWarehouseID)
	if err != nil {
		logrus.WithError(err).Error("failed to get warehouse")
		return &usecases.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	if whTo == nil {
		return &usecases.ErrorResponse{
			Message: fmt.Sprintf("target warehouse with id: %s not found", req.ToWarehouseID),
			Status:  http.StatusNotFound,
		}
	}

	if !whTo.IsActive {
		return &usecases.ErrorResponse{
			Message: fmt.Sprintf("%s already %t", whTo.Name, whTo.IsActive),
			Status:  http.StatusUnprocessableEntity,
		}
	}

	detail, err := w.warehouseRepo.GetWarehousedetailByItemID(req.ToWarehouseID, req.ItemID)
	if err != nil {
		logrus.WithError(err).Error("failed to get warehouse detail")
		return &usecases.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	if detail != nil {
		if detail.Qty < req.Qty {
			return &usecases.ErrorResponse{
				Message: "qty not enough",
				Status:  http.StatusUnprocessableEntity,
			}

		}
	}
	return nil
}
