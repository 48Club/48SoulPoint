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
)

func init() {
	ec, err := ethclient.Dial(config.GlobalConfig.RPC)
	if err != nil {
		panic(err)
	}
	Client = ec
}

func GetAllMembers(ctx context.Context) (addrs []common.Address, err error) {
	logs, err := Client.FilterLogs(ctx, ethereum.FilterQuery{
		Addresses: []common.Address{contract},
		FromBlock: big.NewInt(49660490),
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

func GetAllSp(ctx context.Context, addrs []common.Address) ([]types.CalculatorDetail, error) {
	// addrs 最大一次性查询 100 个, 超过 100 个分批查询
	res := []types.CalculatorDetail{}

	for i := 0; i < len(addrs); i += 100 {
		end := i + 100
		if end > len(addrs) {
			end = len(addrs)
		}
		sp, err := getAllSp(ctx, addrs[i:end])
		if err != nil {
			return nil, err
		}
		res = append(res, sp...)
	}

	return res, nil
}

func getAllSp(ctx context.Context, addrs []common.Address) ([]types.CalculatorDetail, error) {
	mapAddrsSp := []types.CalculatorDetail{}

	calls := []multicall.Struct0{}
	for _, addr := range addrs {
		data, _ := calculatorAbi.Pack("getPointDetail", addr)
		calls = append(calls, multicall.Struct0{Target: calculatorAdd, CallData: data})
	}

	callData, _ := multicallAbi.Pack("aggregate", calls)
	vals, err := Client.CallContract(ctx, ethereum.CallMsg{To: &multicallAdd, Data: callData}, nil)
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
