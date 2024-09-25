package repository

import (
	"context"
	"database/sql"
	"errors"

	domain "github.com/dodirepository/warehouse-svc/internal/domain/repository"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type WarehouseRepository struct {
	*domain.BaseStorage
	db *gorm.DB
}

func WarehouseRepositoryHandler(db *gorm.DB) domain.WarehouseRepositoryInterface {
	return &WarehouseRepository{
		BaseStorage: domain.NewBaseRepo(db),
		db:          db,
	}
}

func (w *WarehouseRepository) CreateWarehouse(wh domain.CreateWarehouse) error {
	db := w.db.DB()
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			logrus.WithError(err).Errorln("failed on rollback transaction")
		}
	}(tx)
	q := `INSERT INTO warehouses (name) VALUES (?)`
	vals := []interface{}{
		wh.Name,
	}

	res, err := tx.Exec(q, vals...)
	if err != nil {
		logrus.WithError(err).Errorln("failed on insert warehouse")
		return err
	}
	wh.ID, _ = res.LastInsertId()

	valsDetails := []interface{}{}
	qDetails := `INSERT INTO warehouse_details (warehouse_id,item_id,qty) VALUES `
	for _, v := range wh.WarehouseDetail {
		qDetails += "(?,?,?),"
		valsDetails = append(valsDetails,
			wh.ID,
			v.ItemID,
			v.Qty,
		)
	}
	qDetails = qDetails[0 : len(qDetails)-1]
	_, err = tx.Exec(qDetails, valsDetails...)
	if err != nil {
		logrus.WithError(err).Error("Failed save warehouse detail")
		return err
	}

	return tx.Commit()
}

func (w *WarehouseRepository) CreateWarehouseDetail(warehouseDetails []domain.WarehouseDetail) error {
	db := w.db.DB()
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			logrus.WithError(err).Errorln("failed on rollback transaction")
		}
	}(tx)
	valsDetails := []interface{}{}
	qDetails := `INSERT INTO warehouse_details (warehouse_id,item_id,qty) VALUES `
	for _, v := range warehouseDetails {
		qDetails += "(?,?,?),"
		valsDetails = append(valsDetails,
			v.WarehouseID,
			v.ItemID,
			v.Qty,
		)
	}
	qDetails = qDetails[0 : len(qDetails)-1]
	_, err = tx.Exec(qDetails, valsDetails...)
	if err != nil {
		logrus.WithError(err).Error("Failed save warehouse detail")
		return err
	}

	return tx.Commit()
}

func (w *WarehouseRepository) UpdateStatusWarehouseByID(ID int64, isActive bool) error {
	warehouse := domain.Warehouse{}
	if err := w.db.Model(&warehouse).Table("warehouses").Where("id = ?", ID).Update("is_active", isActive).Error; err != nil {
		logrus.WithError(err).Error("failed to update warehouse")
		return err
	}
	return nil

}

func (w *WarehouseRepository) AddingQtyToWarehouse(warehouseDetails []domain.WarehouseDetail) error {
	db := w.db.DB()
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer func(tx *sql.Tx) {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			logrus.WithError(err).Errorln("failed on rollback transaction")
		}
	}(tx)
	valsDetails := []interface{}{}
	qDetails := `UPDATE warehouse_details SET qty = CASE `
	for _, v := range warehouseDetails {
		qDetails += "WHEN item_id = ? AND warehouse_id = ? THEN ? ELSE qty END,"
		valsDetails = append(valsDetails,
			v.ItemID,
			v.WarehouseID,
			v.Qty,
		)
	}
	qDetails = qDetails[0 : len(qDetails)-1]
	_, err = tx.Exec(qDetails, valsDetails...)
	if err != nil {
		logrus.WithError(err).Error("Failed save warehouse detail")
		return err
	}

	return tx.Commit()
}
func (w *WarehouseRepository) GetWarehouseByID(ID int64) (*domain.Warehouse, error) {
	warehouse := domain.Warehouse{}
	q := w.db.Table("warehouses").Where("id = ?", ID).First(&warehouse)
	if q.Error != nil {
		if errors.Is(q.Error, gorm.ErrRecordNotFound) {
			logrus.WithFields(logrus.Fields{"id": ID}).Error("warehouse not found")
			return nil, nil
		}
		return nil, q.Error
	}
	return &warehouse, nil

}
func (w *WarehouseRepository) GetWarehouseByIDs(ID []int64, itemID int64) ([]domain.Warehouse, error) {
	warehouse := []domain.Warehouse{}
	q := w.db.Table("warehouses as w").
		Select("w.id, w.name, w.is_active").
		Joins("INNER JOIN warehouse_details wd ON w.id = wd.warehouse_id").
		Where("w.id IN (?)", ID).
		Where("wd.item_id = ?", itemID).
		Find(&warehouse)
	if q.Error != nil {
		if errors.Is(q.Error, gorm.ErrRecordNotFound) {
			logrus.WithFields(logrus.Fields{"id": ID}).Error("warehouse not found")
			return nil, nil
		}
		return nil, q.Error
	}
	return warehouse, nil

}
func (w *WarehouseRepository) GetWarehousedetailByItemID(ID, itemID int64) (*domain.WarehouseDetail, error) {
	warehouse := domain.WarehouseDetail{}
	q := w.db.Table("warehouse_details").Where("warehouse_id = ? AND item_id = ?", ID, itemID).First(&warehouse)
	if q.Error != nil {
		if errors.Is(q.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, q.Error
	}
	return &warehouse, nil

}

func (w *WarehouseRepository) TransferItems(ctx context.Context, data domain.WarehouseDetail, tf domain.TrasferItems) error {
	db := w.db.CommonDB()
	if tx := domain.GetTxFromContext(ctx); tx != nil {
		db = tx
	}

	qInsert := `INSERT INTO warehouse_details (warehouse_id,item_id,qty) VALUES (?,?,?)`
	valsInsert := []interface{}{
		tf.ToWarehouseID,
		data.ItemID,
		tf.ToQty,
	}
	result, err := db.Exec(qInsert, valsInsert...)
	if err != nil {
		logrus.WithError(err).Errorln("failed on insert qty warehouse_details")
		return err
	}
	if affected, err := result.RowsAffected(); err != nil || affected == 0 {
		logrus.WithError(err).Errorln("failed on insert qty warehouse_details, affected rows: ", affected)
		return err
	}
	return nil
}

func (w *WarehouseRepository) UpdateDeductQty(ctx context.Context, data domain.WarehouseDetail) error {
	db := w.db.CommonDB()
	if tx := domain.GetTxFromContext(ctx); tx != nil {
		db = tx
	}

	q, err := db.Exec(`UPDATE warehouse_details SET qty = qty - ?  WHERE warehouse_id = ? AND item_id = ?`, data.Qty, data.WarehouseID, data.ItemID)
	if err != nil {
		logrus.WithError(err).Error("failed on update qty warehouse_details")
		return err
	}

	if affected, err := q.RowsAffected(); err != nil || affected == 0 {
		logrus.WithError(err).Errorln("failed on update qty warehouse_details, affected rows: ", affected)
		return err
	}
	return nil

}
func (w *WarehouseRepository) UpdateAddQty(ctx context.Context, data domain.WarehouseDetail) error {
	db := w.db.CommonDB()
	if tx := domain.GetTxFromContext(ctx); tx != nil {
		db = tx
	}

	q, err := db.Exec(`UPDATE warehouse_details SET qty = qty + ?  WHERE warehouse_id = ? AND item_id = ?`, data.Qty, data.WarehouseID, data.ItemID)
	if err != nil {
		logrus.WithError(err).Error("failed on update qty warehouse_details")
		return err
	}

	if affected, err := q.RowsAffected(); err != nil || affected == 0 {
		logrus.WithError(err).Errorln("failed on update qty warehouse_details, affected rows: ", affected)
		return err
	}
	return nil

}
