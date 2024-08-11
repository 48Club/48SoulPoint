package ethclient

import (
	"context"
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
	Client               *ethclient.Client
	contract             = common.HexToAddress("0x928dC5e31de14114f1486c756C30f39Ab9578A92")
	multicallAdd         = common.HexToAddress("0x41263cBA59EB80dC200F3E2544eda4ed6A90E76C")
	calculatorAdd        = common.HexToAddress("0xa1a46a202cc867826Ed433c361449623c6f28359")
	spabi, _             = SoulPoint_48Club.SoulPoint48ClubMetaData.GetAbi()
	multicallAbi, _      = multicall.MulticallMetaData.GetAbi()
	calculatorAbi, _     = calculator.CalculatorMetaData.GetAbi()
	getAllMembersData, _ = spabi.Pack("getAllMembers")
)

func init() {
	ec, err := ethclient.Dial(config.GlobalConfig.RPC)
	if err != nil {
		panic(err)
	}
	Client = ec
}

func GetAllMembers(ctx context.Context) (addrs []common.Address, err error) {
	hex, err := Client.CallContract(ctx, ethereum.CallMsg{
		To:   &contract,
		Data: getAllMembersData,
	}, nil)
	if err == nil {
		err = spabi.UnpackIntoInterface(&addrs, "getAllMembers", hex)
	}

	for k, v := range addrs {
		if v == (common.Address{}) {
			addrs = addrs[:k]
			break
		}
	}
	return addrs, err
}

func GetAllSp(ctx context.Context, addrs []common.Address) ([]types.CalculatorDetail, error) {
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
