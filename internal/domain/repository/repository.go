package domain

import (
	"context"
	"errors"
	"fmt"

	"database/sql"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

const (
	txKey ctxKey = 1
)

type (
	ctxKey      int
	BaseStorage struct {
		db *gorm.DB
	}
)

func NewBaseRepo(db *gorm.DB) *BaseStorage {
	return &BaseStorage{
		db: db,
	}
}

func GetTxFromContext(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx
	}
	return nil
}

func (b *BaseStorage) WithTransaction(ctx context.Context, fn TransactionFunc) error {
	if parentTx := GetTxFromContext(ctx); parentTx != nil {
		// don't create new tx when context already has one
		return fn(ctx)
	}

	db := b.db.DB()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	defer func(tx *sql.Tx) {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			logrus.WithError(err).Errorln("failed on rollback transaction")
		}
	}(tx)

	ctx = context.WithValue(ctx, txKey, tx)
	if err := fn(ctx); err != nil {
		return fmt.Errorf("error on performing transactional request: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
