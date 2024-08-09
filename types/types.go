package types

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

type Users struct { // mysql table
	ID      uint64 `gorm:"bigint;primaryKey;autoIncrement;column:id"`
	Address string `gorm:"varchar(42);not null;unique;column:address"`
}

type SoulPoints struct { // mysql table
	ID            uint64 `gorm:"bigint;primaryKey;autoIncrement;column:id" json:"-"`
	Address       string `gorm:"->;-:migration" json:"address"`
	UserID        uint64 `gorm:"bigint;not null;column:user_id" json:"-"`
	Points        uint64 `gorm:"bigint;not null;column:points" json:"points"`
	KogePoint     uint64 `gorm:"bigint;column:koge_point;default:0" json:"-"`
	StakePoint    uint64 `gorm:"bigint;column:stake_point;default:0" json:"-"`
	NftPoint      uint64 `gorm:"bigint;column:nft_point;default:0" json:"-"`
	BscStakePoint uint64 `gorm:"bigint;column:bsc_stake_point;default:0" json:"-"`
	CreatedAt     int64  `gorm:"bigint;column:created" json:"-"`
}

type SoulPointsWithDetail struct {
	SoulPoints
	Detail []Detail `json:"detail"`
}

type Detail struct {
	SnapTime string `json:"snap_time"`
	Points   uint64 `json:"points"`
	Koge     uint64 `json:"koge"`
	Stake    uint64 `json:"stake"`
	Nft      uint64 `json:"nft"`
	BscStake uint64 `json:"bsc_stake"`
}

type SnapTime struct { // mysql table
	ID        uint64 `gorm:"bigint;primaryKey;autoIncrement;column:id"`
	CreatedAt int64  `gorm:"bigint;column:created"`
}

type Query struct { // gin query
	Address string `form:"address" json:"address"`
	Detail  bool   `form:"detail" json:"detail"`
}

type CalculatorDetail struct {
	Addr                                           common.Address
	KogePoint, StakePoint, NftPoint, BscStakePoint *big.Int
}

func (c CalculatorDetail) Sum() *big.Int {
	return new(big.Int).SetUint64(
		c.KogePoint.Uint64() + c.StakePoint.Uint64() + c.NftPoint.Uint64() + c.BscStakePoint.Uint64(),
	)
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
