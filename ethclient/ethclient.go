package ethclient

import (
	"context"
	"log"
	"math/big"
	"sp/config"
	"sp/contracts/SoulPoint_48Club"
	"sp/contracts/calculator"
	"sp/contracts/multicall"
	"sp/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	Client           *ethclient.Client
	contract         = common.HexToAddress("0x3FF1ae3ff05d452EF3E26A883158D7AAD95231dB")
	multicallAdd     = common.HexToAddress("0x41263cBA59EB80dC200F3E2544eda4ed6A90E76C")
	calculatorAdd    = common.HexToAddress("0x988C52043B1151f9502670150df7Cf6008558aF2")
	spabi, _         = SoulPoint_48Club.SoulPoint48ClubMetaData.GetAbi()
	multicallAbi, _  = multicall.MulticallMetaData.GetAbi()
	calculatorAbi, _ = calculator.CalculatorMetaData.GetAbi()

	// 当前区块链每个区块的平均时间(秒)
	AvgBlockTimeSec = 0.75
	// 日区块数
	DayBlockCount = 24 * 60 * 60 / AvgBlockTimeSec
)

func init() {
	ec, err := ethclient.Dial(config.GlobalConfig.RPC)
	if err != nil {
		panic(err)
	}
	Client = ec
}

func GetBlockByTime(tt int64) (blockNumber *big.Int, err error) {
	ctx := context.Background()

	// 获取最新区块号
	latestHeader, err := Client.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}

	latest := latestHeader.Number.Uint64()
	var (
		left  uint64 = 1
		right uint64 = latest
	)

	for left <= right {
		mid := (left + right) / 2
		header, err := Client.HeaderByNumber(ctx, big.NewInt(int64(mid)))
		if err != nil {
			return nil, err
		}

		t := int64(header.Time)

		if t == tt {
			return big.NewInt(int64(mid)), nil
		}

		if t < tt {
			// 时间比目标小 → 去右边找
			left = mid + 1
		} else {
			// 时间比目标大 → 去左边找
			if mid == 0 {
				break
			}
			right = mid - 1
		}
	}
	// 如果没有完全相等的区块时间，返回小于目标时间的最大区块号 right
	return big.NewInt(int64(right)), nil
}

func GetAllMembers(ctx context.Context, ToBlock *big.Int) (addrs []common.Address, err error) {
	logs, err := Client.FilterLogs(ctx, ethereum.FilterQuery{
		Addresses: []common.Address{contract},
		FromBlock: big.NewInt(49660490),
		ToBlock:   ToBlock,
		Topics:    [][]common.Hash{{spabi.Events["Minted"].ID}},
	})
	if err != nil {
		return
	}

	for _, _log := range logs {
		addr := common.BytesToAddress(_log.Topics[1].Bytes())
		if addr == (common.Address{}) {
			continue
		}
		addrs = append(addrs, addr)
	}
	log.Println("filter logs count:", len(logs), "addrs count:", len(addrs))

	return
}

func GetAllSp(ctx context.Context, addrs []common.Address, blockAt *big.Int) ([]types.CalculatorDetail, error) {
	// addrs 最大一次性查询 100 个, 超过 100 个分批查询
	res := []types.CalculatorDetail{}

	for i := 0; i < len(addrs); i += 100 {
		end := min(i+100, len(addrs))
		sp, err := getAllSp(ctx, addrs[i:end], blockAt)
		if err != nil {
			return nil, err
		}
		res = append(res, sp...)
	}

	return res, nil
}

func getAllSp(ctx context.Context, addrs []common.Address, blockAt *big.Int) ([]types.CalculatorDetail, error) {
	mapAddrsSp := []types.CalculatorDetail{}

	calls := []multicall.Struct0{}
	for _, addr := range addrs {
		data, _ := calculatorAbi.Pack("getPointDetail", addr)
		calls = append(calls, multicall.Struct0{Target: calculatorAdd, CallData: data})
	}

	callData, _ := multicallAbi.Pack("aggregate", calls)
	vals, err := Client.CallContract(ctx, ethereum.CallMsg{To: &multicallAdd, Data: callData}, blockAt)
	if err != nil {
		return nil, err
	}

	var result struct {
		BlockNumber *big.Int
		ReturnData  []hexutil.Bytes
	}
	if err := multicallAbi.UnpackIntoInterface(&result, "aggregate", vals); err != nil {
		return nil, err
	}

	for _, data := range result.ReturnData {
		var points types.CalculatorDetail
		err := calculatorAbi.UnpackIntoInterface(&points, "getPointDetail", data)
		if err != nil {
			return nil, err
		}

		mapAddrsSp = append(mapAddrsSp, points)
	}

	return mapAddrsSp, nil
}
