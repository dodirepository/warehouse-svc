package domain

import "context"

type WarehouseRepositoryInterface interface {
	TransactionRepository
	CreateWarehouse(CreateWarehouse) error
	CreateWarehouseDetail([]WarehouseDetail) error
	UpdateStatusWarehouseByID(ID int64, isActive bool) error
	GetWarehouseByID(ID int64) (*Warehouse, error)
	GetWarehouseByIDs(ID []int64, itemID int64) ([]Warehouse, error)
	GetWarehousedetailByItemID(ID, itemID int64) (*WarehouseDetail, error)
	AddingQtyToWarehouse([]WarehouseDetail) error
	TransferItems(ctx context.Context, data WarehouseDetail, UpdateQty TrasferItems) error
	UpdateDeductQty(ctx context.Context, data WarehouseDetail) error
	UpdateAddQty(ctx context.Context, data WarehouseDetail) error
}

type TransactionFunc func(ctx context.Context) error

type TransactionRepository interface {
	WithTransaction(ctx context.Context, fn TransactionFunc) error
}
