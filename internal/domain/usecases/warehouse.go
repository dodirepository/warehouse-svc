package domain

type Warehouse struct {
	Name string `json:"name"`
}

type WarehouseDetail struct {
	ID  int64 `json:"id"`
	Qty int   `json:"qty"`
}

type CreateWarehouse struct {
	Warehouse
	WarehouseDetail []WarehouseDetail `json:"details"`
}
type Items struct {
	ID  int64 `json:"id"`
	Qty int   `json:"qty"`
}

type AddingProductToWarehouse struct {
	ID    int64   `json:"id"`
	Items []Items `json:"items"`
}

type ProductTransfer struct {
	ID            int64 `json:"id"`
	ItemID        int64 `json:"item_id"`
	ToWarehouseID int64 `json:"to_warehouse_id"`
	Qty           int   `json:"qty"`
}
