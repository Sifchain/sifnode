// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Oracle

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// OracleABI is the input ABI used to generate the binding from.
const OracleABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prophecyID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"}],\"name\":\"LogNewOracleClaim\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prophecyID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prophecyPowerCurrent\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prophecyPowerThreshold\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_submitter\",\"type\":\"address\"}],\"name\":\"LogProphecyProcessed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_power\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_currentValsetVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_validatorCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_totalPower\",\"type\":\"uint256\"}],\"name\":\"LogValidatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_power\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_currentValsetVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_validatorCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_totalPower\",\"type\":\"uint256\"}],\"name\":\"LogValidatorPowerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_power\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_currentValsetVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_validatorCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_totalPower\",\"type\":\"uint256\"}],\"name\":\"LogValidatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newValsetVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_validatorCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_totalPower\",\"type\":\"uint256\"}],\"name\":\"LogValsetReset\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newValsetVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_validatorCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_totalPower\",\"type\":\"uint256\"}],\"name\":\"LogValsetUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_validatorPower\",\"type\":\"uint256\"}],\"name\":\"addValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"consensusThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cosmosBridge\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentValsetVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"signedPower\",\"type\":\"uint256\"}],\"name\":\"getProphecyStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"}],\"name\":\"getValidatorPower\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"hasMadeClaim\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"}],\"name\":\"isActiveValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastNonceSubmitted\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"operator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"oracleClaimValidators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"powers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"prophecyRedeemed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_valsetVersion\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"}],\"name\":\"recoverGas\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"}],\"name\":\"removeValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalPower\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_newValidatorPower\",\"type\":\"uint256\"}],\"name\":\"updateValidatorPower\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_validators\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_powers\",\"type\":\"uint256[]\"}],\"name\":\"updateValset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validatorCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Oracle is an auto generated Go binding around an Ethereum contract.
type Oracle struct {
	OracleCaller     // Read-only binding to the contract
	OracleTransactor // Write-only binding to the contract
	OracleFilterer   // Log filterer for contract events
}

// OracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type OracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OracleSession struct {
	Contract     *Oracle           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OracleCallerSession struct {
	Contract *OracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// OracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OracleTransactorSession struct {
	Contract     *OracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type OracleRaw struct {
	Contract *Oracle // Generic contract binding to access the raw methods on
}

// OracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OracleCallerRaw struct {
	Contract *OracleCaller // Generic read-only contract binding to access the raw methods on
}

// OracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OracleTransactorRaw struct {
	Contract *OracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOracle creates a new instance of Oracle, bound to a specific deployed contract.
func NewOracle(address common.Address, backend bind.ContractBackend) (*Oracle, error) {
	contract, err := bindOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Oracle{OracleCaller: OracleCaller{contract: contract}, OracleTransactor: OracleTransactor{contract: contract}, OracleFilterer: OracleFilterer{contract: contract}}, nil
}

// NewOracleCaller creates a new read-only instance of Oracle, bound to a specific deployed contract.
func NewOracleCaller(address common.Address, caller bind.ContractCaller) (*OracleCaller, error) {
	contract, err := bindOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OracleCaller{contract: contract}, nil
}

// NewOracleTransactor creates a new write-only instance of Oracle, bound to a specific deployed contract.
func NewOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*OracleTransactor, error) {
	contract, err := bindOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OracleTransactor{contract: contract}, nil
}

// NewOracleFilterer creates a new log filterer instance of Oracle, bound to a specific deployed contract.
func NewOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*OracleFilterer, error) {
	contract, err := bindOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OracleFilterer{contract: contract}, nil
}

// bindOracle binds a generic wrapper to an already deployed contract.
func bindOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracle *OracleRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Oracle.Contract.OracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracle *OracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.Contract.OracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracle *OracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracle.Contract.OracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracle *OracleCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Oracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracle *OracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracle *OracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracle.Contract.contract.Transact(opts, method, params...)
}

// ConsensusThreshold is a free data retrieval call binding the contract method 0xf9b0b5b9.
//
// Solidity: function consensusThreshold() view returns(uint256)
func (_Oracle *OracleCaller) ConsensusThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "consensusThreshold")
	return *ret0, err
}

// ConsensusThreshold is a free data retrieval call binding the contract method 0xf9b0b5b9.
//
// Solidity: function consensusThreshold() view returns(uint256)
func (_Oracle *OracleSession) ConsensusThreshold() (*big.Int, error) {
	return _Oracle.Contract.ConsensusThreshold(&_Oracle.CallOpts)
}

// ConsensusThreshold is a free data retrieval call binding the contract method 0xf9b0b5b9.
//
// Solidity: function consensusThreshold() view returns(uint256)
func (_Oracle *OracleCallerSession) ConsensusThreshold() (*big.Int, error) {
	return _Oracle.Contract.ConsensusThreshold(&_Oracle.CallOpts)
}

// CosmosBridge is a free data retrieval call binding the contract method 0xb0e9ef71.
//
// Solidity: function cosmosBridge() view returns(address)
func (_Oracle *OracleCaller) CosmosBridge(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "cosmosBridge")
	return *ret0, err
}

// CosmosBridge is a free data retrieval call binding the contract method 0xb0e9ef71.
//
// Solidity: function cosmosBridge() view returns(address)
func (_Oracle *OracleSession) CosmosBridge() (common.Address, error) {
	return _Oracle.Contract.CosmosBridge(&_Oracle.CallOpts)
}

// CosmosBridge is a free data retrieval call binding the contract method 0xb0e9ef71.
//
// Solidity: function cosmosBridge() view returns(address)
func (_Oracle *OracleCallerSession) CosmosBridge() (common.Address, error) {
	return _Oracle.Contract.CosmosBridge(&_Oracle.CallOpts)
}

// CurrentValsetVersion is a free data retrieval call binding the contract method 0x8d56c37d.
//
// Solidity: function currentValsetVersion() view returns(uint256)
func (_Oracle *OracleCaller) CurrentValsetVersion(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "currentValsetVersion")
	return *ret0, err
}

// CurrentValsetVersion is a free data retrieval call binding the contract method 0x8d56c37d.
//
// Solidity: function currentValsetVersion() view returns(uint256)
func (_Oracle *OracleSession) CurrentValsetVersion() (*big.Int, error) {
	return _Oracle.Contract.CurrentValsetVersion(&_Oracle.CallOpts)
}

// CurrentValsetVersion is a free data retrieval call binding the contract method 0x8d56c37d.
//
// Solidity: function currentValsetVersion() view returns(uint256)
func (_Oracle *OracleCallerSession) CurrentValsetVersion() (*big.Int, error) {
	return _Oracle.Contract.CurrentValsetVersion(&_Oracle.CallOpts)
}

// GetProphecyStatus is a free data retrieval call binding the contract method 0x77491c75.
//
// Solidity: function getProphecyStatus(uint256 signedPower) view returns(bool)
func (_Oracle *OracleCaller) GetProphecyStatus(opts *bind.CallOpts, signedPower *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "getProphecyStatus", signedPower)
	return *ret0, err
}

// GetProphecyStatus is a free data retrieval call binding the contract method 0x77491c75.
//
// Solidity: function getProphecyStatus(uint256 signedPower) view returns(bool)
func (_Oracle *OracleSession) GetProphecyStatus(signedPower *big.Int) (bool, error) {
	return _Oracle.Contract.GetProphecyStatus(&_Oracle.CallOpts, signedPower)
}

// GetProphecyStatus is a free data retrieval call binding the contract method 0x77491c75.
//
// Solidity: function getProphecyStatus(uint256 signedPower) view returns(bool)
func (_Oracle *OracleCallerSession) GetProphecyStatus(signedPower *big.Int) (bool, error) {
	return _Oracle.Contract.GetProphecyStatus(&_Oracle.CallOpts, signedPower)
}

// GetValidatorPower is a free data retrieval call binding the contract method 0x473691a4.
//
// Solidity: function getValidatorPower(address _validatorAddress) view returns(uint256)
func (_Oracle *OracleCaller) GetValidatorPower(opts *bind.CallOpts, _validatorAddress common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "getValidatorPower", _validatorAddress)
	return *ret0, err
}

// GetValidatorPower is a free data retrieval call binding the contract method 0x473691a4.
//
// Solidity: function getValidatorPower(address _validatorAddress) view returns(uint256)
func (_Oracle *OracleSession) GetValidatorPower(_validatorAddress common.Address) (*big.Int, error) {
	return _Oracle.Contract.GetValidatorPower(&_Oracle.CallOpts, _validatorAddress)
}

// GetValidatorPower is a free data retrieval call binding the contract method 0x473691a4.
//
// Solidity: function getValidatorPower(address _validatorAddress) view returns(uint256)
func (_Oracle *OracleCallerSession) GetValidatorPower(_validatorAddress common.Address) (*big.Int, error) {
	return _Oracle.Contract.GetValidatorPower(&_Oracle.CallOpts, _validatorAddress)
}

// HasMadeClaim is a free data retrieval call binding the contract method 0xa219763e.
//
// Solidity: function hasMadeClaim(uint256 , address ) view returns(bool)
func (_Oracle *OracleCaller) HasMadeClaim(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "hasMadeClaim", arg0, arg1)
	return *ret0, err
}

// HasMadeClaim is a free data retrieval call binding the contract method 0xa219763e.
//
// Solidity: function hasMadeClaim(uint256 , address ) view returns(bool)
func (_Oracle *OracleSession) HasMadeClaim(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _Oracle.Contract.HasMadeClaim(&_Oracle.CallOpts, arg0, arg1)
}

// HasMadeClaim is a free data retrieval call binding the contract method 0xa219763e.
//
// Solidity: function hasMadeClaim(uint256 , address ) view returns(bool)
func (_Oracle *OracleCallerSession) HasMadeClaim(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _Oracle.Contract.HasMadeClaim(&_Oracle.CallOpts, arg0, arg1)
}

// IsActiveValidator is a free data retrieval call binding the contract method 0x40550a1c.
//
// Solidity: function isActiveValidator(address _validatorAddress) view returns(bool)
func (_Oracle *OracleCaller) IsActiveValidator(opts *bind.CallOpts, _validatorAddress common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "isActiveValidator", _validatorAddress)
	return *ret0, err
}

// IsActiveValidator is a free data retrieval call binding the contract method 0x40550a1c.
//
// Solidity: function isActiveValidator(address _validatorAddress) view returns(bool)
func (_Oracle *OracleSession) IsActiveValidator(_validatorAddress common.Address) (bool, error) {
	return _Oracle.Contract.IsActiveValidator(&_Oracle.CallOpts, _validatorAddress)
}

// IsActiveValidator is a free data retrieval call binding the contract method 0x40550a1c.
//
// Solidity: function isActiveValidator(address _validatorAddress) view returns(bool)
func (_Oracle *OracleCallerSession) IsActiveValidator(_validatorAddress common.Address) (bool, error) {
	return _Oracle.Contract.IsActiveValidator(&_Oracle.CallOpts, _validatorAddress)
}

// LastNonceSubmitted is a free data retrieval call binding the contract method 0x457c1288.
//
// Solidity: function lastNonceSubmitted() view returns(uint256)
func (_Oracle *OracleCaller) LastNonceSubmitted(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "lastNonceSubmitted")
	return *ret0, err
}

// LastNonceSubmitted is a free data retrieval call binding the contract method 0x457c1288.
//
// Solidity: function lastNonceSubmitted() view returns(uint256)
func (_Oracle *OracleSession) LastNonceSubmitted() (*big.Int, error) {
	return _Oracle.Contract.LastNonceSubmitted(&_Oracle.CallOpts)
}

// LastNonceSubmitted is a free data retrieval call binding the contract method 0x457c1288.
//
// Solidity: function lastNonceSubmitted() view returns(uint256)
func (_Oracle *OracleCallerSession) LastNonceSubmitted() (*big.Int, error) {
	return _Oracle.Contract.LastNonceSubmitted(&_Oracle.CallOpts)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_Oracle *OracleCaller) Operator(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "operator")
	return *ret0, err
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_Oracle *OracleSession) Operator() (common.Address, error) {
	return _Oracle.Contract.Operator(&_Oracle.CallOpts)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_Oracle *OracleCallerSession) Operator() (common.Address, error) {
	return _Oracle.Contract.Operator(&_Oracle.CallOpts)
}

// OracleClaimValidators is a free data retrieval call binding the contract method 0x78ffb1c6.
//
// Solidity: function oracleClaimValidators(uint256 ) view returns(uint256)
func (_Oracle *OracleCaller) OracleClaimValidators(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "oracleClaimValidators", arg0)
	return *ret0, err
}

// OracleClaimValidators is a free data retrieval call binding the contract method 0x78ffb1c6.
//
// Solidity: function oracleClaimValidators(uint256 ) view returns(uint256)
func (_Oracle *OracleSession) OracleClaimValidators(arg0 *big.Int) (*big.Int, error) {
	return _Oracle.Contract.OracleClaimValidators(&_Oracle.CallOpts, arg0)
}

// OracleClaimValidators is a free data retrieval call binding the contract method 0x78ffb1c6.
//
// Solidity: function oracleClaimValidators(uint256 ) view returns(uint256)
func (_Oracle *OracleCallerSession) OracleClaimValidators(arg0 *big.Int) (*big.Int, error) {
	return _Oracle.Contract.OracleClaimValidators(&_Oracle.CallOpts, arg0)
}

// Powers is a free data retrieval call binding the contract method 0x850f43dd.
//
// Solidity: function powers(address , uint256 ) view returns(uint256)
func (_Oracle *OracleCaller) Powers(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "powers", arg0, arg1)
	return *ret0, err
}

// Powers is a free data retrieval call binding the contract method 0x850f43dd.
//
// Solidity: function powers(address , uint256 ) view returns(uint256)
func (_Oracle *OracleSession) Powers(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _Oracle.Contract.Powers(&_Oracle.CallOpts, arg0, arg1)
}

// Powers is a free data retrieval call binding the contract method 0x850f43dd.
//
// Solidity: function powers(address , uint256 ) view returns(uint256)
func (_Oracle *OracleCallerSession) Powers(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _Oracle.Contract.Powers(&_Oracle.CallOpts, arg0, arg1)
}

// ProphecyRedeemed is a free data retrieval call binding the contract method 0x04e12aa5.
//
// Solidity: function prophecyRedeemed(uint256 ) view returns(bool)
func (_Oracle *OracleCaller) ProphecyRedeemed(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "prophecyRedeemed", arg0)
	return *ret0, err
}

// ProphecyRedeemed is a free data retrieval call binding the contract method 0x04e12aa5.
//
// Solidity: function prophecyRedeemed(uint256 ) view returns(bool)
func (_Oracle *OracleSession) ProphecyRedeemed(arg0 *big.Int) (bool, error) {
	return _Oracle.Contract.ProphecyRedeemed(&_Oracle.CallOpts, arg0)
}

// ProphecyRedeemed is a free data retrieval call binding the contract method 0x04e12aa5.
//
// Solidity: function prophecyRedeemed(uint256 ) view returns(bool)
func (_Oracle *OracleCallerSession) ProphecyRedeemed(arg0 *big.Int) (bool, error) {
	return _Oracle.Contract.ProphecyRedeemed(&_Oracle.CallOpts, arg0)
}

// TotalPower is a free data retrieval call binding the contract method 0xdb3ad22c.
//
// Solidity: function totalPower() view returns(uint256)
func (_Oracle *OracleCaller) TotalPower(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "totalPower")
	return *ret0, err
}

// TotalPower is a free data retrieval call binding the contract method 0xdb3ad22c.
//
// Solidity: function totalPower() view returns(uint256)
func (_Oracle *OracleSession) TotalPower() (*big.Int, error) {
	return _Oracle.Contract.TotalPower(&_Oracle.CallOpts)
}

// TotalPower is a free data retrieval call binding the contract method 0xdb3ad22c.
//
// Solidity: function totalPower() view returns(uint256)
func (_Oracle *OracleCallerSession) TotalPower() (*big.Int, error) {
	return _Oracle.Contract.TotalPower(&_Oracle.CallOpts)
}

// ValidatorCount is a free data retrieval call binding the contract method 0x0f43a677.
//
// Solidity: function validatorCount() view returns(uint256)
func (_Oracle *OracleCaller) ValidatorCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "validatorCount")
	return *ret0, err
}

// ValidatorCount is a free data retrieval call binding the contract method 0x0f43a677.
//
// Solidity: function validatorCount() view returns(uint256)
func (_Oracle *OracleSession) ValidatorCount() (*big.Int, error) {
	return _Oracle.Contract.ValidatorCount(&_Oracle.CallOpts)
}

// ValidatorCount is a free data retrieval call binding the contract method 0x0f43a677.
//
// Solidity: function validatorCount() view returns(uint256)
func (_Oracle *OracleCallerSession) ValidatorCount() (*big.Int, error) {
	return _Oracle.Contract.ValidatorCount(&_Oracle.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0x45aaf18c.
//
// Solidity: function validators(address , uint256 ) view returns(bool)
func (_Oracle *OracleCaller) Validators(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Oracle.contract.Call(opts, out, "validators", arg0, arg1)
	return *ret0, err
}

// Validators is a free data retrieval call binding the contract method 0x45aaf18c.
//
// Solidity: function validators(address , uint256 ) view returns(bool)
func (_Oracle *OracleSession) Validators(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _Oracle.Contract.Validators(&_Oracle.CallOpts, arg0, arg1)
}

// Validators is a free data retrieval call binding the contract method 0x45aaf18c.
//
// Solidity: function validators(address , uint256 ) view returns(bool)
func (_Oracle *OracleCallerSession) Validators(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _Oracle.Contract.Validators(&_Oracle.CallOpts, arg0, arg1)
}

// AddValidator is a paid mutator transaction binding the contract method 0xfc6c1f02.
//
// Solidity: function addValidator(address _validatorAddress, uint256 _validatorPower) returns()
func (_Oracle *OracleTransactor) AddValidator(opts *bind.TransactOpts, _validatorAddress common.Address, _validatorPower *big.Int) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "addValidator", _validatorAddress, _validatorPower)
}

// AddValidator is a paid mutator transaction binding the contract method 0xfc6c1f02.
//
// Solidity: function addValidator(address _validatorAddress, uint256 _validatorPower) returns()
func (_Oracle *OracleSession) AddValidator(_validatorAddress common.Address, _validatorPower *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.AddValidator(&_Oracle.TransactOpts, _validatorAddress, _validatorPower)
}

// AddValidator is a paid mutator transaction binding the contract method 0xfc6c1f02.
//
// Solidity: function addValidator(address _validatorAddress, uint256 _validatorPower) returns()
func (_Oracle *OracleTransactorSession) AddValidator(_validatorAddress common.Address, _validatorPower *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.AddValidator(&_Oracle.TransactOpts, _validatorAddress, _validatorPower)
}

// RecoverGas is a paid mutator transaction binding the contract method 0xb5672be3.
//
// Solidity: function recoverGas(uint256 _valsetVersion, address _validatorAddress) returns()
func (_Oracle *OracleTransactor) RecoverGas(opts *bind.TransactOpts, _valsetVersion *big.Int, _validatorAddress common.Address) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "recoverGas", _valsetVersion, _validatorAddress)
}

// RecoverGas is a paid mutator transaction binding the contract method 0xb5672be3.
//
// Solidity: function recoverGas(uint256 _valsetVersion, address _validatorAddress) returns()
func (_Oracle *OracleSession) RecoverGas(_valsetVersion *big.Int, _validatorAddress common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.RecoverGas(&_Oracle.TransactOpts, _valsetVersion, _validatorAddress)
}

// RecoverGas is a paid mutator transaction binding the contract method 0xb5672be3.
//
// Solidity: function recoverGas(uint256 _valsetVersion, address _validatorAddress) returns()
func (_Oracle *OracleTransactorSession) RecoverGas(_valsetVersion *big.Int, _validatorAddress common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.RecoverGas(&_Oracle.TransactOpts, _valsetVersion, _validatorAddress)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(address _validatorAddress) returns()
func (_Oracle *OracleTransactor) RemoveValidator(opts *bind.TransactOpts, _validatorAddress common.Address) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "removeValidator", _validatorAddress)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(address _validatorAddress) returns()
func (_Oracle *OracleSession) RemoveValidator(_validatorAddress common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.RemoveValidator(&_Oracle.TransactOpts, _validatorAddress)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(address _validatorAddress) returns()
func (_Oracle *OracleTransactorSession) RemoveValidator(_validatorAddress common.Address) (*types.Transaction, error) {
	return _Oracle.Contract.RemoveValidator(&_Oracle.TransactOpts, _validatorAddress)
}

// UpdateValidatorPower is a paid mutator transaction binding the contract method 0x2e75293b.
//
// Solidity: function updateValidatorPower(address _validatorAddress, uint256 _newValidatorPower) returns()
func (_Oracle *OracleTransactor) UpdateValidatorPower(opts *bind.TransactOpts, _validatorAddress common.Address, _newValidatorPower *big.Int) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "updateValidatorPower", _validatorAddress, _newValidatorPower)
}

// UpdateValidatorPower is a paid mutator transaction binding the contract method 0x2e75293b.
//
// Solidity: function updateValidatorPower(address _validatorAddress, uint256 _newValidatorPower) returns()
func (_Oracle *OracleSession) UpdateValidatorPower(_validatorAddress common.Address, _newValidatorPower *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.UpdateValidatorPower(&_Oracle.TransactOpts, _validatorAddress, _newValidatorPower)
}

// UpdateValidatorPower is a paid mutator transaction binding the contract method 0x2e75293b.
//
// Solidity: function updateValidatorPower(address _validatorAddress, uint256 _newValidatorPower) returns()
func (_Oracle *OracleTransactorSession) UpdateValidatorPower(_validatorAddress common.Address, _newValidatorPower *big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.UpdateValidatorPower(&_Oracle.TransactOpts, _validatorAddress, _newValidatorPower)
}

// UpdateValset is a paid mutator transaction binding the contract method 0x788cf92f.
//
// Solidity: function updateValset(address[] _validators, uint256[] _powers) returns()
func (_Oracle *OracleTransactor) UpdateValset(opts *bind.TransactOpts, _validators []common.Address, _powers []*big.Int) (*types.Transaction, error) {
	return _Oracle.contract.Transact(opts, "updateValset", _validators, _powers)
}

// UpdateValset is a paid mutator transaction binding the contract method 0x788cf92f.
//
// Solidity: function updateValset(address[] _validators, uint256[] _powers) returns()
func (_Oracle *OracleSession) UpdateValset(_validators []common.Address, _powers []*big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.UpdateValset(&_Oracle.TransactOpts, _validators, _powers)
}

// UpdateValset is a paid mutator transaction binding the contract method 0x788cf92f.
//
// Solidity: function updateValset(address[] _validators, uint256[] _powers) returns()
func (_Oracle *OracleTransactorSession) UpdateValset(_validators []common.Address, _powers []*big.Int) (*types.Transaction, error) {
	return _Oracle.Contract.UpdateValset(&_Oracle.TransactOpts, _validators, _powers)
}

// OracleLogNewOracleClaimIterator is returned from FilterLogNewOracleClaim and is used to iterate over the raw logs and unpacked data for LogNewOracleClaim events raised by the Oracle contract.
type OracleLogNewOracleClaimIterator struct {
	Event *OracleLogNewOracleClaim // Event containing the contract specifics and raw log

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
func (it *OracleLogNewOracleClaimIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleLogNewOracleClaim)
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
		it.Event = new(OracleLogNewOracleClaim)
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
func (it *OracleLogNewOracleClaimIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleLogNewOracleClaimIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleLogNewOracleClaim represents a LogNewOracleClaim event raised by the Oracle contract.
type OracleLogNewOracleClaim struct {
	ProphecyID       *big.Int
	ValidatorAddress common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogNewOracleClaim is a free log retrieval operation binding the contract event 0x668fce9833323940537a9000d512a6c580a1c0797d2b526db0078ee9c5a087a9.
//
// Solidity: event LogNewOracleClaim(uint256 _prophecyID, address _validatorAddress)
func (_Oracle *OracleFilterer) FilterLogNewOracleClaim(opts *bind.FilterOpts) (*OracleLogNewOracleClaimIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "LogNewOracleClaim")
	if err != nil {
		return nil, err
	}
	return &OracleLogNewOracleClaimIterator{contract: _Oracle.contract, event: "LogNewOracleClaim", logs: logs, sub: sub}, nil
}

// WatchLogNewOracleClaim is a free log subscription operation binding the contract event 0x668fce9833323940537a9000d512a6c580a1c0797d2b526db0078ee9c5a087a9.
//
// Solidity: event LogNewOracleClaim(uint256 _prophecyID, address _validatorAddress)
func (_Oracle *OracleFilterer) WatchLogNewOracleClaim(opts *bind.WatchOpts, sink chan<- *OracleLogNewOracleClaim) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "LogNewOracleClaim")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleLogNewOracleClaim)
				if err := _Oracle.contract.UnpackLog(event, "LogNewOracleClaim", log); err != nil {
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

// ParseLogNewOracleClaim is a log parse operation binding the contract event 0x668fce9833323940537a9000d512a6c580a1c0797d2b526db0078ee9c5a087a9.
//
// Solidity: event LogNewOracleClaim(uint256 _prophecyID, address _validatorAddress)
func (_Oracle *OracleFilterer) ParseLogNewOracleClaim(log types.Log) (*OracleLogNewOracleClaim, error) {
	event := new(OracleLogNewOracleClaim)
	if err := _Oracle.contract.UnpackLog(event, "LogNewOracleClaim", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OracleLogProphecyProcessedIterator is returned from FilterLogProphecyProcessed and is used to iterate over the raw logs and unpacked data for LogProphecyProcessed events raised by the Oracle contract.
type OracleLogProphecyProcessedIterator struct {
	Event *OracleLogProphecyProcessed // Event containing the contract specifics and raw log

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
func (it *OracleLogProphecyProcessedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleLogProphecyProcessed)
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
		it.Event = new(OracleLogProphecyProcessed)
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
func (it *OracleLogProphecyProcessedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleLogProphecyProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleLogProphecyProcessed represents a LogProphecyProcessed event raised by the Oracle contract.
type OracleLogProphecyProcessed struct {
	ProphecyID             *big.Int
	ProphecyPowerCurrent   *big.Int
	ProphecyPowerThreshold *big.Int
	Submitter              common.Address
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterLogProphecyProcessed is a free log retrieval operation binding the contract event 0x1d8e3fbd601d9d92db7022fb97f75e132841b94db732dcecb0c93cb31852fcbc.
//
// Solidity: event LogProphecyProcessed(uint256 _prophecyID, uint256 _prophecyPowerCurrent, uint256 _prophecyPowerThreshold, address _submitter)
func (_Oracle *OracleFilterer) FilterLogProphecyProcessed(opts *bind.FilterOpts) (*OracleLogProphecyProcessedIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "LogProphecyProcessed")
	if err != nil {
		return nil, err
	}
	return &OracleLogProphecyProcessedIterator{contract: _Oracle.contract, event: "LogProphecyProcessed", logs: logs, sub: sub}, nil
}

// WatchLogProphecyProcessed is a free log subscription operation binding the contract event 0x1d8e3fbd601d9d92db7022fb97f75e132841b94db732dcecb0c93cb31852fcbc.
//
// Solidity: event LogProphecyProcessed(uint256 _prophecyID, uint256 _prophecyPowerCurrent, uint256 _prophecyPowerThreshold, address _submitter)
func (_Oracle *OracleFilterer) WatchLogProphecyProcessed(opts *bind.WatchOpts, sink chan<- *OracleLogProphecyProcessed) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "LogProphecyProcessed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleLogProphecyProcessed)
				if err := _Oracle.contract.UnpackLog(event, "LogProphecyProcessed", log); err != nil {
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

// ParseLogProphecyProcessed is a log parse operation binding the contract event 0x1d8e3fbd601d9d92db7022fb97f75e132841b94db732dcecb0c93cb31852fcbc.
//
// Solidity: event LogProphecyProcessed(uint256 _prophecyID, uint256 _prophecyPowerCurrent, uint256 _prophecyPowerThreshold, address _submitter)
func (_Oracle *OracleFilterer) ParseLogProphecyProcessed(log types.Log) (*OracleLogProphecyProcessed, error) {
	event := new(OracleLogProphecyProcessed)
	if err := _Oracle.contract.UnpackLog(event, "LogProphecyProcessed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OracleLogValidatorAddedIterator is returned from FilterLogValidatorAdded and is used to iterate over the raw logs and unpacked data for LogValidatorAdded events raised by the Oracle contract.
type OracleLogValidatorAddedIterator struct {
	Event *OracleLogValidatorAdded // Event containing the contract specifics and raw log

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
func (it *OracleLogValidatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleLogValidatorAdded)
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
		it.Event = new(OracleLogValidatorAdded)
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
func (it *OracleLogValidatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleLogValidatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleLogValidatorAdded represents a LogValidatorAdded event raised by the Oracle contract.
type OracleLogValidatorAdded struct {
	Validator            common.Address
	Power                *big.Int
	CurrentValsetVersion *big.Int
	ValidatorCount       *big.Int
	TotalPower           *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterLogValidatorAdded is a free log retrieval operation binding the contract event 0x1a396bcf647502e902dce665d58a0c1b25f982f193ab9a1d0f1500d8d927bf2a.
//
// Solidity: event LogValidatorAdded(address _validator, uint256 _power, uint256 _currentValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Oracle *OracleFilterer) FilterLogValidatorAdded(opts *bind.FilterOpts) (*OracleLogValidatorAddedIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "LogValidatorAdded")
	if err != nil {
		return nil, err
	}
	return &OracleLogValidatorAddedIterator{contract: _Oracle.contract, event: "LogValidatorAdded", logs: logs, sub: sub}, nil
}

// WatchLogValidatorAdded is a free log subscription operation binding the contract event 0x1a396bcf647502e902dce665d58a0c1b25f982f193ab9a1d0f1500d8d927bf2a.
//
// Solidity: event LogValidatorAdded(address _validator, uint256 _power, uint256 _currentValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Oracle *OracleFilterer) WatchLogValidatorAdded(opts *bind.WatchOpts, sink chan<- *OracleLogValidatorAdded) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "LogValidatorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleLogValidatorAdded)
				if err := _Oracle.contract.UnpackLog(event, "LogValidatorAdded", log); err != nil {
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

// ParseLogValidatorAdded is a log parse operation binding the contract event 0x1a396bcf647502e902dce665d58a0c1b25f982f193ab9a1d0f1500d8d927bf2a.
//
// Solidity: event LogValidatorAdded(address _validator, uint256 _power, uint256 _currentValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Oracle *OracleFilterer) ParseLogValidatorAdded(log types.Log) (*OracleLogValidatorAdded, error) {
	event := new(OracleLogValidatorAdded)
	if err := _Oracle.contract.UnpackLog(event, "LogValidatorAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OracleLogValidatorPowerUpdatedIterator is returned from FilterLogValidatorPowerUpdated and is used to iterate over the raw logs and unpacked data for LogValidatorPowerUpdated events raised by the Oracle contract.
type OracleLogValidatorPowerUpdatedIterator struct {
	Event *OracleLogValidatorPowerUpdated // Event containing the contract specifics and raw log

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
func (it *OracleLogValidatorPowerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleLogValidatorPowerUpdated)
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
		it.Event = new(OracleLogValidatorPowerUpdated)
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
func (it *OracleLogValidatorPowerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleLogValidatorPowerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleLogValidatorPowerUpdated represents a LogValidatorPowerUpdated event raised by the Oracle contract.
type OracleLogValidatorPowerUpdated struct {
	Validator            common.Address
	Power                *big.Int
	CurrentValsetVersion *big.Int
	ValidatorCount       *big.Int
	TotalPower           *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterLogValidatorPowerUpdated is a free log retrieval operation binding the contract event 0x335940ce4119f8aae891d73dba74510a3d51f6210134d058237f26e6a31d5340.
//
// Solidity: event LogValidatorPowerUpdated(address _validator, uint256 _power, uint256 _currentValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Oracle *OracleFilterer) FilterLogValidatorPowerUpdated(opts *bind.FilterOpts) (*OracleLogValidatorPowerUpdatedIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "LogValidatorPowerUpdated")
	if err != nil {
		return nil, err
	}
	return &OracleLogValidatorPowerUpdatedIterator{contract: _Oracle.contract, event: "LogValidatorPowerUpdated", logs: logs, sub: sub}, nil
}

// WatchLogValidatorPowerUpdated is a free log subscription operation binding the contract event 0x335940ce4119f8aae891d73dba74510a3d51f6210134d058237f26e6a31d5340.
//
// Solidity: event LogValidatorPowerUpdated(address _validator, uint256 _power, uint256 _currentValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Oracle *OracleFilterer) WatchLogValidatorPowerUpdated(opts *bind.WatchOpts, sink chan<- *OracleLogValidatorPowerUpdated) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "LogValidatorPowerUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleLogValidatorPowerUpdated)
				if err := _Oracle.contract.UnpackLog(event, "LogValidatorPowerUpdated", log); err != nil {
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

// ParseLogValidatorPowerUpdated is a log parse operation binding the contract event 0x335940ce4119f8aae891d73dba74510a3d51f6210134d058237f26e6a31d5340.
//
// Solidity: event LogValidatorPowerUpdated(address _validator, uint256 _power, uint256 _currentValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Oracle *OracleFilterer) ParseLogValidatorPowerUpdated(log types.Log) (*OracleLogValidatorPowerUpdated, error) {
	event := new(OracleLogValidatorPowerUpdated)
	if err := _Oracle.contract.UnpackLog(event, "LogValidatorPowerUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OracleLogValidatorRemovedIterator is returned from FilterLogValidatorRemoved and is used to iterate over the raw logs and unpacked data for LogValidatorRemoved events raised by the Oracle contract.
type OracleLogValidatorRemovedIterator struct {
	Event *OracleLogValidatorRemoved // Event containing the contract specifics and raw log

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
func (it *OracleLogValidatorRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleLogValidatorRemoved)
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
		it.Event = new(OracleLogValidatorRemoved)
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
func (it *OracleLogValidatorRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleLogValidatorRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleLogValidatorRemoved represents a LogValidatorRemoved event raised by the Oracle contract.
type OracleLogValidatorRemoved struct {
	Validator            common.Address
	Power                *big.Int
	CurrentValsetVersion *big.Int
	ValidatorCount       *big.Int
	TotalPower           *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterLogValidatorRemoved is a free log retrieval operation binding the contract event 0x1241fb43a101ff98ab819a1882097d4ccada51ba60f326c1981cc48840f55b8c.
//
// Solidity: event LogValidatorRemoved(address _validator, uint256 _power, uint256 _currentValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Oracle *OracleFilterer) FilterLogValidatorRemoved(opts *bind.FilterOpts) (*OracleLogValidatorRemovedIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "LogValidatorRemoved")
	if err != nil {
		return nil, err
	}
	return &OracleLogValidatorRemovedIterator{contract: _Oracle.contract, event: "LogValidatorRemoved", logs: logs, sub: sub}, nil
}

// WatchLogValidatorRemoved is a free log subscription operation binding the contract event 0x1241fb43a101ff98ab819a1882097d4ccada51ba60f326c1981cc48840f55b8c.
//
// Solidity: event LogValidatorRemoved(address _validator, uint256 _power, uint256 _currentValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Oracle *OracleFilterer) WatchLogValidatorRemoved(opts *bind.WatchOpts, sink chan<- *OracleLogValidatorRemoved) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "LogValidatorRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleLogValidatorRemoved)
				if err := _Oracle.contract.UnpackLog(event, "LogValidatorRemoved", log); err != nil {
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

// ParseLogValidatorRemoved is a log parse operation binding the contract event 0x1241fb43a101ff98ab819a1882097d4ccada51ba60f326c1981cc48840f55b8c.
//
// Solidity: event LogValidatorRemoved(address _validator, uint256 _power, uint256 _currentValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Oracle *OracleFilterer) ParseLogValidatorRemoved(log types.Log) (*OracleLogValidatorRemoved, error) {
	event := new(OracleLogValidatorRemoved)
	if err := _Oracle.contract.UnpackLog(event, "LogValidatorRemoved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OracleLogValsetResetIterator is returned from FilterLogValsetReset and is used to iterate over the raw logs and unpacked data for LogValsetReset events raised by the Oracle contract.
type OracleLogValsetResetIterator struct {
	Event *OracleLogValsetReset // Event containing the contract specifics and raw log

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
func (it *OracleLogValsetResetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleLogValsetReset)
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
		it.Event = new(OracleLogValsetReset)
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
func (it *OracleLogValsetResetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleLogValsetResetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleLogValsetReset represents a LogValsetReset event raised by the Oracle contract.
type OracleLogValsetReset struct {
	NewValsetVersion *big.Int
	ValidatorCount   *big.Int
	TotalPower       *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogValsetReset is a free log retrieval operation binding the contract event 0xd870653e19f161500290fd0c4ca41bf5cf2bcb1ba66448f41c66c512dabd65f2.
//
// Solidity: event LogValsetReset(uint256 _newValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Oracle *OracleFilterer) FilterLogValsetReset(opts *bind.FilterOpts) (*OracleLogValsetResetIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "LogValsetReset")
	if err != nil {
		return nil, err
	}
	return &OracleLogValsetResetIterator{contract: _Oracle.contract, event: "LogValsetReset", logs: logs, sub: sub}, nil
}

// WatchLogValsetReset is a free log subscription operation binding the contract event 0xd870653e19f161500290fd0c4ca41bf5cf2bcb1ba66448f41c66c512dabd65f2.
//
// Solidity: event LogValsetReset(uint256 _newValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Oracle *OracleFilterer) WatchLogValsetReset(opts *bind.WatchOpts, sink chan<- *OracleLogValsetReset) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "LogValsetReset")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleLogValsetReset)
				if err := _Oracle.contract.UnpackLog(event, "LogValsetReset", log); err != nil {
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

// ParseLogValsetReset is a log parse operation binding the contract event 0xd870653e19f161500290fd0c4ca41bf5cf2bcb1ba66448f41c66c512dabd65f2.
//
// Solidity: event LogValsetReset(uint256 _newValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Oracle *OracleFilterer) ParseLogValsetReset(log types.Log) (*OracleLogValsetReset, error) {
	event := new(OracleLogValsetReset)
	if err := _Oracle.contract.UnpackLog(event, "LogValsetReset", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OracleLogValsetUpdatedIterator is returned from FilterLogValsetUpdated and is used to iterate over the raw logs and unpacked data for LogValsetUpdated events raised by the Oracle contract.
type OracleLogValsetUpdatedIterator struct {
	Event *OracleLogValsetUpdated // Event containing the contract specifics and raw log

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
func (it *OracleLogValsetUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleLogValsetUpdated)
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
		it.Event = new(OracleLogValsetUpdated)
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
func (it *OracleLogValsetUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleLogValsetUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleLogValsetUpdated represents a LogValsetUpdated event raised by the Oracle contract.
type OracleLogValsetUpdated struct {
	NewValsetVersion *big.Int
	ValidatorCount   *big.Int
	TotalPower       *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogValsetUpdated is a free log retrieval operation binding the contract event 0x3a7ef0da3179668af8114719645585b5a37092ef2d66f187dcf63d83a221eaa6.
//
// Solidity: event LogValsetUpdated(uint256 _newValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Oracle *OracleFilterer) FilterLogValsetUpdated(opts *bind.FilterOpts) (*OracleLogValsetUpdatedIterator, error) {

	logs, sub, err := _Oracle.contract.FilterLogs(opts, "LogValsetUpdated")
	if err != nil {
		return nil, err
	}
	return &OracleLogValsetUpdatedIterator{contract: _Oracle.contract, event: "LogValsetUpdated", logs: logs, sub: sub}, nil
}

// WatchLogValsetUpdated is a free log subscription operation binding the contract event 0x3a7ef0da3179668af8114719645585b5a37092ef2d66f187dcf63d83a221eaa6.
//
// Solidity: event LogValsetUpdated(uint256 _newValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Oracle *OracleFilterer) WatchLogValsetUpdated(opts *bind.WatchOpts, sink chan<- *OracleLogValsetUpdated) (event.Subscription, error) {

	logs, sub, err := _Oracle.contract.WatchLogs(opts, "LogValsetUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleLogValsetUpdated)
				if err := _Oracle.contract.UnpackLog(event, "LogValsetUpdated", log); err != nil {
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

// ParseLogValsetUpdated is a log parse operation binding the contract event 0x3a7ef0da3179668af8114719645585b5a37092ef2d66f187dcf63d83a221eaa6.
//
// Solidity: event LogValsetUpdated(uint256 _newValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Oracle *OracleFilterer) ParseLogValsetUpdated(log types.Log) (*OracleLogValsetUpdated, error) {
	event := new(OracleLogValsetUpdated)
	if err := _Oracle.contract.UnpackLog(event, "LogValsetUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}
