package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

type Users struct { // mysql table
	ID      uint64 `gorm:"primaryKey;autoIncrement;column:id"`
	Address string `gorm:"size:42;not null;unique;column:address"`
}

type SoulPoints struct { // mysql table
	ID            uint64 `gorm:"primaryKey;autoIncrement;column:id" json:"-"`
	Address       string `gorm:"->;-:migration" json:"address"`
	Count         uint64 `gorm:"->;-:migration" json:"count"` // count of lines, min 1, max 48
	Level         string `gorm:"-" json:"level"`
	UserID        uint64 `gorm:"not null;column:user_id" json:"-"`
	Points        uint64 `gorm:"not null;column:points" json:"points"`
	KogePoint     uint64 `gorm:"column:koge_point;default:0" json:"-"`
	StakePoint    uint64 `gorm:"column:stake_point;default:0" json:"-"`
	NftPoint      uint64 `gorm:"column:nft_point;default:0" json:"-"`
	BscStakePoint uint64 `gorm:"column:bsc_stake_point;default:0" json:"-"`
	GovBNBPoint   uint64 `gorm:"column:gov_bnb_point;default:0" json:"-"`
	CreatedAt     int64  `gorm:"column:created" json:"-"`
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
	ID        uint64 `gorm:"primaryKey;autoIncrement;column:id"`
	CreatedAt int64  `gorm:"column:created"`
}

type Query struct { // gin query
	Address string `form:"address" json:"address"`
	Detail  bool   `form:"detail" json:"detail"`
}

type CalculatorDetail struct {
	Addr                                                        common.Address
	KogePoint, StakePoint, NftPoint, BscStakePoint, GovBNBPoint *big.Int
}

func (c CalculatorDetail) Sum() *big.Int {
	return new(big.Int).SetUint64(
		c.KogePoint.Uint64() + c.StakePoint.Uint64() + c.NftPoint.Uint64() + c.BscStakePoint.Uint64() + c.GovBNBPoint.Uint64(),
	)
}

func (Users) TableName() string {
	return "users"
}

// UserIDs 确保所有地址存在于 users，返回地址 -> id；插入后重读而非依赖批量插入回填的连续自增 id
func UserIDs(db *gorm.DB, addrs []common.Address) (map[string]uint64, error) {
	load := func() (map[string]uint64, error) {
		var users []Users
		if err := db.Find(&users).Error; err != nil {
			return nil, err
		}
		m := make(map[string]uint64, len(users))
		for _, u := range users {
			m[u.Address] = u.ID
		}
		return m, nil
	}

	ids, err := load()
	if err != nil {
		return nil, err
	}
	var missing []Users
	for _, addr := range addrs {
		s := addr.String()
		if _, ok := ids[s]; !ok {
			ids[s] = 0
			missing = append(missing, Users{Address: s})
		}
	}
	if len(missing) == 0 {
		return ids, nil
	}
	if err := db.CreateInBatches(missing, 500).Error; err != nil {
		return nil, err
	}
	return load()
}

func (SoulPoints) TableName() string {
	return "soul_points"
}

func (SnapTime) TableName() string {
	return "snap_time"
}
