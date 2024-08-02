package ethclient

import (
	"context"
	"math/big"
	"sp/config"
	"sp/contracts/SoulPoint_48Club"
	"sp/contracts/multicall"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	Client               *ethclient.Client
	contract             = common.HexToAddress("0xdB1295c57f62a713f10DdD9DA73a4C6a1700B8b3")
	multicallAdd         = common.HexToAddress("0x41263cBA59EB80dC200F3E2544eda4ed6A90E76C")
	spabi, _             = SoulPoint_48Club.SoulPoint48ClubMetaData.GetAbi()
	multicallAbi, _      = multicall.MulticallMetaData.GetAbi()
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
	return addrs, err
}

func GetAllSp(ctx context.Context, addrs []common.Address) (map[common.Address]*big.Int, error) {
	mapAddrsSp := make(map[common.Address]*big.Int)
	calls := []multicall.Struct0{}

	for _, addr := range addrs {
		data, _ := spabi.Pack("getPoint", addr)
		calls = append(calls, multicall.Struct0{Target: contract, CallData: data})
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

	for i, data := range result.ReturnData {
		var point *big.Int
		err := spabi.UnpackIntoInterface(&point, "getPoint", data)
		if err != nil {
			return nil, err
		}

		mapAddrsSp[addrs[i]] = point
	}

	return mapAddrsSp, nil
}
