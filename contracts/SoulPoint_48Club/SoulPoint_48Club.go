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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_calculator\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"member\",\"type\":\"address\"}],\"name\":\"Minted\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getMember\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"getMembers\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getMembersCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getPoint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"isMember\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_calculator\",\"type\":\"address\"}],\"name\":\"setCalculator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"upgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// GetMember is a free data retrieval call binding the contract method 0xab3545e5.
//
// Solidity: function getMember(uint256 index) view returns(address)
func (_SoulPoint48Club *SoulPoint48ClubCaller) GetMember(opts *bind.CallOpts, index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _SoulPoint48Club.contract.Call(opts, &out, "getMember", index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetMember is a free data retrieval call binding the contract method 0xab3545e5.
//
// Solidity: function getMember(uint256 index) view returns(address)
func (_SoulPoint48Club *SoulPoint48ClubSession) GetMember(index *big.Int) (common.Address, error) {
	return _SoulPoint48Club.Contract.GetMember(&_SoulPoint48Club.CallOpts, index)
}

// GetMember is a free data retrieval call binding the contract method 0xab3545e5.
//
// Solidity: function getMember(uint256 index) view returns(address)
func (_SoulPoint48Club *SoulPoint48ClubCallerSession) GetMember(index *big.Int) (common.Address, error) {
	return _SoulPoint48Club.Contract.GetMember(&_SoulPoint48Club.CallOpts, index)
}

// GetMembers is a free data retrieval call binding the contract method 0x6e7ba939.
//
// Solidity: function getMembers(uint256 start, uint256 limit) view returns(address[])
func (_SoulPoint48Club *SoulPoint48ClubCaller) GetMembers(opts *bind.CallOpts, start *big.Int, limit *big.Int) ([]common.Address, error) {
	var out []interface{}
	err := _SoulPoint48Club.contract.Call(opts, &out, "getMembers", start, limit)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetMembers is a free data retrieval call binding the contract method 0x6e7ba939.
//
// Solidity: function getMembers(uint256 start, uint256 limit) view returns(address[])
func (_SoulPoint48Club *SoulPoint48ClubSession) GetMembers(start *big.Int, limit *big.Int) ([]common.Address, error) {
	return _SoulPoint48Club.Contract.GetMembers(&_SoulPoint48Club.CallOpts, start, limit)
}

// GetMembers is a free data retrieval call binding the contract method 0x6e7ba939.
//
// Solidity: function getMembers(uint256 start, uint256 limit) view returns(address[])
func (_SoulPoint48Club *SoulPoint48ClubCallerSession) GetMembers(start *big.Int, limit *big.Int) ([]common.Address, error) {
	return _SoulPoint48Club.Contract.GetMembers(&_SoulPoint48Club.CallOpts, start, limit)
}

// GetMembersCount is a free data retrieval call binding the contract method 0x09772f8f.
//
// Solidity: function getMembersCount() view returns(uint256)
func (_SoulPoint48Club *SoulPoint48ClubCaller) GetMembersCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SoulPoint48Club.contract.Call(opts, &out, "getMembersCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMembersCount is a free data retrieval call binding the contract method 0x09772f8f.
//
// Solidity: function getMembersCount() view returns(uint256)
func (_SoulPoint48Club *SoulPoint48ClubSession) GetMembersCount() (*big.Int, error) {
	return _SoulPoint48Club.Contract.GetMembersCount(&_SoulPoint48Club.CallOpts)
}

// GetMembersCount is a free data retrieval call binding the contract method 0x09772f8f.
//
// Solidity: function getMembersCount() view returns(uint256)
func (_SoulPoint48Club *SoulPoint48ClubCallerSession) GetMembersCount() (*big.Int, error) {
	return _SoulPoint48Club.Contract.GetMembersCount(&_SoulPoint48Club.CallOpts)
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

// IsMember is a free data retrieval call binding the contract method 0xa230c524.
//
// Solidity: function isMember(address ) view returns(bool)
func (_SoulPoint48Club *SoulPoint48ClubCaller) IsMember(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _SoulPoint48Club.contract.Call(opts, &out, "isMember", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsMember is a free data retrieval call binding the contract method 0xa230c524.
//
// Solidity: function isMember(address ) view returns(bool)
func (_SoulPoint48Club *SoulPoint48ClubSession) IsMember(arg0 common.Address) (bool, error) {
	return _SoulPoint48Club.Contract.IsMember(&_SoulPoint48Club.CallOpts, arg0)
}

// IsMember is a free data retrieval call binding the contract method 0xa230c524.
//
// Solidity: function isMember(address ) view returns(bool)
func (_SoulPoint48Club *SoulPoint48ClubCallerSession) IsMember(arg0 common.Address) (bool, error) {
	return _SoulPoint48Club.Contract.IsMember(&_SoulPoint48Club.CallOpts, arg0)
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() returns()
func (_SoulPoint48Club *SoulPoint48ClubTransactor) Mint(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SoulPoint48Club.contract.Transact(opts, "mint")
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() returns()
func (_SoulPoint48Club *SoulPoint48ClubSession) Mint() (*types.Transaction, error) {
	return _SoulPoint48Club.Contract.Mint(&_SoulPoint48Club.TransactOpts)
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() returns()
func (_SoulPoint48Club *SoulPoint48ClubTransactorSession) Mint() (*types.Transaction, error) {
	return _SoulPoint48Club.Contract.Mint(&_SoulPoint48Club.TransactOpts)
}

// SetCalculator is a paid mutator transaction binding the contract method 0xc53468f0.
//
// Solidity: function setCalculator(address _calculator) returns()
func (_SoulPoint48Club *SoulPoint48ClubTransactor) SetCalculator(opts *bind.TransactOpts, _calculator common.Address) (*types.Transaction, error) {
	return _SoulPoint48Club.contract.Transact(opts, "setCalculator", _calculator)
}

// SetCalculator is a paid mutator transaction binding the contract method 0xc53468f0.
//
// Solidity: function setCalculator(address _calculator) returns()
func (_SoulPoint48Club *SoulPoint48ClubSession) SetCalculator(_calculator common.Address) (*types.Transaction, error) {
	return _SoulPoint48Club.Contract.SetCalculator(&_SoulPoint48Club.TransactOpts, _calculator)
}

// SetCalculator is a paid mutator transaction binding the contract method 0xc53468f0.
//
// Solidity: function setCalculator(address _calculator) returns()
func (_SoulPoint48Club *SoulPoint48ClubTransactorSession) SetCalculator(_calculator common.Address) (*types.Transaction, error) {
	return _SoulPoint48Club.Contract.SetCalculator(&_SoulPoint48Club.TransactOpts, _calculator)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SoulPoint48Club *SoulPoint48ClubTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SoulPoint48Club.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SoulPoint48Club *SoulPoint48ClubSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SoulPoint48Club.Contract.TransferOwnership(&_SoulPoint48Club.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SoulPoint48Club *SoulPoint48ClubTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SoulPoint48Club.Contract.TransferOwnership(&_SoulPoint48Club.TransactOpts, newOwner)
}

// Upgrade is a paid mutator transaction binding the contract method 0x451450ec.
//
// Solidity: function upgrade(uint256 start, uint256 limit) returns()
func (_SoulPoint48Club *SoulPoint48ClubTransactor) Upgrade(opts *bind.TransactOpts, start *big.Int, limit *big.Int) (*types.Transaction, error) {
	return _SoulPoint48Club.contract.Transact(opts, "upgrade", start, limit)
}

// Upgrade is a paid mutator transaction binding the contract method 0x451450ec.
//
// Solidity: function upgrade(uint256 start, uint256 limit) returns()
func (_SoulPoint48Club *SoulPoint48ClubSession) Upgrade(start *big.Int, limit *big.Int) (*types.Transaction, error) {
	return _SoulPoint48Club.Contract.Upgrade(&_SoulPoint48Club.TransactOpts, start, limit)
}

// Upgrade is a paid mutator transaction binding the contract method 0x451450ec.
//
// Solidity: function upgrade(uint256 start, uint256 limit) returns()
func (_SoulPoint48Club *SoulPoint48ClubTransactorSession) Upgrade(start *big.Int, limit *big.Int) (*types.Transaction, error) {
	return _SoulPoint48Club.Contract.Upgrade(&_SoulPoint48Club.TransactOpts, start, limit)
}

// SoulPoint48ClubMintedIterator is returned from FilterMinted and is used to iterate over the raw logs and unpacked data for Minted events raised by the SoulPoint48Club contract.
type SoulPoint48ClubMintedIterator struct {
	Event *SoulPoint48ClubMinted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SoulPoint48ClubMintedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SoulPoint48ClubMinted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SoulPoint48ClubMinted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SoulPoint48ClubMintedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SoulPoint48ClubMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SoulPoint48ClubMinted represents a Minted event raised by the SoulPoint48Club contract.
type SoulPoint48ClubMinted struct {
	Member common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMinted is a free log retrieval operation binding the contract event 0x90ddedd5a25821bba11fbb98de02ec1f75c1be90ae147d6450ce873e7b78b5d8.
//
// Solidity: event Minted(address indexed member)
func (_SoulPoint48Club *SoulPoint48ClubFilterer) FilterMinted(opts *bind.FilterOpts, member []common.Address) (*SoulPoint48ClubMintedIterator, error) {

	var memberRule []interface{}
	for _, memberItem := range member {
		memberRule = append(memberRule, memberItem)
	}

	logs, sub, err := _SoulPoint48Club.contract.FilterLogs(opts, "Minted", memberRule)
	if err != nil {
		return nil, err
	}
	return &SoulPoint48ClubMintedIterator{contract: _SoulPoint48Club.contract, event: "Minted", logs: logs, sub: sub}, nil
}

// WatchMinted is a free log subscription operation binding the contract event 0x90ddedd5a25821bba11fbb98de02ec1f75c1be90ae147d6450ce873e7b78b5d8.
//
// Solidity: event Minted(address indexed member)
func (_SoulPoint48Club *SoulPoint48ClubFilterer) WatchMinted(opts *bind.WatchOpts, sink chan<- *SoulPoint48ClubMinted, member []common.Address) (event.Subscription, error) {

	var memberRule []interface{}
	for _, memberItem := range member {
		memberRule = append(memberRule, memberItem)
	}

	logs, sub, err := _SoulPoint48Club.contract.WatchLogs(opts, "Minted", memberRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SoulPoint48ClubMinted)
				if err := _SoulPoint48Club.contract.UnpackLog(event, "Minted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMinted is a log parse operation binding the contract event 0x90ddedd5a25821bba11fbb98de02ec1f75c1be90ae147d6450ce873e7b78b5d8.
//
// Solidity: event Minted(address indexed member)
func (_SoulPoint48Club *SoulPoint48ClubFilterer) ParseMinted(log types.Log) (*SoulPoint48ClubMinted, error) {
	event := new(SoulPoint48ClubMinted)
	if err := _SoulPoint48Club.contract.UnpackLog(event, "Minted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
