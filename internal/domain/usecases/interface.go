package domain

import "context"

type WarehouseUseCaseInterface interface {
	CreateHouse(CreateWarehouse) *ErrorResponse
	AddProduct(AddingProductToWarehouse) *ErrorResponse
	UpdateStatusWarehouseByID(ID int64, isActive bool) *ErrorResponse
	TransferProduct(ctx context.Context, req ProductTransfer) *ErrorResponse
}
