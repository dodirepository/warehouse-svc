package domain

type Warehouse struct {
	ID       int64  `gorm:"column:id;primaryKey"`
	Name     string `gorm:"column:name"`
	IsActive bool   `gorm:"column:is_active"`
}

type WarehouseDetail struct {
	ID          int64 `gorm:"column:id;primaryKey"`
	WarehouseID int64 `gorm:"column:warehouse_id"`
	ItemID      int64 `gorm:"column:item_id"`
	Qty         int   `gorm:"column:qty"`
}

type CreateWarehouse struct {
	Warehouse
	WarehouseDetail []WarehouseDetail
}

type TrasferItems struct {
	ToQty         int   `gorm:"column:qty"`
	ToWarehouseID int64 `gorm:"column:warehouse_id"`
}
