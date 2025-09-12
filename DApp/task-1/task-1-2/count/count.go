// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package count

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

// CountMetaData contains all meta data concerning the Count contract.
var CountMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newValue\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"changedBy\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"CountChanged\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"decrement\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCounterInfo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"increment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastModifiedBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// CountABI is the input ABI used to generate the binding from.
// Deprecated: Use CountMetaData.ABI instead.
var CountABI = CountMetaData.ABI

// Count is an auto generated Go binding around an Ethereum contract.
type Count struct {
	CountCaller     // Read-only binding to the contract
	CountTransactor // Write-only binding to the contract
	CountFilterer   // Log filterer for contract events
}

// CountCaller is an auto generated read-only Go binding around an Ethereum contract.
type CountCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CountTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CountTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CountFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CountFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CountSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CountSession struct {
	Contract     *Count            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CountCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CountCallerSession struct {
	Contract *CountCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// CountTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CountTransactorSession struct {
	Contract     *CountTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CountRaw is an auto generated low-level Go binding around an Ethereum contract.
type CountRaw struct {
	Contract *Count // Generic contract binding to access the raw methods on
}

// CountCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CountCallerRaw struct {
	Contract *CountCaller // Generic read-only contract binding to access the raw methods on
}

// CountTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CountTransactorRaw struct {
	Contract *CountTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCount creates a new instance of Count, bound to a specific deployed contract.
func NewCount(address common.Address, backend bind.ContractBackend) (*Count, error) {
	contract, err := bindCount(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Count{CountCaller: CountCaller{contract: contract}, CountTransactor: CountTransactor{contract: contract}, CountFilterer: CountFilterer{contract: contract}}, nil
}

// NewCountCaller creates a new read-only instance of Count, bound to a specific deployed contract.
func NewCountCaller(address common.Address, caller bind.ContractCaller) (*CountCaller, error) {
	contract, err := bindCount(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CountCaller{contract: contract}, nil
}

// NewCountTransactor creates a new write-only instance of Count, bound to a specific deployed contract.
func NewCountTransactor(address common.Address, transactor bind.ContractTransactor) (*CountTransactor, error) {
	contract, err := bindCount(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CountTransactor{contract: contract}, nil
}

// NewCountFilterer creates a new log filterer instance of Count, bound to a specific deployed contract.
func NewCountFilterer(address common.Address, filterer bind.ContractFilterer) (*CountFilterer, error) {
	contract, err := bindCount(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CountFilterer{contract: contract}, nil
}

// bindCount binds a generic wrapper to an already deployed contract.
func bindCount(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CountMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Count *CountRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Count.Contract.CountCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Count *CountRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Count.Contract.CountTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Count *CountRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Count.Contract.CountTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Count *CountCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Count.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Count *CountTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Count.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Count *CountTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Count.Contract.contract.Transact(opts, method, params...)
}

// GetCount is a free data retrieval call binding the contract method 0xa87d942c.
//
// Solidity: function getCount() view returns(uint256)
func (_Count *CountCaller) GetCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Count.contract.Call(opts, &out, "getCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCount is a free data retrieval call binding the contract method 0xa87d942c.
//
// Solidity: function getCount() view returns(uint256)
func (_Count *CountSession) GetCount() (*big.Int, error) {
	return _Count.Contract.GetCount(&_Count.CallOpts)
}

// GetCount is a free data retrieval call binding the contract method 0xa87d942c.
//
// Solidity: function getCount() view returns(uint256)
func (_Count *CountCallerSession) GetCount() (*big.Int, error) {
	return _Count.Contract.GetCount(&_Count.CallOpts)
}

// GetCounterInfo is a free data retrieval call binding the contract method 0xa2618209.
//
// Solidity: function getCounterInfo() view returns(uint256, address, uint256)
func (_Count *CountCaller) GetCounterInfo(opts *bind.CallOpts) (*big.Int, common.Address, *big.Int, error) {
	var out []interface{}
	err := _Count.contract.Call(opts, &out, "getCounterInfo")

	if err != nil {
		return *new(*big.Int), *new(common.Address), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return out0, out1, out2, err

}

// GetCounterInfo is a free data retrieval call binding the contract method 0xa2618209.
//
// Solidity: function getCounterInfo() view returns(uint256, address, uint256)
func (_Count *CountSession) GetCounterInfo() (*big.Int, common.Address, *big.Int, error) {
	return _Count.Contract.GetCounterInfo(&_Count.CallOpts)
}

// GetCounterInfo is a free data retrieval call binding the contract method 0xa2618209.
//
// Solidity: function getCounterInfo() view returns(uint256, address, uint256)
func (_Count *CountCallerSession) GetCounterInfo() (*big.Int, common.Address, *big.Int, error) {
	return _Count.Contract.GetCounterInfo(&_Count.CallOpts)
}

// LastModifiedBlock is a free data retrieval call binding the contract method 0x3079d307.
//
// Solidity: function lastModifiedBlock() view returns(uint256)
func (_Count *CountCaller) LastModifiedBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Count.contract.Call(opts, &out, "lastModifiedBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastModifiedBlock is a free data retrieval call binding the contract method 0x3079d307.
//
// Solidity: function lastModifiedBlock() view returns(uint256)
func (_Count *CountSession) LastModifiedBlock() (*big.Int, error) {
	return _Count.Contract.LastModifiedBlock(&_Count.CallOpts)
}

// LastModifiedBlock is a free data retrieval call binding the contract method 0x3079d307.
//
// Solidity: function lastModifiedBlock() view returns(uint256)
func (_Count *CountCallerSession) LastModifiedBlock() (*big.Int, error) {
	return _Count.Contract.LastModifiedBlock(&_Count.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Count *CountCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Count.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Count *CountSession) Owner() (common.Address, error) {
	return _Count.Contract.Owner(&_Count.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Count *CountCallerSession) Owner() (common.Address, error) {
	return _Count.Contract.Owner(&_Count.CallOpts)
}

// Decrement is a paid mutator transaction binding the contract method 0x2baeceb7.
//
// Solidity: function decrement() returns()
func (_Count *CountTransactor) Decrement(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Count.contract.Transact(opts, "decrement")
}

// Decrement is a paid mutator transaction binding the contract method 0x2baeceb7.
//
// Solidity: function decrement() returns()
func (_Count *CountSession) Decrement() (*types.Transaction, error) {
	return _Count.Contract.Decrement(&_Count.TransactOpts)
}

// Decrement is a paid mutator transaction binding the contract method 0x2baeceb7.
//
// Solidity: function decrement() returns()
func (_Count *CountTransactorSession) Decrement() (*types.Transaction, error) {
	return _Count.Contract.Decrement(&_Count.TransactOpts)
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_Count *CountTransactor) Increment(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Count.contract.Transact(opts, "increment")
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_Count *CountSession) Increment() (*types.Transaction, error) {
	return _Count.Contract.Increment(&_Count.TransactOpts)
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_Count *CountTransactorSession) Increment() (*types.Transaction, error) {
	return _Count.Contract.Increment(&_Count.TransactOpts)
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_Count *CountTransactor) Reset(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Count.contract.Transact(opts, "reset")
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_Count *CountSession) Reset() (*types.Transaction, error) {
	return _Count.Contract.Reset(&_Count.TransactOpts)
}

// Reset is a paid mutator transaction binding the contract method 0xd826f88f.
//
// Solidity: function reset() returns()
func (_Count *CountTransactorSession) Reset() (*types.Transaction, error) {
	return _Count.Contract.Reset(&_Count.TransactOpts)
}

// CountCountChangedIterator is returned from FilterCountChanged and is used to iterate over the raw logs and unpacked data for CountChanged events raised by the Count contract.
type CountCountChangedIterator struct {
	Event *CountCountChanged // Event containing the contract specifics and raw log

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
func (it *CountCountChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CountCountChanged)
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
		it.Event = new(CountCountChanged)
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
func (it *CountCountChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CountCountChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CountCountChanged represents a CountChanged event raised by the Count contract.
type CountCountChanged struct {
	NewValue    *big.Int
	ChangedBy   common.Address
	BlockNumber *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterCountChanged is a free log retrieval operation binding the contract event 0xc3d684daf59eb687343d4193c26f0c998ad931d6b7ecfe78f94e9104d699346f.
//
// Solidity: event CountChanged(uint256 newValue, address changedBy, uint256 blockNumber)
func (_Count *CountFilterer) FilterCountChanged(opts *bind.FilterOpts) (*CountCountChangedIterator, error) {

	logs, sub, err := _Count.contract.FilterLogs(opts, "CountChanged")
	if err != nil {
		return nil, err
	}
	return &CountCountChangedIterator{contract: _Count.contract, event: "CountChanged", logs: logs, sub: sub}, nil
}

// WatchCountChanged is a free log subscription operation binding the contract event 0xc3d684daf59eb687343d4193c26f0c998ad931d6b7ecfe78f94e9104d699346f.
//
// Solidity: event CountChanged(uint256 newValue, address changedBy, uint256 blockNumber)
func (_Count *CountFilterer) WatchCountChanged(opts *bind.WatchOpts, sink chan<- *CountCountChanged) (event.Subscription, error) {

	logs, sub, err := _Count.contract.WatchLogs(opts, "CountChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CountCountChanged)
				if err := _Count.contract.UnpackLog(event, "CountChanged", log); err != nil {
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

// ParseCountChanged is a log parse operation binding the contract event 0xc3d684daf59eb687343d4193c26f0c998ad931d6b7ecfe78f94e9104d699346f.
//
// Solidity: event CountChanged(uint256 newValue, address changedBy, uint256 blockNumber)
func (_Count *CountFilterer) ParseCountChanged(log types.Log) (*CountCountChanged, error) {
	event := new(CountCountChanged)
	if err := _Count.contract.UnpackLog(event, "CountChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
