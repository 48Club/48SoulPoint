package types

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

type Users struct { // mysql table
	ID      uint64 `gorm:"bigint;primaryKey;autoIncrement;column:id"`
	Address string `gorm:"varchar(42);not null;unique;column:address"`
}

type SoulPoints struct { // mysql table
	ID        uint64 `gorm:"bigint;primaryKey;autoIncrement;column:id" json:"-"`
	Address   string `gorm:"->;-:migration" json:"address"`
	UserID    uint64 `gorm:"bigint;not null;column:user_id" json:"-"`
	Points    uint64 `gorm:"bigint;not null;column:points" json:"points"`
	CreatedAt int64  `gorm:"bigint;column:created" json:"-"`
}

type SnapTime struct { // mysql table
	ID        uint64 `gorm:"bigint;primaryKey;autoIncrement;column:id"`
	CreatedAt int64  `gorm:"bigint;column:created"`
}

type Query struct { // gin query
	Address string `form:"address" json:"address"`
}

type Response struct { // gin response
	Address string `json:"address"`
	Points  uint64 `json:"points"`
}

func (Users) TableName() string {
	return "users"
}

func (u *Users) GetID(db *gorm.DB) error {
	tx := db.First(u, "address = ?", common.HexToAddress(u.Address).String())
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		if err := u.Insert(db); err != nil {
			return err
		}
	}
	return nil
}

func (u *Users) Insert(db *gorm.DB) error {
	tx := db.Create(u)
	return tx.Error
}

func (SoulPoints) TableName() string {
	return "soul_points"
}

func (sp *SoulPoints) Insert(db *gorm.DB, add common.Address) error {
	user := Users{Address: add.String()}
	if err := user.GetID(db); err != nil {
		return err
	}
	sp.UserID = user.ID
	tx := db.Create(sp)
	return tx.Error
}

func (SnapTime) TableName() string {
	return "snap_time"
}
