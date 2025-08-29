package snapshot

import (
	"context"
	"errors"
	"log"
	"math/big"
	"math/rand/v2"
	"sp/db"
	"sp/ethclient"
	"sp/types"
	"time"

	"gorm.io/gorm"
)

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
			continue
		}

		if tx.RowsAffected == 0 {
			_takeSnapshotNow = true
		} else {
			tt := time.Unix(lastSnapTime.CreatedAt, 0)

			if time.Since(tt) > 24*time.Hour {
				_takeSnapshotHistory = true
			} else if time.Now().Format("20060102") != tt.Format("20060102") {
				_takeSnapshotNow = true
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
				err = sql.Transaction(func(_tx *gorm.DB) error {
					return TakeSnapshotNow(ctx, _tx, tt)
				})
				if err != nil {
					log.Printf("takeSnapshotNow error: %v", err)
					break
				}
				log.Printf("take snapshot @ %v", tt)
			}
			if !snapshotFinish {
				break
			}
		}

		if _takeSnapshotNow {
			err := sql.Transaction(func(_tx *gorm.DB) error {
				return TakeSnapshotNow(ctx, _tx)
			})
			if err != nil {
				log.Printf("takeSnapshotNow error: %v", err)
				continue
			}
			log.Println("take snapshot now")

		}

		stt := time.Duration(rand.Int64N(720-10)+10) * time.Minute
		log.Printf("sleep %v to next snapshot", stt)
		time.Sleep(stt)
	}
}

func TakeSnapshotNow(ctx context.Context, sql *gorm.DB, _tt ...time.Time) error {
	getByHistory := true
	if len(_tt) == 0 {
		getByHistory = false
		_tt = append(_tt, time.Now())
	}
	tt := _tt[0].Unix()
	tx := sql.Create(&types.SnapTime{CreatedAt: tt})
	if tx.Error != nil {
		return tx.Error
	}

	var block *big.Int = nil
	if getByHistory {
		var err error
		block, err = ethclient.GetBlockByTime(tt)
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

	for _, point := range maps {
		sp := types.SoulPoints{
			UserID:        0,
			Points:        point.Sum().Uint64(),
			KogePoint:     point.KogePoint.Uint64(),
			StakePoint:    point.StakePoint.Uint64(),
			NftPoint:      point.NftPoint.Uint64(),
			BscStakePoint: point.BscStakePoint.Uint64(),
			GovBNBPoint:   point.GovBNBPoint.Uint64(),
			CreatedAt:     tt,
		}
		if err := sp.Insert(sql, point.Addr); err != nil {
			log.Printf("sp.InsertOne error: %v", err)
			return err
		}
	}

	return nil
}
