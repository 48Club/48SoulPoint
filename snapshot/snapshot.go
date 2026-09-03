package snapshot

import (
	"context"
	"errors"
	"log"
	"math/big"
	"math/rand/v2"
	"sp/config"
	"sp/db"
	"sp/ethclient"
	"sp/types"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

// retryInterval 数据库/链上失败后的重试等待，避免失败路径空转
const retryInterval = time.Minute

// sleepCtx 等待 d 或 ctx 结束；ctx 结束时返回 false
func sleepCtx(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-t.C:
		return true
	}
}

func TakeSnapshot(ctx context.Context) {
	sql := db.Server

	for {
		var (
			lastSnapTime         types.SnapTime
			_takeSnapshotNow     bool
			_takeSnapshotHistory bool // 数据库内中断, 需要快照历史数据到今天
		)
		tx := sql.Order("created desc").First(&lastSnapTime)
		if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Printf("mysql error: %v", tx.Error)
			if !sleepCtx(ctx, retryInterval) {
				return
			}
			continue
		}

		if tx.RowsAffected == 0 {
			_takeSnapshotNow = true
		} else {
			tt := time.Unix(lastSnapTime.CreatedAt, 0)
			today := time.Now().Format("20060102")

			if today != tt.Format("20060102") {
				_takeSnapshotNow = true
				// 最后快照早于昨天，中间有遗漏天数需要补录（按日期而非时长判定，避免跨日不满 48h 漏天）
				if tt.AddDate(0, 0, 1).Format("20060102") != today {
					_takeSnapshotHistory = true
				}
			}
		}
		if _takeSnapshotHistory {
			tt := time.Unix(lastSnapTime.CreatedAt, 0)
			snapshotFinish := false
			for {
				tt = tt.Add(24 * time.Hour)
				if time.Now().Format("20060102") == tt.Format("20060102") {
					snapshotFinish = true
					break
				}
				beginTt, err := time.ParseInLocation("20060102", tt.Format("20060102"), time.Local)
				if err != nil {
					log.Printf("time.ParseInLocation error: %v", err)
					break
				}
				tt = beginTt.Add(time.Second * time.Duration(rand.Int64N(86400-3600)+1800)) // random 30min~23h
				if err := TakeSnapshotNow(ctx, sql.DB, tt); err != nil {
					log.Printf("takeSnapshotNow error: %v", err)
					break
				}
				log.Printf("take snapshot @ %v", tt)
			}
			if !snapshotFinish {
				// 历史补录未完成，退避后继续外层循环重试
				if !sleepCtx(ctx, retryInterval) {
					return
				}
				continue
			}
		}

		if _takeSnapshotNow {
			if err := TakeSnapshotNow(ctx, sql.DB); err != nil {
				log.Printf("takeSnapshotNow error: %v", err)
				if !sleepCtx(ctx, retryInterval) {
					return
				}
				continue
			}
			log.Println("take snapshot now")

			if err := db.ReloadDetailCache(); err != nil {
				log.Printf("reload detail cache error: %v", err)
			}
		}

		stt := time.Duration(rand.Int64N(720-10)+10) * time.Minute
		log.Printf("sleep %v to next snapshot", stt)
		if !sleepCtx(ctx, stt) {
			return
		}
	}
}

// TakeSnapshotNow 链上抓取在事务外完成，仅写库收进短事务，避免长事务持锁
func TakeSnapshotNow(ctx context.Context, sql *gorm.DB, _tt ...time.Time) error {
	getByHistory := true
	if len(_tt) == 0 {
		getByHistory = false
		_tt = append(_tt, time.Now())
	}
	tt := _tt[0].Unix()

	var block *big.Int = nil
	if getByHistory {
		var err error
		block, err = ethclient.GetBlockByTime(ctx, tt)
		if err != nil {
			return err
		}
	}

	addrs, err := ethclient.GetAllMembers(ctx, block)
	if err != nil {
		log.Printf("ethclient.GetAllMembers error: %v", err)
		return err
	}
	maps, err := ethclient.GetAllSp(ctx, addrs, block)
	if err != nil {
		log.Printf("ethclient.GetAllSp error: %v", err)
		return err
	}

	maps = append(maps, config.GlobalConfig.TestAddr...) // 添加测试钱包地址

	return sql.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&types.SnapTime{CreatedAt: tt}).Error; err != nil {
			return err
		}

		addrs := make([]common.Address, len(maps))
		for i, point := range maps {
			addrs[i] = point.Addr
		}
		ids, err := types.UserIDs(tx, addrs)
		if err != nil {
			return err
		}

		rows := make([]types.SoulPoints, 0, len(maps))
		for _, point := range maps {
			rows = append(rows, types.SoulPoints{
				UserID:        ids[point.Addr.String()],
				Points:        point.Sum().Uint64(),
				KogePoint:     point.KogePoint.Uint64(),
				StakePoint:    point.StakePoint.Uint64(),
				NftPoint:      point.NftPoint.Uint64(),
				BscStakePoint: point.BscStakePoint.Uint64(),
				GovBNBPoint:   point.GovBNBPoint.Uint64(),
				CreatedAt:     tt,
			})
		}
		return tx.CreateInBatches(rows, 500).Error
	})
}
