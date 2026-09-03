package test

import (
	"context"
	"errors"
	"sp/db"
	"sp/snapshot"
	"sp/types"
	"testing"

	"gorm.io/gorm"
)

func TestTakeSnapshotNow(t *testing.T) {
	var lastSnapTime types.SnapTime
	sql := db.Server
	tx := sql.Order("created desc").First(&lastSnapTime)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		t.Error(tx.Error)
	}
	var ctx = context.Background()
	if err := snapshot.TakeSnapshotNow(ctx, sql.DB); err != nil {
		t.Error(err)
	}
}
