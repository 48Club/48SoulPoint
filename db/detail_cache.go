package db

import (
	"errors"
	"sync/atomic"
	"time"

	"sp/types"

	"gorm.io/gorm"
)

// detailCache 地址 -> 每日明细（按 created 降序）；启动时加载，写入快照后重建，
// 全量在新 map 上构建后原子替换，读侧无锁。
var detailCache atomic.Pointer[map[string][]types.Detail]

// ReloadDetailCache 全量重建明细缓存
func ReloadDetailCache() error {
	var st []types.SnapTime
	tx := Server.Order("created DESC").Limit(48).Find(&st)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return tx.Error
	}
	var snapshotCount int64
	switch {
	case tx.RowsAffected == 48:
		snapshotCount = 48
	case tx.RowsAffected >= 7:
		snapshotCount = 7
	case tx.RowsAffected >= 2:
		snapshotCount = 2
	default:
		m := map[string][]types.Detail{}
		detailCache.Store(&m)
		return nil
	}

	var rows []types.SoulPoints
	tx = Server.Model(&types.SoulPoints{}).
		Select("user_id, users.address AS address, points, koge_point, stake_point, nft_point, bsc_stake_point, gov_bnb_point, created").
		Joins("JOIN users ON user_id = users.id").
		Where("created BETWEEN ? AND ?", st[snapshotCount-1].CreatedAt, st[0].CreatedAt).
		Order("created DESC").Find(&rows)
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return tx.Error
	}

	m := make(map[string][]types.Detail)
	for _, p := range rows {
		m[p.Address] = append(m[p.Address], types.Detail{
			SnapTime: time.Unix(p.CreatedAt, 0).Format("2006/01/02"),
			Points:   p.Points,
			Koge:     p.KogePoint,
			Stake:    p.StakePoint,
			Nft:      p.NftPoint,
			BscStake: p.BscStakePoint + p.GovBNBPoint,
		})
	}
	detailCache.Store(&m)
	return nil
}

// DetailOf 返回地址的每日明细；缓存未就绪时同步加载一次兜底
func DetailOf(address string) ([]types.Detail, error) {
	m := detailCache.Load()
	if m == nil {
		if err := ReloadDetailCache(); err != nil {
			return nil, err
		}
		m = detailCache.Load()
	}
	return (*m)[address], nil
}
