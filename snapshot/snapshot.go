package snapshot

import (
	"context"
	"errors"
	"log"
	"sp/db"
	"sp/ethclient"
	"sp/types"
	"time"

	"golang.org/x/exp/rand"
	"gorm.io/gorm"
)

func TakeSnapshot(ctx context.Context) {
	sql := db.Server

	for {
		var (
			lastSnapTime     types.SnapTime
			_takeSnapshotNow bool
		)
		tx := sql.Order("created desc").First(&lastSnapTime)
		if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Printf("mysql error: %v", tx.Error)
			continue
		}

		if tx.RowsAffected == 0 {
			_takeSnapshotNow = true
		} else {
			if tt := time.Unix(lastSnapTime.CreatedAt, 0); time.Now().Format("20060102") != tt.Format("20060102") {
				_takeSnapshotNow = true
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

		rand.Seed(uint64(time.Now().UnixMilli()))
		stt := time.Duration(rand.Int63n(720-10)+10) * time.Minute
		log.Printf("sleep %v to next snapshot", stt)
		time.Sleep(stt)
	}
}

func TakeSnapshotNow(ctx context.Context, sql *gorm.DB) error {
	tt := time.Now().Unix()
	tx := sql.Create(&types.SnapTime{CreatedAt: tt})
	if tx.Error != nil {
		return tx.Error
	}

	addrs, err := ethclient.GetAllMembers(ctx)
	if err != nil {
		log.Printf("ethclient.GetAllMembers error: %v", err)
		return err
	}
	maps, err := ethclient.GetAllSp(ctx, addrs)
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
			CreatedAt:     tt,
		}
		if err := sp.Insert(sql, point.Addr); err != nil {
			log.Printf("sp.InsertOne error: %v", err)
			return err
		}
	}

	return nil
}
