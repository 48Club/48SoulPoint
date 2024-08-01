// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package SoulPoint_48Club

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// SoulPoint48ClubMetaData contains all meta data concerning the SoulPoint48Club contract.
var SoulPoint48ClubMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getAllMembers\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getPoint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// SoulPoint48ClubABI is the input ABI used to generate the binding from.
// Deprecated: Use SoulPoint48ClubMetaData.ABI instead.
var SoulPoint48ClubABI = SoulPoint48ClubMetaData.ABI

// SoulPoint48Club is an auto generated Go binding around an Ethereum contract.
type SoulPoint48Club struct {
	SoulPoint48ClubCaller     // Read-only binding to the contract
	SoulPoint48ClubTransactor // Write-only binding to the contract
	SoulPoint48ClubFilterer   // Log filterer for contract events
}

// SoulPoint48ClubCaller is an auto generated read-only Go binding around an Ethereum contract.
type SoulPoint48ClubCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SoulPoint48ClubTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SoulPoint48ClubTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SoulPoint48ClubFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SoulPoint48ClubFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SoulPoint48ClubSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SoulPoint48ClubSession struct {
	Contract     *SoulPoint48Club  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SoulPoint48ClubCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SoulPoint48ClubCallerSession struct {
	Contract *SoulPoint48ClubCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// SoulPoint48ClubTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SoulPoint48ClubTransactorSession struct {
	Contract     *SoulPoint48ClubTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// SoulPoint48ClubRaw is an auto generated low-level Go binding around an Ethereum contract.
type SoulPoint48ClubRaw struct {
	Contract *SoulPoint48Club // Generic contract binding to access the raw methods on
}

// SoulPoint48ClubCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SoulPoint48ClubCallerRaw struct {
	Contract *SoulPoint48ClubCaller // Generic read-only contract binding to access the raw methods on
}

// SoulPoint48ClubTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SoulPoint48ClubTransactorRaw struct {
	Contract *SoulPoint48ClubTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSoulPoint48Club creates a new instance of SoulPoint48Club, bound to a specific deployed contract.
func NewSoulPoint48Club(address common.Address, backend bind.ContractBackend) (*SoulPoint48Club, error) {
	contract, err := bindSoulPoint48Club(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SoulPoint48Club{SoulPoint48ClubCaller: SoulPoint48ClubCaller{contract: contract}, SoulPoint48ClubTransactor: SoulPoint48ClubTransactor{contract: contract}, SoulPoint48ClubFilterer: SoulPoint48ClubFilterer{contract: contract}}, nil
}

// NewSoulPoint48ClubCaller creates a new read-only instance of SoulPoint48Club, bound to a specific deployed contract.
func NewSoulPoint48ClubCaller(address common.Address, caller bind.ContractCaller) (*SoulPoint48ClubCaller, error) {
	contract, err := bindSoulPoint48Club(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SoulPoint48ClubCaller{contract: contract}, nil
}

// NewSoulPoint48ClubTransactor creates a new write-only instance of SoulPoint48Club, bound to a specific deployed contract.
func NewSoulPoint48ClubTransactor(address common.Address, transactor bind.ContractTransactor) (*SoulPoint48ClubTransactor, error) {
	contract, err := bindSoulPoint48Club(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SoulPoint48ClubTransactor{contract: contract}, nil
}

// NewSoulPoint48ClubFilterer creates a new log filterer instance of SoulPoint48Club, bound to a specific deployed contract.
func NewSoulPoint48ClubFilterer(address common.Address, filterer bind.ContractFilterer) (*SoulPoint48ClubFilterer, error) {
	contract, err := bindSoulPoint48Club(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SoulPoint48ClubFilterer{contract: contract}, nil
}

// bindSoulPoint48Club binds a generic wrapper to an already deployed contract.
func bindSoulPoint48Club(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SoulPoint48ClubMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SoulPoint48Club *SoulPoint48ClubRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SoulPoint48Club.Contract.SoulPoint48ClubCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SoulPoint48Club *SoulPoint48ClubRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SoulPoint48Club.Contract.SoulPoint48ClubTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SoulPoint48Club *SoulPoint48ClubRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SoulPoint48Club.Contract.SoulPoint48ClubTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SoulPoint48Club *SoulPoint48ClubCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SoulPoint48Club.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SoulPoint48Club *SoulPoint48ClubTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SoulPoint48Club.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SoulPoint48Club *SoulPoint48ClubTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SoulPoint48Club.Contract.contract.Transact(opts, method, params...)
}

// GetAllMembers is a free data retrieval call binding the contract method 0x7c0f6b35.
//
// Solidity: function getAllMembers() view returns(address[])
func (_SoulPoint48Club *SoulPoint48ClubCaller) GetAllMembers(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _SoulPoint48Club.contract.Call(opts, &out, "getAllMembers")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetAllMembers is a free data retrieval call binding the contract method 0x7c0f6b35.
//
// Solidity: function getAllMembers() view returns(address[])
func (_SoulPoint48Club *SoulPoint48ClubSession) GetAllMembers() ([]common.Address, error) {
	return _SoulPoint48Club.Contract.GetAllMembers(&_SoulPoint48Club.CallOpts)
}

// GetAllMembers is a free data retrieval call binding the contract method 0x7c0f6b35.
//
// Solidity: function getAllMembers() view returns(address[])
func (_SoulPoint48Club *SoulPoint48ClubCallerSession) GetAllMembers() ([]common.Address, error) {
	return _SoulPoint48Club.Contract.GetAllMembers(&_SoulPoint48Club.CallOpts)
}

// GetPoint is a free data retrieval call binding the contract method 0x4ff531b6.
//
// Solidity: function getPoint(address account) view returns(uint256)
func (_SoulPoint48Club *SoulPoint48ClubCaller) GetPoint(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SoulPoint48Club.contract.Call(opts, &out, "getPoint", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPoint is a free data retrieval call binding the contract method 0x4ff531b6.
//
// Solidity: function getPoint(address account) view returns(uint256)
func (_SoulPoint48Club *SoulPoint48ClubSession) GetPoint(account common.Address) (*big.Int, error) {
	return _SoulPoint48Club.Contract.GetPoint(&_SoulPoint48Club.CallOpts, account)
}

// GetPoint is a free data retrieval call binding the contract method 0x4ff531b6.
//
// Solidity: function getPoint(address account) view returns(uint256)
func (_SoulPoint48Club *SoulPoint48ClubCallerSession) GetPoint(account common.Address) (*big.Int, error) {
	return _SoulPoint48Club.Contract.GetPoint(&_SoulPoint48Club.CallOpts, account)
}
