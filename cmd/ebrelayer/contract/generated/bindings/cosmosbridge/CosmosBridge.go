// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package CosmosBridge

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

// CosmosBridgeClaimData is an auto generated low-level Go binding around an user-defined struct.
type CosmosBridgeClaimData struct {
	CosmosSender         []byte
	CosmosSenderSequence *big.Int
	EthereumReceiver     common.Address
	TokenAddress         common.Address
	Amount               *big.Int
	DoublePeg            bool
	Nonce                *big.Int
}

// CosmosBridgeSignatureData is an auto generated low-level Go binding around an user-defined struct.
type CosmosBridgeSignatureData struct {
	Signer common.Address
	V      uint8
	R      [32]byte
	S      [32]byte
}

// CosmosBridgeABI is the input ABI used to generate the binding from.
const CosmosBridgeABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"bridgeBank\",\"type\":\"address\"}],\"name\":\"LogBridgeBankSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"sourceChainDescriptor\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sourceContractAddress\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"bridgeTokenAddress\",\"type\":\"address\"}],\"name\":\"LogNewBridgeTokenCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prophecyID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"}],\"name\":\"LogNewOracleClaim\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"prophecyID\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"ethereumReceiver\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"LogNewProphecyClaim\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"prophecyID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"LogProphecyCompleted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prophecyID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prophecyPowerCurrent\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_prophecyPowerThreshold\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_submitter\",\"type\":\"address\"}],\"name\":\"LogProphecyProcessed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_power\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_currentValsetVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_validatorCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_totalPower\",\"type\":\"uint256\"}],\"name\":\"LogValidatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_power\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_currentValsetVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_validatorCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_totalPower\",\"type\":\"uint256\"}],\"name\":\"LogValidatorPowerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_power\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_currentValsetVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_validatorCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_totalPower\",\"type\":\"uint256\"}],\"name\":\"LogValidatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newValsetVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_validatorCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_totalPower\",\"type\":\"uint256\"}],\"name\":\"LogValsetReset\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newValsetVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_validatorCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_totalPower\",\"type\":\"uint256\"}],\"name\":\"LogValsetUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_validatorPower\",\"type\":\"uint256\"}],\"name\":\"addValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"sigs\",\"type\":\"bytes32[]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"cosmosSender\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"cosmosSenderSequence\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"ethereumReceiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"doublePeg\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"nonce\",\"type\":\"uint128\"}],\"internalType\":\"structCosmosBridge.ClaimData[]\",\"name\":\"claims\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_s\",\"type\":\"bytes32\"}],\"internalType\":\"structCosmosBridge.SignatureData[][]\",\"name\":\"signatureData\",\"type\":\"tuple[][]\"}],\"name\":\"batchSubmitProphecyClaimAggregatedSigs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridgeBank\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newOperator\",\"type\":\"address\"}],\"name\":\"changeOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"consensusThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"cosmosBridge\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"sourceChainTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"decimals\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"chainDescriptor\",\"type\":\"uint256\"}],\"name\":\"createNewBridgeToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentValsetVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_s\",\"type\":\"bytes32\"}],\"internalType\":\"structCosmosBridge.SignatureData[]\",\"name\":\"validators\",\"type\":\"tuple[]\"}],\"name\":\"findDup\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"cosmosSender\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"cosmosSenderSequence\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"ethereumReceiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"doublePeg\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"nonce\",\"type\":\"uint128\"}],\"name\":\"getProphecyID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"signedPower\",\"type\":\"uint256\"}],\"name\":\"getProphecyStatus\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_s\",\"type\":\"bytes32\"}],\"internalType\":\"structCosmosBridge.SignatureData[]\",\"name\":\"validators\",\"type\":\"tuple[]\"}],\"name\":\"getSignedPower\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"}],\"name\":\"getValidatorPower\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"hasBridgeBank\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"hasMadeClaim\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_consensusThreshold\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"_initValidators\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_initPowers\",\"type\":\"uint256[]\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"}],\"name\":\"isActiveValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastNonceSubmitted\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"operator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracle\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"oracleClaimValidators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"powers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"prophecyClaims\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"ethereumReceiver\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"prophecyRedeemed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_valsetVersion\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"}],\"name\":\"recoverGas\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"}],\"name\":\"removeValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_bridgeBank\",\"type\":\"address\"}],\"name\":\"setBridgeBank\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"sourceAddressToDestinationAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashDigest\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"cosmosSender\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"cosmosSenderSequence\",\"type\":\"uint256\"},{\"internalType\":\"addresspayable\",\"name\":\"ethereumReceiver\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"doublePeg\",\"type\":\"bool\"},{\"internalType\":\"uint128\",\"name\":\"nonce\",\"type\":\"uint128\"}],\"internalType\":\"structCosmosBridge.ClaimData\",\"name\":\"claimData\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_s\",\"type\":\"bytes32\"}],\"internalType\":\"structCosmosBridge.SignatureData[]\",\"name\":\"signatureData\",\"type\":\"tuple[]\"}],\"name\":\"submitProphecyClaimAggregatedSigs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalPower\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_newValidatorPower\",\"type\":\"uint256\"}],\"name\":\"updateValidatorPower\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_validators\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_powers\",\"type\":\"uint256[]\"}],\"name\":\"updateValset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"validatorCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"valset\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// CosmosBridge is an auto generated Go binding around an Ethereum contract.
type CosmosBridge struct {
	CosmosBridgeCaller     // Read-only binding to the contract
	CosmosBridgeTransactor // Write-only binding to the contract
	CosmosBridgeFilterer   // Log filterer for contract events
}

// CosmosBridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type CosmosBridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CosmosBridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CosmosBridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CosmosBridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CosmosBridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CosmosBridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CosmosBridgeSession struct {
	Contract     *CosmosBridge     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CosmosBridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CosmosBridgeCallerSession struct {
	Contract *CosmosBridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// CosmosBridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CosmosBridgeTransactorSession struct {
	Contract     *CosmosBridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// CosmosBridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type CosmosBridgeRaw struct {
	Contract *CosmosBridge // Generic contract binding to access the raw methods on
}

// CosmosBridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CosmosBridgeCallerRaw struct {
	Contract *CosmosBridgeCaller // Generic read-only contract binding to access the raw methods on
}

// CosmosBridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CosmosBridgeTransactorRaw struct {
	Contract *CosmosBridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCosmosBridge creates a new instance of CosmosBridge, bound to a specific deployed contract.
func NewCosmosBridge(address common.Address, backend bind.ContractBackend) (*CosmosBridge, error) {
	contract, err := bindCosmosBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CosmosBridge{CosmosBridgeCaller: CosmosBridgeCaller{contract: contract}, CosmosBridgeTransactor: CosmosBridgeTransactor{contract: contract}, CosmosBridgeFilterer: CosmosBridgeFilterer{contract: contract}}, nil
}

// NewCosmosBridgeCaller creates a new read-only instance of CosmosBridge, bound to a specific deployed contract.
func NewCosmosBridgeCaller(address common.Address, caller bind.ContractCaller) (*CosmosBridgeCaller, error) {
	contract, err := bindCosmosBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeCaller{contract: contract}, nil
}

// NewCosmosBridgeTransactor creates a new write-only instance of CosmosBridge, bound to a specific deployed contract.
func NewCosmosBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*CosmosBridgeTransactor, error) {
	contract, err := bindCosmosBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeTransactor{contract: contract}, nil
}

// NewCosmosBridgeFilterer creates a new log filterer instance of CosmosBridge, bound to a specific deployed contract.
func NewCosmosBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*CosmosBridgeFilterer, error) {
	contract, err := bindCosmosBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeFilterer{contract: contract}, nil
}

// bindCosmosBridge binds a generic wrapper to an already deployed contract.
func bindCosmosBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CosmosBridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CosmosBridge *CosmosBridgeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CosmosBridge.Contract.CosmosBridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CosmosBridge *CosmosBridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CosmosBridge.Contract.CosmosBridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CosmosBridge *CosmosBridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CosmosBridge.Contract.CosmosBridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CosmosBridge *CosmosBridgeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _CosmosBridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CosmosBridge *CosmosBridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CosmosBridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CosmosBridge *CosmosBridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CosmosBridge.Contract.contract.Transact(opts, method, params...)
}

// BridgeBank is a free data retrieval call binding the contract method 0x0e41f373.
//
// Solidity: function bridgeBank() view returns(address)
func (_CosmosBridge *CosmosBridgeCaller) BridgeBank(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "bridgeBank")
	return *ret0, err
}

// BridgeBank is a free data retrieval call binding the contract method 0x0e41f373.
//
// Solidity: function bridgeBank() view returns(address)
func (_CosmosBridge *CosmosBridgeSession) BridgeBank() (common.Address, error) {
	return _CosmosBridge.Contract.BridgeBank(&_CosmosBridge.CallOpts)
}

// BridgeBank is a free data retrieval call binding the contract method 0x0e41f373.
//
// Solidity: function bridgeBank() view returns(address)
func (_CosmosBridge *CosmosBridgeCallerSession) BridgeBank() (common.Address, error) {
	return _CosmosBridge.Contract.BridgeBank(&_CosmosBridge.CallOpts)
}

// ConsensusThreshold is a free data retrieval call binding the contract method 0xf9b0b5b9.
//
// Solidity: function consensusThreshold() view returns(uint256)
func (_CosmosBridge *CosmosBridgeCaller) ConsensusThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "consensusThreshold")
	return *ret0, err
}

// ConsensusThreshold is a free data retrieval call binding the contract method 0xf9b0b5b9.
//
// Solidity: function consensusThreshold() view returns(uint256)
func (_CosmosBridge *CosmosBridgeSession) ConsensusThreshold() (*big.Int, error) {
	return _CosmosBridge.Contract.ConsensusThreshold(&_CosmosBridge.CallOpts)
}

// ConsensusThreshold is a free data retrieval call binding the contract method 0xf9b0b5b9.
//
// Solidity: function consensusThreshold() view returns(uint256)
func (_CosmosBridge *CosmosBridgeCallerSession) ConsensusThreshold() (*big.Int, error) {
	return _CosmosBridge.Contract.ConsensusThreshold(&_CosmosBridge.CallOpts)
}

// CosmosBridge is a free data retrieval call binding the contract method 0xb0e9ef71.
//
// Solidity: function cosmosBridge() view returns(address)
func (_CosmosBridge *CosmosBridgeCaller) CosmosBridge(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "cosmosBridge")
	return *ret0, err
}

// CosmosBridge is a free data retrieval call binding the contract method 0xb0e9ef71.
//
// Solidity: function cosmosBridge() view returns(address)
func (_CosmosBridge *CosmosBridgeSession) CosmosBridge() (common.Address, error) {
	return _CosmosBridge.Contract.CosmosBridge(&_CosmosBridge.CallOpts)
}

// CosmosBridge is a free data retrieval call binding the contract method 0xb0e9ef71.
//
// Solidity: function cosmosBridge() view returns(address)
func (_CosmosBridge *CosmosBridgeCallerSession) CosmosBridge() (common.Address, error) {
	return _CosmosBridge.Contract.CosmosBridge(&_CosmosBridge.CallOpts)
}

// CurrentValsetVersion is a free data retrieval call binding the contract method 0x8d56c37d.
//
// Solidity: function currentValsetVersion() view returns(uint256)
func (_CosmosBridge *CosmosBridgeCaller) CurrentValsetVersion(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "currentValsetVersion")
	return *ret0, err
}

// CurrentValsetVersion is a free data retrieval call binding the contract method 0x8d56c37d.
//
// Solidity: function currentValsetVersion() view returns(uint256)
func (_CosmosBridge *CosmosBridgeSession) CurrentValsetVersion() (*big.Int, error) {
	return _CosmosBridge.Contract.CurrentValsetVersion(&_CosmosBridge.CallOpts)
}

// CurrentValsetVersion is a free data retrieval call binding the contract method 0x8d56c37d.
//
// Solidity: function currentValsetVersion() view returns(uint256)
func (_CosmosBridge *CosmosBridgeCallerSession) CurrentValsetVersion() (*big.Int, error) {
	return _CosmosBridge.Contract.CurrentValsetVersion(&_CosmosBridge.CallOpts)
}

// FindDup is a free data retrieval call binding the contract method 0x48426c99.
//
// Solidity: function findDup((address,uint8,bytes32,bytes32)[] validators) pure returns(bool)
func (_CosmosBridge *CosmosBridgeCaller) FindDup(opts *bind.CallOpts, validators []CosmosBridgeSignatureData) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "findDup", validators)
	return *ret0, err
}

// FindDup is a free data retrieval call binding the contract method 0x48426c99.
//
// Solidity: function findDup((address,uint8,bytes32,bytes32)[] validators) pure returns(bool)
func (_CosmosBridge *CosmosBridgeSession) FindDup(validators []CosmosBridgeSignatureData) (bool, error) {
	return _CosmosBridge.Contract.FindDup(&_CosmosBridge.CallOpts, validators)
}

// FindDup is a free data retrieval call binding the contract method 0x48426c99.
//
// Solidity: function findDup((address,uint8,bytes32,bytes32)[] validators) pure returns(bool)
func (_CosmosBridge *CosmosBridgeCallerSession) FindDup(validators []CosmosBridgeSignatureData) (bool, error) {
	return _CosmosBridge.Contract.FindDup(&_CosmosBridge.CallOpts, validators)
}

// GetProphecyID is a free data retrieval call binding the contract method 0x2d8a21d5.
//
// Solidity: function getProphecyID(bytes cosmosSender, uint256 cosmosSenderSequence, address ethereumReceiver, address tokenAddress, uint256 amount, bool doublePeg, uint128 nonce) pure returns(uint256)
func (_CosmosBridge *CosmosBridgeCaller) GetProphecyID(opts *bind.CallOpts, cosmosSender []byte, cosmosSenderSequence *big.Int, ethereumReceiver common.Address, tokenAddress common.Address, amount *big.Int, doublePeg bool, nonce *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "getProphecyID", cosmosSender, cosmosSenderSequence, ethereumReceiver, tokenAddress, amount, doublePeg, nonce)
	return *ret0, err
}

// GetProphecyID is a free data retrieval call binding the contract method 0x2d8a21d5.
//
// Solidity: function getProphecyID(bytes cosmosSender, uint256 cosmosSenderSequence, address ethereumReceiver, address tokenAddress, uint256 amount, bool doublePeg, uint128 nonce) pure returns(uint256)
func (_CosmosBridge *CosmosBridgeSession) GetProphecyID(cosmosSender []byte, cosmosSenderSequence *big.Int, ethereumReceiver common.Address, tokenAddress common.Address, amount *big.Int, doublePeg bool, nonce *big.Int) (*big.Int, error) {
	return _CosmosBridge.Contract.GetProphecyID(&_CosmosBridge.CallOpts, cosmosSender, cosmosSenderSequence, ethereumReceiver, tokenAddress, amount, doublePeg, nonce)
}

// GetProphecyID is a free data retrieval call binding the contract method 0x2d8a21d5.
//
// Solidity: function getProphecyID(bytes cosmosSender, uint256 cosmosSenderSequence, address ethereumReceiver, address tokenAddress, uint256 amount, bool doublePeg, uint128 nonce) pure returns(uint256)
func (_CosmosBridge *CosmosBridgeCallerSession) GetProphecyID(cosmosSender []byte, cosmosSenderSequence *big.Int, ethereumReceiver common.Address, tokenAddress common.Address, amount *big.Int, doublePeg bool, nonce *big.Int) (*big.Int, error) {
	return _CosmosBridge.Contract.GetProphecyID(&_CosmosBridge.CallOpts, cosmosSender, cosmosSenderSequence, ethereumReceiver, tokenAddress, amount, doublePeg, nonce)
}

// GetProphecyStatus is a free data retrieval call binding the contract method 0x77491c75.
//
// Solidity: function getProphecyStatus(uint256 signedPower) view returns(bool)
func (_CosmosBridge *CosmosBridgeCaller) GetProphecyStatus(opts *bind.CallOpts, signedPower *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "getProphecyStatus", signedPower)
	return *ret0, err
}

// GetProphecyStatus is a free data retrieval call binding the contract method 0x77491c75.
//
// Solidity: function getProphecyStatus(uint256 signedPower) view returns(bool)
func (_CosmosBridge *CosmosBridgeSession) GetProphecyStatus(signedPower *big.Int) (bool, error) {
	return _CosmosBridge.Contract.GetProphecyStatus(&_CosmosBridge.CallOpts, signedPower)
}

// GetProphecyStatus is a free data retrieval call binding the contract method 0x77491c75.
//
// Solidity: function getProphecyStatus(uint256 signedPower) view returns(bool)
func (_CosmosBridge *CosmosBridgeCallerSession) GetProphecyStatus(signedPower *big.Int) (bool, error) {
	return _CosmosBridge.Contract.GetProphecyStatus(&_CosmosBridge.CallOpts, signedPower)
}

// GetSignedPower is a free data retrieval call binding the contract method 0x3561516d.
//
// Solidity: function getSignedPower((address,uint8,bytes32,bytes32)[] validators) view returns(uint256)
func (_CosmosBridge *CosmosBridgeCaller) GetSignedPower(opts *bind.CallOpts, validators []CosmosBridgeSignatureData) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "getSignedPower", validators)
	return *ret0, err
}

// GetSignedPower is a free data retrieval call binding the contract method 0x3561516d.
//
// Solidity: function getSignedPower((address,uint8,bytes32,bytes32)[] validators) view returns(uint256)
func (_CosmosBridge *CosmosBridgeSession) GetSignedPower(validators []CosmosBridgeSignatureData) (*big.Int, error) {
	return _CosmosBridge.Contract.GetSignedPower(&_CosmosBridge.CallOpts, validators)
}

// GetSignedPower is a free data retrieval call binding the contract method 0x3561516d.
//
// Solidity: function getSignedPower((address,uint8,bytes32,bytes32)[] validators) view returns(uint256)
func (_CosmosBridge *CosmosBridgeCallerSession) GetSignedPower(validators []CosmosBridgeSignatureData) (*big.Int, error) {
	return _CosmosBridge.Contract.GetSignedPower(&_CosmosBridge.CallOpts, validators)
}

// GetValidatorPower is a free data retrieval call binding the contract method 0x473691a4.
//
// Solidity: function getValidatorPower(address _validatorAddress) view returns(uint256)
func (_CosmosBridge *CosmosBridgeCaller) GetValidatorPower(opts *bind.CallOpts, _validatorAddress common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "getValidatorPower", _validatorAddress)
	return *ret0, err
}

// GetValidatorPower is a free data retrieval call binding the contract method 0x473691a4.
//
// Solidity: function getValidatorPower(address _validatorAddress) view returns(uint256)
func (_CosmosBridge *CosmosBridgeSession) GetValidatorPower(_validatorAddress common.Address) (*big.Int, error) {
	return _CosmosBridge.Contract.GetValidatorPower(&_CosmosBridge.CallOpts, _validatorAddress)
}

// GetValidatorPower is a free data retrieval call binding the contract method 0x473691a4.
//
// Solidity: function getValidatorPower(address _validatorAddress) view returns(uint256)
func (_CosmosBridge *CosmosBridgeCallerSession) GetValidatorPower(_validatorAddress common.Address) (*big.Int, error) {
	return _CosmosBridge.Contract.GetValidatorPower(&_CosmosBridge.CallOpts, _validatorAddress)
}

// HasBridgeBank is a free data retrieval call binding the contract method 0x69294a4e.
//
// Solidity: function hasBridgeBank() view returns(bool)
func (_CosmosBridge *CosmosBridgeCaller) HasBridgeBank(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "hasBridgeBank")
	return *ret0, err
}

// HasBridgeBank is a free data retrieval call binding the contract method 0x69294a4e.
//
// Solidity: function hasBridgeBank() view returns(bool)
func (_CosmosBridge *CosmosBridgeSession) HasBridgeBank() (bool, error) {
	return _CosmosBridge.Contract.HasBridgeBank(&_CosmosBridge.CallOpts)
}

// HasBridgeBank is a free data retrieval call binding the contract method 0x69294a4e.
//
// Solidity: function hasBridgeBank() view returns(bool)
func (_CosmosBridge *CosmosBridgeCallerSession) HasBridgeBank() (bool, error) {
	return _CosmosBridge.Contract.HasBridgeBank(&_CosmosBridge.CallOpts)
}

// HasMadeClaim is a free data retrieval call binding the contract method 0xa219763e.
//
// Solidity: function hasMadeClaim(uint256 , address ) view returns(bool)
func (_CosmosBridge *CosmosBridgeCaller) HasMadeClaim(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "hasMadeClaim", arg0, arg1)
	return *ret0, err
}

// HasMadeClaim is a free data retrieval call binding the contract method 0xa219763e.
//
// Solidity: function hasMadeClaim(uint256 , address ) view returns(bool)
func (_CosmosBridge *CosmosBridgeSession) HasMadeClaim(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _CosmosBridge.Contract.HasMadeClaim(&_CosmosBridge.CallOpts, arg0, arg1)
}

// HasMadeClaim is a free data retrieval call binding the contract method 0xa219763e.
//
// Solidity: function hasMadeClaim(uint256 , address ) view returns(bool)
func (_CosmosBridge *CosmosBridgeCallerSession) HasMadeClaim(arg0 *big.Int, arg1 common.Address) (bool, error) {
	return _CosmosBridge.Contract.HasMadeClaim(&_CosmosBridge.CallOpts, arg0, arg1)
}

// IsActiveValidator is a free data retrieval call binding the contract method 0x40550a1c.
//
// Solidity: function isActiveValidator(address _validatorAddress) view returns(bool)
func (_CosmosBridge *CosmosBridgeCaller) IsActiveValidator(opts *bind.CallOpts, _validatorAddress common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "isActiveValidator", _validatorAddress)
	return *ret0, err
}

// IsActiveValidator is a free data retrieval call binding the contract method 0x40550a1c.
//
// Solidity: function isActiveValidator(address _validatorAddress) view returns(bool)
func (_CosmosBridge *CosmosBridgeSession) IsActiveValidator(_validatorAddress common.Address) (bool, error) {
	return _CosmosBridge.Contract.IsActiveValidator(&_CosmosBridge.CallOpts, _validatorAddress)
}

// IsActiveValidator is a free data retrieval call binding the contract method 0x40550a1c.
//
// Solidity: function isActiveValidator(address _validatorAddress) view returns(bool)
func (_CosmosBridge *CosmosBridgeCallerSession) IsActiveValidator(_validatorAddress common.Address) (bool, error) {
	return _CosmosBridge.Contract.IsActiveValidator(&_CosmosBridge.CallOpts, _validatorAddress)
}

// LastNonceSubmitted is a free data retrieval call binding the contract method 0x457c1288.
//
// Solidity: function lastNonceSubmitted() view returns(uint256)
func (_CosmosBridge *CosmosBridgeCaller) LastNonceSubmitted(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "lastNonceSubmitted")
	return *ret0, err
}

// LastNonceSubmitted is a free data retrieval call binding the contract method 0x457c1288.
//
// Solidity: function lastNonceSubmitted() view returns(uint256)
func (_CosmosBridge *CosmosBridgeSession) LastNonceSubmitted() (*big.Int, error) {
	return _CosmosBridge.Contract.LastNonceSubmitted(&_CosmosBridge.CallOpts)
}

// LastNonceSubmitted is a free data retrieval call binding the contract method 0x457c1288.
//
// Solidity: function lastNonceSubmitted() view returns(uint256)
func (_CosmosBridge *CosmosBridgeCallerSession) LastNonceSubmitted() (*big.Int, error) {
	return _CosmosBridge.Contract.LastNonceSubmitted(&_CosmosBridge.CallOpts)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_CosmosBridge *CosmosBridgeCaller) Operator(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "operator")
	return *ret0, err
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_CosmosBridge *CosmosBridgeSession) Operator() (common.Address, error) {
	return _CosmosBridge.Contract.Operator(&_CosmosBridge.CallOpts)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_CosmosBridge *CosmosBridgeCallerSession) Operator() (common.Address, error) {
	return _CosmosBridge.Contract.Operator(&_CosmosBridge.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_CosmosBridge *CosmosBridgeCaller) Oracle(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "oracle")
	return *ret0, err
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_CosmosBridge *CosmosBridgeSession) Oracle() (common.Address, error) {
	return _CosmosBridge.Contract.Oracle(&_CosmosBridge.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_CosmosBridge *CosmosBridgeCallerSession) Oracle() (common.Address, error) {
	return _CosmosBridge.Contract.Oracle(&_CosmosBridge.CallOpts)
}

// OracleClaimValidators is a free data retrieval call binding the contract method 0x78ffb1c6.
//
// Solidity: function oracleClaimValidators(uint256 ) view returns(uint256)
func (_CosmosBridge *CosmosBridgeCaller) OracleClaimValidators(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "oracleClaimValidators", arg0)
	return *ret0, err
}

// OracleClaimValidators is a free data retrieval call binding the contract method 0x78ffb1c6.
//
// Solidity: function oracleClaimValidators(uint256 ) view returns(uint256)
func (_CosmosBridge *CosmosBridgeSession) OracleClaimValidators(arg0 *big.Int) (*big.Int, error) {
	return _CosmosBridge.Contract.OracleClaimValidators(&_CosmosBridge.CallOpts, arg0)
}

// OracleClaimValidators is a free data retrieval call binding the contract method 0x78ffb1c6.
//
// Solidity: function oracleClaimValidators(uint256 ) view returns(uint256)
func (_CosmosBridge *CosmosBridgeCallerSession) OracleClaimValidators(arg0 *big.Int) (*big.Int, error) {
	return _CosmosBridge.Contract.OracleClaimValidators(&_CosmosBridge.CallOpts, arg0)
}

// Powers is a free data retrieval call binding the contract method 0x850f43dd.
//
// Solidity: function powers(address , uint256 ) view returns(uint256)
func (_CosmosBridge *CosmosBridgeCaller) Powers(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "powers", arg0, arg1)
	return *ret0, err
}

// Powers is a free data retrieval call binding the contract method 0x850f43dd.
//
// Solidity: function powers(address , uint256 ) view returns(uint256)
func (_CosmosBridge *CosmosBridgeSession) Powers(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _CosmosBridge.Contract.Powers(&_CosmosBridge.CallOpts, arg0, arg1)
}

// Powers is a free data retrieval call binding the contract method 0x850f43dd.
//
// Solidity: function powers(address , uint256 ) view returns(uint256)
func (_CosmosBridge *CosmosBridgeCallerSession) Powers(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _CosmosBridge.Contract.Powers(&_CosmosBridge.CallOpts, arg0, arg1)
}

// ProphecyClaims is a free data retrieval call binding the contract method 0xdb4237af.
//
// Solidity: function prophecyClaims(uint256 ) view returns(address ethereumReceiver, string symbol, uint256 amount)
func (_CosmosBridge *CosmosBridgeCaller) ProphecyClaims(opts *bind.CallOpts, arg0 *big.Int) (struct {
	EthereumReceiver common.Address
	Symbol           string
	Amount           *big.Int
}, error) {
	ret := new(struct {
		EthereumReceiver common.Address
		Symbol           string
		Amount           *big.Int
	})
	out := ret
	err := _CosmosBridge.contract.Call(opts, out, "prophecyClaims", arg0)
	return *ret, err
}

// ProphecyClaims is a free data retrieval call binding the contract method 0xdb4237af.
//
// Solidity: function prophecyClaims(uint256 ) view returns(address ethereumReceiver, string symbol, uint256 amount)
func (_CosmosBridge *CosmosBridgeSession) ProphecyClaims(arg0 *big.Int) (struct {
	EthereumReceiver common.Address
	Symbol           string
	Amount           *big.Int
}, error) {
	return _CosmosBridge.Contract.ProphecyClaims(&_CosmosBridge.CallOpts, arg0)
}

// ProphecyClaims is a free data retrieval call binding the contract method 0xdb4237af.
//
// Solidity: function prophecyClaims(uint256 ) view returns(address ethereumReceiver, string symbol, uint256 amount)
func (_CosmosBridge *CosmosBridgeCallerSession) ProphecyClaims(arg0 *big.Int) (struct {
	EthereumReceiver common.Address
	Symbol           string
	Amount           *big.Int
}, error) {
	return _CosmosBridge.Contract.ProphecyClaims(&_CosmosBridge.CallOpts, arg0)
}

// ProphecyRedeemed is a free data retrieval call binding the contract method 0x04e12aa5.
//
// Solidity: function prophecyRedeemed(uint256 ) view returns(bool)
func (_CosmosBridge *CosmosBridgeCaller) ProphecyRedeemed(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "prophecyRedeemed", arg0)
	return *ret0, err
}

// ProphecyRedeemed is a free data retrieval call binding the contract method 0x04e12aa5.
//
// Solidity: function prophecyRedeemed(uint256 ) view returns(bool)
func (_CosmosBridge *CosmosBridgeSession) ProphecyRedeemed(arg0 *big.Int) (bool, error) {
	return _CosmosBridge.Contract.ProphecyRedeemed(&_CosmosBridge.CallOpts, arg0)
}

// ProphecyRedeemed is a free data retrieval call binding the contract method 0x04e12aa5.
//
// Solidity: function prophecyRedeemed(uint256 ) view returns(bool)
func (_CosmosBridge *CosmosBridgeCallerSession) ProphecyRedeemed(arg0 *big.Int) (bool, error) {
	return _CosmosBridge.Contract.ProphecyRedeemed(&_CosmosBridge.CallOpts, arg0)
}

// SourceAddressToDestinationAddress is a free data retrieval call binding the contract method 0x7b010263.
//
// Solidity: function sourceAddressToDestinationAddress(address ) view returns(address)
func (_CosmosBridge *CosmosBridgeCaller) SourceAddressToDestinationAddress(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "sourceAddressToDestinationAddress", arg0)
	return *ret0, err
}

// SourceAddressToDestinationAddress is a free data retrieval call binding the contract method 0x7b010263.
//
// Solidity: function sourceAddressToDestinationAddress(address ) view returns(address)
func (_CosmosBridge *CosmosBridgeSession) SourceAddressToDestinationAddress(arg0 common.Address) (common.Address, error) {
	return _CosmosBridge.Contract.SourceAddressToDestinationAddress(&_CosmosBridge.CallOpts, arg0)
}

// SourceAddressToDestinationAddress is a free data retrieval call binding the contract method 0x7b010263.
//
// Solidity: function sourceAddressToDestinationAddress(address ) view returns(address)
func (_CosmosBridge *CosmosBridgeCallerSession) SourceAddressToDestinationAddress(arg0 common.Address) (common.Address, error) {
	return _CosmosBridge.Contract.SourceAddressToDestinationAddress(&_CosmosBridge.CallOpts, arg0)
}

// TotalPower is a free data retrieval call binding the contract method 0xdb3ad22c.
//
// Solidity: function totalPower() view returns(uint256)
func (_CosmosBridge *CosmosBridgeCaller) TotalPower(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "totalPower")
	return *ret0, err
}

// TotalPower is a free data retrieval call binding the contract method 0xdb3ad22c.
//
// Solidity: function totalPower() view returns(uint256)
func (_CosmosBridge *CosmosBridgeSession) TotalPower() (*big.Int, error) {
	return _CosmosBridge.Contract.TotalPower(&_CosmosBridge.CallOpts)
}

// TotalPower is a free data retrieval call binding the contract method 0xdb3ad22c.
//
// Solidity: function totalPower() view returns(uint256)
func (_CosmosBridge *CosmosBridgeCallerSession) TotalPower() (*big.Int, error) {
	return _CosmosBridge.Contract.TotalPower(&_CosmosBridge.CallOpts)
}

// ValidatorCount is a free data retrieval call binding the contract method 0x0f43a677.
//
// Solidity: function validatorCount() view returns(uint256)
func (_CosmosBridge *CosmosBridgeCaller) ValidatorCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "validatorCount")
	return *ret0, err
}

// ValidatorCount is a free data retrieval call binding the contract method 0x0f43a677.
//
// Solidity: function validatorCount() view returns(uint256)
func (_CosmosBridge *CosmosBridgeSession) ValidatorCount() (*big.Int, error) {
	return _CosmosBridge.Contract.ValidatorCount(&_CosmosBridge.CallOpts)
}

// ValidatorCount is a free data retrieval call binding the contract method 0x0f43a677.
//
// Solidity: function validatorCount() view returns(uint256)
func (_CosmosBridge *CosmosBridgeCallerSession) ValidatorCount() (*big.Int, error) {
	return _CosmosBridge.Contract.ValidatorCount(&_CosmosBridge.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0x45aaf18c.
//
// Solidity: function validators(address , uint256 ) view returns(bool)
func (_CosmosBridge *CosmosBridgeCaller) Validators(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "validators", arg0, arg1)
	return *ret0, err
}

// Validators is a free data retrieval call binding the contract method 0x45aaf18c.
//
// Solidity: function validators(address , uint256 ) view returns(bool)
func (_CosmosBridge *CosmosBridgeSession) Validators(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _CosmosBridge.Contract.Validators(&_CosmosBridge.CallOpts, arg0, arg1)
}

// Validators is a free data retrieval call binding the contract method 0x45aaf18c.
//
// Solidity: function validators(address , uint256 ) view returns(bool)
func (_CosmosBridge *CosmosBridgeCallerSession) Validators(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _CosmosBridge.Contract.Validators(&_CosmosBridge.CallOpts, arg0, arg1)
}

// Valset is a free data retrieval call binding the contract method 0x7f54af0c.
//
// Solidity: function valset() view returns(address)
func (_CosmosBridge *CosmosBridgeCaller) Valset(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _CosmosBridge.contract.Call(opts, out, "valset")
	return *ret0, err
}

// Valset is a free data retrieval call binding the contract method 0x7f54af0c.
//
// Solidity: function valset() view returns(address)
func (_CosmosBridge *CosmosBridgeSession) Valset() (common.Address, error) {
	return _CosmosBridge.Contract.Valset(&_CosmosBridge.CallOpts)
}

// Valset is a free data retrieval call binding the contract method 0x7f54af0c.
//
// Solidity: function valset() view returns(address)
func (_CosmosBridge *CosmosBridgeCallerSession) Valset() (common.Address, error) {
	return _CosmosBridge.Contract.Valset(&_CosmosBridge.CallOpts)
}

// AddValidator is a paid mutator transaction binding the contract method 0xfc6c1f02.
//
// Solidity: function addValidator(address _validatorAddress, uint256 _validatorPower) returns()
func (_CosmosBridge *CosmosBridgeTransactor) AddValidator(opts *bind.TransactOpts, _validatorAddress common.Address, _validatorPower *big.Int) (*types.Transaction, error) {
	return _CosmosBridge.contract.Transact(opts, "addValidator", _validatorAddress, _validatorPower)
}

// AddValidator is a paid mutator transaction binding the contract method 0xfc6c1f02.
//
// Solidity: function addValidator(address _validatorAddress, uint256 _validatorPower) returns()
func (_CosmosBridge *CosmosBridgeSession) AddValidator(_validatorAddress common.Address, _validatorPower *big.Int) (*types.Transaction, error) {
	return _CosmosBridge.Contract.AddValidator(&_CosmosBridge.TransactOpts, _validatorAddress, _validatorPower)
}

// AddValidator is a paid mutator transaction binding the contract method 0xfc6c1f02.
//
// Solidity: function addValidator(address _validatorAddress, uint256 _validatorPower) returns()
func (_CosmosBridge *CosmosBridgeTransactorSession) AddValidator(_validatorAddress common.Address, _validatorPower *big.Int) (*types.Transaction, error) {
	return _CosmosBridge.Contract.AddValidator(&_CosmosBridge.TransactOpts, _validatorAddress, _validatorPower)
}

// BatchSubmitProphecyClaimAggregatedSigs is a paid mutator transaction binding the contract method 0xfb5c9a58.
//
// Solidity: function batchSubmitProphecyClaimAggregatedSigs(bytes32[] sigs, (bytes,uint256,address,address,uint256,bool,uint128)[] claims, (address,uint8,bytes32,bytes32)[][] signatureData) returns()
func (_CosmosBridge *CosmosBridgeTransactor) BatchSubmitProphecyClaimAggregatedSigs(opts *bind.TransactOpts, sigs [][32]byte, claims []CosmosBridgeClaimData, signatureData [][]CosmosBridgeSignatureData) (*types.Transaction, error) {
	return _CosmosBridge.contract.Transact(opts, "batchSubmitProphecyClaimAggregatedSigs", sigs, claims, signatureData)
}

// BatchSubmitProphecyClaimAggregatedSigs is a paid mutator transaction binding the contract method 0xfb5c9a58.
//
// Solidity: function batchSubmitProphecyClaimAggregatedSigs(bytes32[] sigs, (bytes,uint256,address,address,uint256,bool,uint128)[] claims, (address,uint8,bytes32,bytes32)[][] signatureData) returns()
func (_CosmosBridge *CosmosBridgeSession) BatchSubmitProphecyClaimAggregatedSigs(sigs [][32]byte, claims []CosmosBridgeClaimData, signatureData [][]CosmosBridgeSignatureData) (*types.Transaction, error) {
	return _CosmosBridge.Contract.BatchSubmitProphecyClaimAggregatedSigs(&_CosmosBridge.TransactOpts, sigs, claims, signatureData)
}

// BatchSubmitProphecyClaimAggregatedSigs is a paid mutator transaction binding the contract method 0xfb5c9a58.
//
// Solidity: function batchSubmitProphecyClaimAggregatedSigs(bytes32[] sigs, (bytes,uint256,address,address,uint256,bool,uint128)[] claims, (address,uint8,bytes32,bytes32)[][] signatureData) returns()
func (_CosmosBridge *CosmosBridgeTransactorSession) BatchSubmitProphecyClaimAggregatedSigs(sigs [][32]byte, claims []CosmosBridgeClaimData, signatureData [][]CosmosBridgeSignatureData) (*types.Transaction, error) {
	return _CosmosBridge.Contract.BatchSubmitProphecyClaimAggregatedSigs(&_CosmosBridge.TransactOpts, sigs, claims, signatureData)
}

// ChangeOperator is a paid mutator transaction binding the contract method 0x06394c9b.
//
// Solidity: function changeOperator(address _newOperator) returns()
func (_CosmosBridge *CosmosBridgeTransactor) ChangeOperator(opts *bind.TransactOpts, _newOperator common.Address) (*types.Transaction, error) {
	return _CosmosBridge.contract.Transact(opts, "changeOperator", _newOperator)
}

// ChangeOperator is a paid mutator transaction binding the contract method 0x06394c9b.
//
// Solidity: function changeOperator(address _newOperator) returns()
func (_CosmosBridge *CosmosBridgeSession) ChangeOperator(_newOperator common.Address) (*types.Transaction, error) {
	return _CosmosBridge.Contract.ChangeOperator(&_CosmosBridge.TransactOpts, _newOperator)
}

// ChangeOperator is a paid mutator transaction binding the contract method 0x06394c9b.
//
// Solidity: function changeOperator(address _newOperator) returns()
func (_CosmosBridge *CosmosBridgeTransactorSession) ChangeOperator(_newOperator common.Address) (*types.Transaction, error) {
	return _CosmosBridge.Contract.ChangeOperator(&_CosmosBridge.TransactOpts, _newOperator)
}

// CreateNewBridgeToken is a paid mutator transaction binding the contract method 0x2bdb26ee.
//
// Solidity: function createNewBridgeToken(string symbol, string name, address sourceChainTokenAddress, uint8 decimals, uint256 chainDescriptor) returns()
func (_CosmosBridge *CosmosBridgeTransactor) CreateNewBridgeToken(opts *bind.TransactOpts, symbol string, name string, sourceChainTokenAddress common.Address, decimals uint8, chainDescriptor *big.Int) (*types.Transaction, error) {
	return _CosmosBridge.contract.Transact(opts, "createNewBridgeToken", symbol, name, sourceChainTokenAddress, decimals, chainDescriptor)
}

// CreateNewBridgeToken is a paid mutator transaction binding the contract method 0x2bdb26ee.
//
// Solidity: function createNewBridgeToken(string symbol, string name, address sourceChainTokenAddress, uint8 decimals, uint256 chainDescriptor) returns()
func (_CosmosBridge *CosmosBridgeSession) CreateNewBridgeToken(symbol string, name string, sourceChainTokenAddress common.Address, decimals uint8, chainDescriptor *big.Int) (*types.Transaction, error) {
	return _CosmosBridge.Contract.CreateNewBridgeToken(&_CosmosBridge.TransactOpts, symbol, name, sourceChainTokenAddress, decimals, chainDescriptor)
}

// CreateNewBridgeToken is a paid mutator transaction binding the contract method 0x2bdb26ee.
//
// Solidity: function createNewBridgeToken(string symbol, string name, address sourceChainTokenAddress, uint8 decimals, uint256 chainDescriptor) returns()
func (_CosmosBridge *CosmosBridgeTransactorSession) CreateNewBridgeToken(symbol string, name string, sourceChainTokenAddress common.Address, decimals uint8, chainDescriptor *big.Int) (*types.Transaction, error) {
	return _CosmosBridge.Contract.CreateNewBridgeToken(&_CosmosBridge.TransactOpts, symbol, name, sourceChainTokenAddress, decimals, chainDescriptor)
}

// Initialize is a paid mutator transaction binding the contract method 0x9a19224f.
//
// Solidity: function initialize(address _operator, uint256 _consensusThreshold, address[] _initValidators, uint256[] _initPowers) returns()
func (_CosmosBridge *CosmosBridgeTransactor) Initialize(opts *bind.TransactOpts, _operator common.Address, _consensusThreshold *big.Int, _initValidators []common.Address, _initPowers []*big.Int) (*types.Transaction, error) {
	return _CosmosBridge.contract.Transact(opts, "initialize", _operator, _consensusThreshold, _initValidators, _initPowers)
}

// Initialize is a paid mutator transaction binding the contract method 0x9a19224f.
//
// Solidity: function initialize(address _operator, uint256 _consensusThreshold, address[] _initValidators, uint256[] _initPowers) returns()
func (_CosmosBridge *CosmosBridgeSession) Initialize(_operator common.Address, _consensusThreshold *big.Int, _initValidators []common.Address, _initPowers []*big.Int) (*types.Transaction, error) {
	return _CosmosBridge.Contract.Initialize(&_CosmosBridge.TransactOpts, _operator, _consensusThreshold, _initValidators, _initPowers)
}

// Initialize is a paid mutator transaction binding the contract method 0x9a19224f.
//
// Solidity: function initialize(address _operator, uint256 _consensusThreshold, address[] _initValidators, uint256[] _initPowers) returns()
func (_CosmosBridge *CosmosBridgeTransactorSession) Initialize(_operator common.Address, _consensusThreshold *big.Int, _initValidators []common.Address, _initPowers []*big.Int) (*types.Transaction, error) {
	return _CosmosBridge.Contract.Initialize(&_CosmosBridge.TransactOpts, _operator, _consensusThreshold, _initValidators, _initPowers)
}

// RecoverGas is a paid mutator transaction binding the contract method 0xb5672be3.
//
// Solidity: function recoverGas(uint256 _valsetVersion, address _validatorAddress) returns()
func (_CosmosBridge *CosmosBridgeTransactor) RecoverGas(opts *bind.TransactOpts, _valsetVersion *big.Int, _validatorAddress common.Address) (*types.Transaction, error) {
	return _CosmosBridge.contract.Transact(opts, "recoverGas", _valsetVersion, _validatorAddress)
}

// RecoverGas is a paid mutator transaction binding the contract method 0xb5672be3.
//
// Solidity: function recoverGas(uint256 _valsetVersion, address _validatorAddress) returns()
func (_CosmosBridge *CosmosBridgeSession) RecoverGas(_valsetVersion *big.Int, _validatorAddress common.Address) (*types.Transaction, error) {
	return _CosmosBridge.Contract.RecoverGas(&_CosmosBridge.TransactOpts, _valsetVersion, _validatorAddress)
}

// RecoverGas is a paid mutator transaction binding the contract method 0xb5672be3.
//
// Solidity: function recoverGas(uint256 _valsetVersion, address _validatorAddress) returns()
func (_CosmosBridge *CosmosBridgeTransactorSession) RecoverGas(_valsetVersion *big.Int, _validatorAddress common.Address) (*types.Transaction, error) {
	return _CosmosBridge.Contract.RecoverGas(&_CosmosBridge.TransactOpts, _valsetVersion, _validatorAddress)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(address _validatorAddress) returns()
func (_CosmosBridge *CosmosBridgeTransactor) RemoveValidator(opts *bind.TransactOpts, _validatorAddress common.Address) (*types.Transaction, error) {
	return _CosmosBridge.contract.Transact(opts, "removeValidator", _validatorAddress)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(address _validatorAddress) returns()
func (_CosmosBridge *CosmosBridgeSession) RemoveValidator(_validatorAddress common.Address) (*types.Transaction, error) {
	return _CosmosBridge.Contract.RemoveValidator(&_CosmosBridge.TransactOpts, _validatorAddress)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(address _validatorAddress) returns()
func (_CosmosBridge *CosmosBridgeTransactorSession) RemoveValidator(_validatorAddress common.Address) (*types.Transaction, error) {
	return _CosmosBridge.Contract.RemoveValidator(&_CosmosBridge.TransactOpts, _validatorAddress)
}

// SetBridgeBank is a paid mutator transaction binding the contract method 0x814c92c3.
//
// Solidity: function setBridgeBank(address _bridgeBank) returns()
func (_CosmosBridge *CosmosBridgeTransactor) SetBridgeBank(opts *bind.TransactOpts, _bridgeBank common.Address) (*types.Transaction, error) {
	return _CosmosBridge.contract.Transact(opts, "setBridgeBank", _bridgeBank)
}

// SetBridgeBank is a paid mutator transaction binding the contract method 0x814c92c3.
//
// Solidity: function setBridgeBank(address _bridgeBank) returns()
func (_CosmosBridge *CosmosBridgeSession) SetBridgeBank(_bridgeBank common.Address) (*types.Transaction, error) {
	return _CosmosBridge.Contract.SetBridgeBank(&_CosmosBridge.TransactOpts, _bridgeBank)
}

// SetBridgeBank is a paid mutator transaction binding the contract method 0x814c92c3.
//
// Solidity: function setBridgeBank(address _bridgeBank) returns()
func (_CosmosBridge *CosmosBridgeTransactorSession) SetBridgeBank(_bridgeBank common.Address) (*types.Transaction, error) {
	return _CosmosBridge.Contract.SetBridgeBank(&_CosmosBridge.TransactOpts, _bridgeBank)
}

// SubmitProphecyClaimAggregatedSigs is a paid mutator transaction binding the contract method 0x3b18ff9f.
//
// Solidity: function submitProphecyClaimAggregatedSigs(bytes32 hashDigest, (bytes,uint256,address,address,uint256,bool,uint128) claimData, (address,uint8,bytes32,bytes32)[] signatureData) returns()
func (_CosmosBridge *CosmosBridgeTransactor) SubmitProphecyClaimAggregatedSigs(opts *bind.TransactOpts, hashDigest [32]byte, claimData CosmosBridgeClaimData, signatureData []CosmosBridgeSignatureData) (*types.Transaction, error) {
	return _CosmosBridge.contract.Transact(opts, "submitProphecyClaimAggregatedSigs", hashDigest, claimData, signatureData)
}

// SubmitProphecyClaimAggregatedSigs is a paid mutator transaction binding the contract method 0x3b18ff9f.
//
// Solidity: function submitProphecyClaimAggregatedSigs(bytes32 hashDigest, (bytes,uint256,address,address,uint256,bool,uint128) claimData, (address,uint8,bytes32,bytes32)[] signatureData) returns()
func (_CosmosBridge *CosmosBridgeSession) SubmitProphecyClaimAggregatedSigs(hashDigest [32]byte, claimData CosmosBridgeClaimData, signatureData []CosmosBridgeSignatureData) (*types.Transaction, error) {
	return _CosmosBridge.Contract.SubmitProphecyClaimAggregatedSigs(&_CosmosBridge.TransactOpts, hashDigest, claimData, signatureData)
}

// SubmitProphecyClaimAggregatedSigs is a paid mutator transaction binding the contract method 0x3b18ff9f.
//
// Solidity: function submitProphecyClaimAggregatedSigs(bytes32 hashDigest, (bytes,uint256,address,address,uint256,bool,uint128) claimData, (address,uint8,bytes32,bytes32)[] signatureData) returns()
func (_CosmosBridge *CosmosBridgeTransactorSession) SubmitProphecyClaimAggregatedSigs(hashDigest [32]byte, claimData CosmosBridgeClaimData, signatureData []CosmosBridgeSignatureData) (*types.Transaction, error) {
	return _CosmosBridge.Contract.SubmitProphecyClaimAggregatedSigs(&_CosmosBridge.TransactOpts, hashDigest, claimData, signatureData)
}

// UpdateValidatorPower is a paid mutator transaction binding the contract method 0x2e75293b.
//
// Solidity: function updateValidatorPower(address _validatorAddress, uint256 _newValidatorPower) returns()
func (_CosmosBridge *CosmosBridgeTransactor) UpdateValidatorPower(opts *bind.TransactOpts, _validatorAddress common.Address, _newValidatorPower *big.Int) (*types.Transaction, error) {
	return _CosmosBridge.contract.Transact(opts, "updateValidatorPower", _validatorAddress, _newValidatorPower)
}

// UpdateValidatorPower is a paid mutator transaction binding the contract method 0x2e75293b.
//
// Solidity: function updateValidatorPower(address _validatorAddress, uint256 _newValidatorPower) returns()
func (_CosmosBridge *CosmosBridgeSession) UpdateValidatorPower(_validatorAddress common.Address, _newValidatorPower *big.Int) (*types.Transaction, error) {
	return _CosmosBridge.Contract.UpdateValidatorPower(&_CosmosBridge.TransactOpts, _validatorAddress, _newValidatorPower)
}

// UpdateValidatorPower is a paid mutator transaction binding the contract method 0x2e75293b.
//
// Solidity: function updateValidatorPower(address _validatorAddress, uint256 _newValidatorPower) returns()
func (_CosmosBridge *CosmosBridgeTransactorSession) UpdateValidatorPower(_validatorAddress common.Address, _newValidatorPower *big.Int) (*types.Transaction, error) {
	return _CosmosBridge.Contract.UpdateValidatorPower(&_CosmosBridge.TransactOpts, _validatorAddress, _newValidatorPower)
}

// UpdateValset is a paid mutator transaction binding the contract method 0x788cf92f.
//
// Solidity: function updateValset(address[] _validators, uint256[] _powers) returns()
func (_CosmosBridge *CosmosBridgeTransactor) UpdateValset(opts *bind.TransactOpts, _validators []common.Address, _powers []*big.Int) (*types.Transaction, error) {
	return _CosmosBridge.contract.Transact(opts, "updateValset", _validators, _powers)
}

// UpdateValset is a paid mutator transaction binding the contract method 0x788cf92f.
//
// Solidity: function updateValset(address[] _validators, uint256[] _powers) returns()
func (_CosmosBridge *CosmosBridgeSession) UpdateValset(_validators []common.Address, _powers []*big.Int) (*types.Transaction, error) {
	return _CosmosBridge.Contract.UpdateValset(&_CosmosBridge.TransactOpts, _validators, _powers)
}

// UpdateValset is a paid mutator transaction binding the contract method 0x788cf92f.
//
// Solidity: function updateValset(address[] _validators, uint256[] _powers) returns()
func (_CosmosBridge *CosmosBridgeTransactorSession) UpdateValset(_validators []common.Address, _powers []*big.Int) (*types.Transaction, error) {
	return _CosmosBridge.Contract.UpdateValset(&_CosmosBridge.TransactOpts, _validators, _powers)
}

// CosmosBridgeLogBridgeBankSetIterator is returned from FilterLogBridgeBankSet and is used to iterate over the raw logs and unpacked data for LogBridgeBankSet events raised by the CosmosBridge contract.
type CosmosBridgeLogBridgeBankSetIterator struct {
	Event *CosmosBridgeLogBridgeBankSet // Event containing the contract specifics and raw log

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
func (it *CosmosBridgeLogBridgeBankSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CosmosBridgeLogBridgeBankSet)
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
		it.Event = new(CosmosBridgeLogBridgeBankSet)
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
func (it *CosmosBridgeLogBridgeBankSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CosmosBridgeLogBridgeBankSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CosmosBridgeLogBridgeBankSet represents a LogBridgeBankSet event raised by the CosmosBridge contract.
type CosmosBridgeLogBridgeBankSet struct {
	BridgeBank common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogBridgeBankSet is a free log retrieval operation binding the contract event 0xc8b65043fb196ac032b79a435397d1d14a96b4e9d12e366c3b1f550cb01d2dfa.
//
// Solidity: event LogBridgeBankSet(address bridgeBank)
func (_CosmosBridge *CosmosBridgeFilterer) FilterLogBridgeBankSet(opts *bind.FilterOpts) (*CosmosBridgeLogBridgeBankSetIterator, error) {

	logs, sub, err := _CosmosBridge.contract.FilterLogs(opts, "LogBridgeBankSet")
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeLogBridgeBankSetIterator{contract: _CosmosBridge.contract, event: "LogBridgeBankSet", logs: logs, sub: sub}, nil
}

// WatchLogBridgeBankSet is a free log subscription operation binding the contract event 0xc8b65043fb196ac032b79a435397d1d14a96b4e9d12e366c3b1f550cb01d2dfa.
//
// Solidity: event LogBridgeBankSet(address bridgeBank)
func (_CosmosBridge *CosmosBridgeFilterer) WatchLogBridgeBankSet(opts *bind.WatchOpts, sink chan<- *CosmosBridgeLogBridgeBankSet) (event.Subscription, error) {

	logs, sub, err := _CosmosBridge.contract.WatchLogs(opts, "LogBridgeBankSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CosmosBridgeLogBridgeBankSet)
				if err := _CosmosBridge.contract.UnpackLog(event, "LogBridgeBankSet", log); err != nil {
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

// ParseLogBridgeBankSet is a log parse operation binding the contract event 0xc8b65043fb196ac032b79a435397d1d14a96b4e9d12e366c3b1f550cb01d2dfa.
//
// Solidity: event LogBridgeBankSet(address bridgeBank)
func (_CosmosBridge *CosmosBridgeFilterer) ParseLogBridgeBankSet(log types.Log) (*CosmosBridgeLogBridgeBankSet, error) {
	event := new(CosmosBridgeLogBridgeBankSet)
	if err := _CosmosBridge.contract.UnpackLog(event, "LogBridgeBankSet", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CosmosBridgeLogNewBridgeTokenCreatedIterator is returned from FilterLogNewBridgeTokenCreated and is used to iterate over the raw logs and unpacked data for LogNewBridgeTokenCreated events raised by the CosmosBridge contract.
type CosmosBridgeLogNewBridgeTokenCreatedIterator struct {
	Event *CosmosBridgeLogNewBridgeTokenCreated // Event containing the contract specifics and raw log

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
func (it *CosmosBridgeLogNewBridgeTokenCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CosmosBridgeLogNewBridgeTokenCreated)
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
		it.Event = new(CosmosBridgeLogNewBridgeTokenCreated)
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
func (it *CosmosBridgeLogNewBridgeTokenCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CosmosBridgeLogNewBridgeTokenCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CosmosBridgeLogNewBridgeTokenCreated represents a LogNewBridgeTokenCreated event raised by the CosmosBridge contract.
type CosmosBridgeLogNewBridgeTokenCreated struct {
	Decimals              uint8
	SourceChainDescriptor *big.Int
	Name                  string
	Symbol                string
	SourceContractAddress common.Address
	BridgeTokenAddress    common.Address
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterLogNewBridgeTokenCreated is a free log retrieval operation binding the contract event 0xa3866dbc9098b0c8ef4b4aa3dc7c0c5f86be05de8205b28bc2734ca9b530e321.
//
// Solidity: event LogNewBridgeTokenCreated(uint8 decimals, uint256 indexed sourceChainDescriptor, string name, string symbol, address indexed sourceContractAddress, address indexed bridgeTokenAddress)
func (_CosmosBridge *CosmosBridgeFilterer) FilterLogNewBridgeTokenCreated(opts *bind.FilterOpts, sourceChainDescriptor []*big.Int, sourceContractAddress []common.Address, bridgeTokenAddress []common.Address) (*CosmosBridgeLogNewBridgeTokenCreatedIterator, error) {

	var sourceChainDescriptorRule []interface{}
	for _, sourceChainDescriptorItem := range sourceChainDescriptor {
		sourceChainDescriptorRule = append(sourceChainDescriptorRule, sourceChainDescriptorItem)
	}

	var sourceContractAddressRule []interface{}
	for _, sourceContractAddressItem := range sourceContractAddress {
		sourceContractAddressRule = append(sourceContractAddressRule, sourceContractAddressItem)
	}
	var bridgeTokenAddressRule []interface{}
	for _, bridgeTokenAddressItem := range bridgeTokenAddress {
		bridgeTokenAddressRule = append(bridgeTokenAddressRule, bridgeTokenAddressItem)
	}

	logs, sub, err := _CosmosBridge.contract.FilterLogs(opts, "LogNewBridgeTokenCreated", sourceChainDescriptorRule, sourceContractAddressRule, bridgeTokenAddressRule)
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeLogNewBridgeTokenCreatedIterator{contract: _CosmosBridge.contract, event: "LogNewBridgeTokenCreated", logs: logs, sub: sub}, nil
}

// WatchLogNewBridgeTokenCreated is a free log subscription operation binding the contract event 0xa3866dbc9098b0c8ef4b4aa3dc7c0c5f86be05de8205b28bc2734ca9b530e321.
//
// Solidity: event LogNewBridgeTokenCreated(uint8 decimals, uint256 indexed sourceChainDescriptor, string name, string symbol, address indexed sourceContractAddress, address indexed bridgeTokenAddress)
func (_CosmosBridge *CosmosBridgeFilterer) WatchLogNewBridgeTokenCreated(opts *bind.WatchOpts, sink chan<- *CosmosBridgeLogNewBridgeTokenCreated, sourceChainDescriptor []*big.Int, sourceContractAddress []common.Address, bridgeTokenAddress []common.Address) (event.Subscription, error) {

	var sourceChainDescriptorRule []interface{}
	for _, sourceChainDescriptorItem := range sourceChainDescriptor {
		sourceChainDescriptorRule = append(sourceChainDescriptorRule, sourceChainDescriptorItem)
	}

	var sourceContractAddressRule []interface{}
	for _, sourceContractAddressItem := range sourceContractAddress {
		sourceContractAddressRule = append(sourceContractAddressRule, sourceContractAddressItem)
	}
	var bridgeTokenAddressRule []interface{}
	for _, bridgeTokenAddressItem := range bridgeTokenAddress {
		bridgeTokenAddressRule = append(bridgeTokenAddressRule, bridgeTokenAddressItem)
	}

	logs, sub, err := _CosmosBridge.contract.WatchLogs(opts, "LogNewBridgeTokenCreated", sourceChainDescriptorRule, sourceContractAddressRule, bridgeTokenAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CosmosBridgeLogNewBridgeTokenCreated)
				if err := _CosmosBridge.contract.UnpackLog(event, "LogNewBridgeTokenCreated", log); err != nil {
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

// ParseLogNewBridgeTokenCreated is a log parse operation binding the contract event 0xa3866dbc9098b0c8ef4b4aa3dc7c0c5f86be05de8205b28bc2734ca9b530e321.
//
// Solidity: event LogNewBridgeTokenCreated(uint8 decimals, uint256 indexed sourceChainDescriptor, string name, string symbol, address indexed sourceContractAddress, address indexed bridgeTokenAddress)
func (_CosmosBridge *CosmosBridgeFilterer) ParseLogNewBridgeTokenCreated(log types.Log) (*CosmosBridgeLogNewBridgeTokenCreated, error) {
	event := new(CosmosBridgeLogNewBridgeTokenCreated)
	if err := _CosmosBridge.contract.UnpackLog(event, "LogNewBridgeTokenCreated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CosmosBridgeLogNewOracleClaimIterator is returned from FilterLogNewOracleClaim and is used to iterate over the raw logs and unpacked data for LogNewOracleClaim events raised by the CosmosBridge contract.
type CosmosBridgeLogNewOracleClaimIterator struct {
	Event *CosmosBridgeLogNewOracleClaim // Event containing the contract specifics and raw log

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
func (it *CosmosBridgeLogNewOracleClaimIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CosmosBridgeLogNewOracleClaim)
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
		it.Event = new(CosmosBridgeLogNewOracleClaim)
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
func (it *CosmosBridgeLogNewOracleClaimIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CosmosBridgeLogNewOracleClaimIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CosmosBridgeLogNewOracleClaim represents a LogNewOracleClaim event raised by the CosmosBridge contract.
type CosmosBridgeLogNewOracleClaim struct {
	ProphecyID       *big.Int
	ValidatorAddress common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogNewOracleClaim is a free log retrieval operation binding the contract event 0x668fce9833323940537a9000d512a6c580a1c0797d2b526db0078ee9c5a087a9.
//
// Solidity: event LogNewOracleClaim(uint256 _prophecyID, address _validatorAddress)
func (_CosmosBridge *CosmosBridgeFilterer) FilterLogNewOracleClaim(opts *bind.FilterOpts) (*CosmosBridgeLogNewOracleClaimIterator, error) {

	logs, sub, err := _CosmosBridge.contract.FilterLogs(opts, "LogNewOracleClaim")
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeLogNewOracleClaimIterator{contract: _CosmosBridge.contract, event: "LogNewOracleClaim", logs: logs, sub: sub}, nil
}

// WatchLogNewOracleClaim is a free log subscription operation binding the contract event 0x668fce9833323940537a9000d512a6c580a1c0797d2b526db0078ee9c5a087a9.
//
// Solidity: event LogNewOracleClaim(uint256 _prophecyID, address _validatorAddress)
func (_CosmosBridge *CosmosBridgeFilterer) WatchLogNewOracleClaim(opts *bind.WatchOpts, sink chan<- *CosmosBridgeLogNewOracleClaim) (event.Subscription, error) {

	logs, sub, err := _CosmosBridge.contract.WatchLogs(opts, "LogNewOracleClaim")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CosmosBridgeLogNewOracleClaim)
				if err := _CosmosBridge.contract.UnpackLog(event, "LogNewOracleClaim", log); err != nil {
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
func (_CosmosBridge *CosmosBridgeFilterer) ParseLogNewOracleClaim(log types.Log) (*CosmosBridgeLogNewOracleClaim, error) {
	event := new(CosmosBridgeLogNewOracleClaim)
	if err := _CosmosBridge.contract.UnpackLog(event, "LogNewOracleClaim", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CosmosBridgeLogNewProphecyClaimIterator is returned from FilterLogNewProphecyClaim and is used to iterate over the raw logs and unpacked data for LogNewProphecyClaim events raised by the CosmosBridge contract.
type CosmosBridgeLogNewProphecyClaimIterator struct {
	Event *CosmosBridgeLogNewProphecyClaim // Event containing the contract specifics and raw log

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
func (it *CosmosBridgeLogNewProphecyClaimIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CosmosBridgeLogNewProphecyClaim)
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
		it.Event = new(CosmosBridgeLogNewProphecyClaim)
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
func (it *CosmosBridgeLogNewProphecyClaimIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CosmosBridgeLogNewProphecyClaimIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CosmosBridgeLogNewProphecyClaim represents a LogNewProphecyClaim event raised by the CosmosBridge contract.
type CosmosBridgeLogNewProphecyClaim struct {
	ProphecyID       *big.Int
	EthereumReceiver common.Address
	Amount           *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogNewProphecyClaim is a free log retrieval operation binding the contract event 0x392e2c08a6c5b92e1d4775a1e43829652f156f1310bd3fe2935f586e2d4ce36e.
//
// Solidity: event LogNewProphecyClaim(uint256 indexed prophecyID, address indexed ethereumReceiver, uint256 indexed amount)
func (_CosmosBridge *CosmosBridgeFilterer) FilterLogNewProphecyClaim(opts *bind.FilterOpts, prophecyID []*big.Int, ethereumReceiver []common.Address, amount []*big.Int) (*CosmosBridgeLogNewProphecyClaimIterator, error) {

	var prophecyIDRule []interface{}
	for _, prophecyIDItem := range prophecyID {
		prophecyIDRule = append(prophecyIDRule, prophecyIDItem)
	}
	var ethereumReceiverRule []interface{}
	for _, ethereumReceiverItem := range ethereumReceiver {
		ethereumReceiverRule = append(ethereumReceiverRule, ethereumReceiverItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _CosmosBridge.contract.FilterLogs(opts, "LogNewProphecyClaim", prophecyIDRule, ethereumReceiverRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeLogNewProphecyClaimIterator{contract: _CosmosBridge.contract, event: "LogNewProphecyClaim", logs: logs, sub: sub}, nil
}

// WatchLogNewProphecyClaim is a free log subscription operation binding the contract event 0x392e2c08a6c5b92e1d4775a1e43829652f156f1310bd3fe2935f586e2d4ce36e.
//
// Solidity: event LogNewProphecyClaim(uint256 indexed prophecyID, address indexed ethereumReceiver, uint256 indexed amount)
func (_CosmosBridge *CosmosBridgeFilterer) WatchLogNewProphecyClaim(opts *bind.WatchOpts, sink chan<- *CosmosBridgeLogNewProphecyClaim, prophecyID []*big.Int, ethereumReceiver []common.Address, amount []*big.Int) (event.Subscription, error) {

	var prophecyIDRule []interface{}
	for _, prophecyIDItem := range prophecyID {
		prophecyIDRule = append(prophecyIDRule, prophecyIDItem)
	}
	var ethereumReceiverRule []interface{}
	for _, ethereumReceiverItem := range ethereumReceiver {
		ethereumReceiverRule = append(ethereumReceiverRule, ethereumReceiverItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _CosmosBridge.contract.WatchLogs(opts, "LogNewProphecyClaim", prophecyIDRule, ethereumReceiverRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CosmosBridgeLogNewProphecyClaim)
				if err := _CosmosBridge.contract.UnpackLog(event, "LogNewProphecyClaim", log); err != nil {
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

// ParseLogNewProphecyClaim is a log parse operation binding the contract event 0x392e2c08a6c5b92e1d4775a1e43829652f156f1310bd3fe2935f586e2d4ce36e.
//
// Solidity: event LogNewProphecyClaim(uint256 indexed prophecyID, address indexed ethereumReceiver, uint256 indexed amount)
func (_CosmosBridge *CosmosBridgeFilterer) ParseLogNewProphecyClaim(log types.Log) (*CosmosBridgeLogNewProphecyClaim, error) {
	event := new(CosmosBridgeLogNewProphecyClaim)
	if err := _CosmosBridge.contract.UnpackLog(event, "LogNewProphecyClaim", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CosmosBridgeLogProphecyCompletedIterator is returned from FilterLogProphecyCompleted and is used to iterate over the raw logs and unpacked data for LogProphecyCompleted events raised by the CosmosBridge contract.
type CosmosBridgeLogProphecyCompletedIterator struct {
	Event *CosmosBridgeLogProphecyCompleted // Event containing the contract specifics and raw log

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
func (it *CosmosBridgeLogProphecyCompletedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CosmosBridgeLogProphecyCompleted)
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
		it.Event = new(CosmosBridgeLogProphecyCompleted)
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
func (it *CosmosBridgeLogProphecyCompletedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CosmosBridgeLogProphecyCompletedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CosmosBridgeLogProphecyCompleted represents a LogProphecyCompleted event raised by the CosmosBridge contract.
type CosmosBridgeLogProphecyCompleted struct {
	ProphecyID *big.Int
	Success    bool
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterLogProphecyCompleted is a free log retrieval operation binding the contract event 0x6cf2aa1395dae0a852879973e322e6d1f5a48485a3a41c28a2fd8f82d3aae487.
//
// Solidity: event LogProphecyCompleted(uint256 prophecyID, bool success)
func (_CosmosBridge *CosmosBridgeFilterer) FilterLogProphecyCompleted(opts *bind.FilterOpts) (*CosmosBridgeLogProphecyCompletedIterator, error) {

	logs, sub, err := _CosmosBridge.contract.FilterLogs(opts, "LogProphecyCompleted")
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeLogProphecyCompletedIterator{contract: _CosmosBridge.contract, event: "LogProphecyCompleted", logs: logs, sub: sub}, nil
}

// WatchLogProphecyCompleted is a free log subscription operation binding the contract event 0x6cf2aa1395dae0a852879973e322e6d1f5a48485a3a41c28a2fd8f82d3aae487.
//
// Solidity: event LogProphecyCompleted(uint256 prophecyID, bool success)
func (_CosmosBridge *CosmosBridgeFilterer) WatchLogProphecyCompleted(opts *bind.WatchOpts, sink chan<- *CosmosBridgeLogProphecyCompleted) (event.Subscription, error) {

	logs, sub, err := _CosmosBridge.contract.WatchLogs(opts, "LogProphecyCompleted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CosmosBridgeLogProphecyCompleted)
				if err := _CosmosBridge.contract.UnpackLog(event, "LogProphecyCompleted", log); err != nil {
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

// ParseLogProphecyCompleted is a log parse operation binding the contract event 0x6cf2aa1395dae0a852879973e322e6d1f5a48485a3a41c28a2fd8f82d3aae487.
//
// Solidity: event LogProphecyCompleted(uint256 prophecyID, bool success)
func (_CosmosBridge *CosmosBridgeFilterer) ParseLogProphecyCompleted(log types.Log) (*CosmosBridgeLogProphecyCompleted, error) {
	event := new(CosmosBridgeLogProphecyCompleted)
	if err := _CosmosBridge.contract.UnpackLog(event, "LogProphecyCompleted", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CosmosBridgeLogProphecyProcessedIterator is returned from FilterLogProphecyProcessed and is used to iterate over the raw logs and unpacked data for LogProphecyProcessed events raised by the CosmosBridge contract.
type CosmosBridgeLogProphecyProcessedIterator struct {
	Event *CosmosBridgeLogProphecyProcessed // Event containing the contract specifics and raw log

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
func (it *CosmosBridgeLogProphecyProcessedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CosmosBridgeLogProphecyProcessed)
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
		it.Event = new(CosmosBridgeLogProphecyProcessed)
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
func (it *CosmosBridgeLogProphecyProcessedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CosmosBridgeLogProphecyProcessedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CosmosBridgeLogProphecyProcessed represents a LogProphecyProcessed event raised by the CosmosBridge contract.
type CosmosBridgeLogProphecyProcessed struct {
	ProphecyID             *big.Int
	ProphecyPowerCurrent   *big.Int
	ProphecyPowerThreshold *big.Int
	Submitter              common.Address
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterLogProphecyProcessed is a free log retrieval operation binding the contract event 0x1d8e3fbd601d9d92db7022fb97f75e132841b94db732dcecb0c93cb31852fcbc.
//
// Solidity: event LogProphecyProcessed(uint256 _prophecyID, uint256 _prophecyPowerCurrent, uint256 _prophecyPowerThreshold, address _submitter)
func (_CosmosBridge *CosmosBridgeFilterer) FilterLogProphecyProcessed(opts *bind.FilterOpts) (*CosmosBridgeLogProphecyProcessedIterator, error) {

	logs, sub, err := _CosmosBridge.contract.FilterLogs(opts, "LogProphecyProcessed")
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeLogProphecyProcessedIterator{contract: _CosmosBridge.contract, event: "LogProphecyProcessed", logs: logs, sub: sub}, nil
}

// WatchLogProphecyProcessed is a free log subscription operation binding the contract event 0x1d8e3fbd601d9d92db7022fb97f75e132841b94db732dcecb0c93cb31852fcbc.
//
// Solidity: event LogProphecyProcessed(uint256 _prophecyID, uint256 _prophecyPowerCurrent, uint256 _prophecyPowerThreshold, address _submitter)
func (_CosmosBridge *CosmosBridgeFilterer) WatchLogProphecyProcessed(opts *bind.WatchOpts, sink chan<- *CosmosBridgeLogProphecyProcessed) (event.Subscription, error) {

	logs, sub, err := _CosmosBridge.contract.WatchLogs(opts, "LogProphecyProcessed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CosmosBridgeLogProphecyProcessed)
				if err := _CosmosBridge.contract.UnpackLog(event, "LogProphecyProcessed", log); err != nil {
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
func (_CosmosBridge *CosmosBridgeFilterer) ParseLogProphecyProcessed(log types.Log) (*CosmosBridgeLogProphecyProcessed, error) {
	event := new(CosmosBridgeLogProphecyProcessed)
	if err := _CosmosBridge.contract.UnpackLog(event, "LogProphecyProcessed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CosmosBridgeLogValidatorAddedIterator is returned from FilterLogValidatorAdded and is used to iterate over the raw logs and unpacked data for LogValidatorAdded events raised by the CosmosBridge contract.
type CosmosBridgeLogValidatorAddedIterator struct {
	Event *CosmosBridgeLogValidatorAdded // Event containing the contract specifics and raw log

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
func (it *CosmosBridgeLogValidatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CosmosBridgeLogValidatorAdded)
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
		it.Event = new(CosmosBridgeLogValidatorAdded)
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
func (it *CosmosBridgeLogValidatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CosmosBridgeLogValidatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CosmosBridgeLogValidatorAdded represents a LogValidatorAdded event raised by the CosmosBridge contract.
type CosmosBridgeLogValidatorAdded struct {
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
func (_CosmosBridge *CosmosBridgeFilterer) FilterLogValidatorAdded(opts *bind.FilterOpts) (*CosmosBridgeLogValidatorAddedIterator, error) {

	logs, sub, err := _CosmosBridge.contract.FilterLogs(opts, "LogValidatorAdded")
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeLogValidatorAddedIterator{contract: _CosmosBridge.contract, event: "LogValidatorAdded", logs: logs, sub: sub}, nil
}

// WatchLogValidatorAdded is a free log subscription operation binding the contract event 0x1a396bcf647502e902dce665d58a0c1b25f982f193ab9a1d0f1500d8d927bf2a.
//
// Solidity: event LogValidatorAdded(address _validator, uint256 _power, uint256 _currentValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_CosmosBridge *CosmosBridgeFilterer) WatchLogValidatorAdded(opts *bind.WatchOpts, sink chan<- *CosmosBridgeLogValidatorAdded) (event.Subscription, error) {

	logs, sub, err := _CosmosBridge.contract.WatchLogs(opts, "LogValidatorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CosmosBridgeLogValidatorAdded)
				if err := _CosmosBridge.contract.UnpackLog(event, "LogValidatorAdded", log); err != nil {
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
func (_CosmosBridge *CosmosBridgeFilterer) ParseLogValidatorAdded(log types.Log) (*CosmosBridgeLogValidatorAdded, error) {
	event := new(CosmosBridgeLogValidatorAdded)
	if err := _CosmosBridge.contract.UnpackLog(event, "LogValidatorAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CosmosBridgeLogValidatorPowerUpdatedIterator is returned from FilterLogValidatorPowerUpdated and is used to iterate over the raw logs and unpacked data for LogValidatorPowerUpdated events raised by the CosmosBridge contract.
type CosmosBridgeLogValidatorPowerUpdatedIterator struct {
	Event *CosmosBridgeLogValidatorPowerUpdated // Event containing the contract specifics and raw log

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
func (it *CosmosBridgeLogValidatorPowerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CosmosBridgeLogValidatorPowerUpdated)
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
		it.Event = new(CosmosBridgeLogValidatorPowerUpdated)
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
func (it *CosmosBridgeLogValidatorPowerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CosmosBridgeLogValidatorPowerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CosmosBridgeLogValidatorPowerUpdated represents a LogValidatorPowerUpdated event raised by the CosmosBridge contract.
type CosmosBridgeLogValidatorPowerUpdated struct {
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
func (_CosmosBridge *CosmosBridgeFilterer) FilterLogValidatorPowerUpdated(opts *bind.FilterOpts) (*CosmosBridgeLogValidatorPowerUpdatedIterator, error) {

	logs, sub, err := _CosmosBridge.contract.FilterLogs(opts, "LogValidatorPowerUpdated")
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeLogValidatorPowerUpdatedIterator{contract: _CosmosBridge.contract, event: "LogValidatorPowerUpdated", logs: logs, sub: sub}, nil
}

// WatchLogValidatorPowerUpdated is a free log subscription operation binding the contract event 0x335940ce4119f8aae891d73dba74510a3d51f6210134d058237f26e6a31d5340.
//
// Solidity: event LogValidatorPowerUpdated(address _validator, uint256 _power, uint256 _currentValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_CosmosBridge *CosmosBridgeFilterer) WatchLogValidatorPowerUpdated(opts *bind.WatchOpts, sink chan<- *CosmosBridgeLogValidatorPowerUpdated) (event.Subscription, error) {

	logs, sub, err := _CosmosBridge.contract.WatchLogs(opts, "LogValidatorPowerUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CosmosBridgeLogValidatorPowerUpdated)
				if err := _CosmosBridge.contract.UnpackLog(event, "LogValidatorPowerUpdated", log); err != nil {
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
func (_CosmosBridge *CosmosBridgeFilterer) ParseLogValidatorPowerUpdated(log types.Log) (*CosmosBridgeLogValidatorPowerUpdated, error) {
	event := new(CosmosBridgeLogValidatorPowerUpdated)
	if err := _CosmosBridge.contract.UnpackLog(event, "LogValidatorPowerUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CosmosBridgeLogValidatorRemovedIterator is returned from FilterLogValidatorRemoved and is used to iterate over the raw logs and unpacked data for LogValidatorRemoved events raised by the CosmosBridge contract.
type CosmosBridgeLogValidatorRemovedIterator struct {
	Event *CosmosBridgeLogValidatorRemoved // Event containing the contract specifics and raw log

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
func (it *CosmosBridgeLogValidatorRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CosmosBridgeLogValidatorRemoved)
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
		it.Event = new(CosmosBridgeLogValidatorRemoved)
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
func (it *CosmosBridgeLogValidatorRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CosmosBridgeLogValidatorRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CosmosBridgeLogValidatorRemoved represents a LogValidatorRemoved event raised by the CosmosBridge contract.
type CosmosBridgeLogValidatorRemoved struct {
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
func (_CosmosBridge *CosmosBridgeFilterer) FilterLogValidatorRemoved(opts *bind.FilterOpts) (*CosmosBridgeLogValidatorRemovedIterator, error) {

	logs, sub, err := _CosmosBridge.contract.FilterLogs(opts, "LogValidatorRemoved")
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeLogValidatorRemovedIterator{contract: _CosmosBridge.contract, event: "LogValidatorRemoved", logs: logs, sub: sub}, nil
}

// WatchLogValidatorRemoved is a free log subscription operation binding the contract event 0x1241fb43a101ff98ab819a1882097d4ccada51ba60f326c1981cc48840f55b8c.
//
// Solidity: event LogValidatorRemoved(address _validator, uint256 _power, uint256 _currentValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_CosmosBridge *CosmosBridgeFilterer) WatchLogValidatorRemoved(opts *bind.WatchOpts, sink chan<- *CosmosBridgeLogValidatorRemoved) (event.Subscription, error) {

	logs, sub, err := _CosmosBridge.contract.WatchLogs(opts, "LogValidatorRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CosmosBridgeLogValidatorRemoved)
				if err := _CosmosBridge.contract.UnpackLog(event, "LogValidatorRemoved", log); err != nil {
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
func (_CosmosBridge *CosmosBridgeFilterer) ParseLogValidatorRemoved(log types.Log) (*CosmosBridgeLogValidatorRemoved, error) {
	event := new(CosmosBridgeLogValidatorRemoved)
	if err := _CosmosBridge.contract.UnpackLog(event, "LogValidatorRemoved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CosmosBridgeLogValsetResetIterator is returned from FilterLogValsetReset and is used to iterate over the raw logs and unpacked data for LogValsetReset events raised by the CosmosBridge contract.
type CosmosBridgeLogValsetResetIterator struct {
	Event *CosmosBridgeLogValsetReset // Event containing the contract specifics and raw log

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
func (it *CosmosBridgeLogValsetResetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CosmosBridgeLogValsetReset)
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
		it.Event = new(CosmosBridgeLogValsetReset)
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
func (it *CosmosBridgeLogValsetResetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CosmosBridgeLogValsetResetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CosmosBridgeLogValsetReset represents a LogValsetReset event raised by the CosmosBridge contract.
type CosmosBridgeLogValsetReset struct {
	NewValsetVersion *big.Int
	ValidatorCount   *big.Int
	TotalPower       *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogValsetReset is a free log retrieval operation binding the contract event 0xd870653e19f161500290fd0c4ca41bf5cf2bcb1ba66448f41c66c512dabd65f2.
//
// Solidity: event LogValsetReset(uint256 _newValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_CosmosBridge *CosmosBridgeFilterer) FilterLogValsetReset(opts *bind.FilterOpts) (*CosmosBridgeLogValsetResetIterator, error) {

	logs, sub, err := _CosmosBridge.contract.FilterLogs(opts, "LogValsetReset")
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeLogValsetResetIterator{contract: _CosmosBridge.contract, event: "LogValsetReset", logs: logs, sub: sub}, nil
}

// WatchLogValsetReset is a free log subscription operation binding the contract event 0xd870653e19f161500290fd0c4ca41bf5cf2bcb1ba66448f41c66c512dabd65f2.
//
// Solidity: event LogValsetReset(uint256 _newValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_CosmosBridge *CosmosBridgeFilterer) WatchLogValsetReset(opts *bind.WatchOpts, sink chan<- *CosmosBridgeLogValsetReset) (event.Subscription, error) {

	logs, sub, err := _CosmosBridge.contract.WatchLogs(opts, "LogValsetReset")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CosmosBridgeLogValsetReset)
				if err := _CosmosBridge.contract.UnpackLog(event, "LogValsetReset", log); err != nil {
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
func (_CosmosBridge *CosmosBridgeFilterer) ParseLogValsetReset(log types.Log) (*CosmosBridgeLogValsetReset, error) {
	event := new(CosmosBridgeLogValsetReset)
	if err := _CosmosBridge.contract.UnpackLog(event, "LogValsetReset", log); err != nil {
		return nil, err
	}
	return event, nil
}

// CosmosBridgeLogValsetUpdatedIterator is returned from FilterLogValsetUpdated and is used to iterate over the raw logs and unpacked data for LogValsetUpdated events raised by the CosmosBridge contract.
type CosmosBridgeLogValsetUpdatedIterator struct {
	Event *CosmosBridgeLogValsetUpdated // Event containing the contract specifics and raw log

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
func (it *CosmosBridgeLogValsetUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CosmosBridgeLogValsetUpdated)
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
		it.Event = new(CosmosBridgeLogValsetUpdated)
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
func (it *CosmosBridgeLogValsetUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CosmosBridgeLogValsetUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CosmosBridgeLogValsetUpdated represents a LogValsetUpdated event raised by the CosmosBridge contract.
type CosmosBridgeLogValsetUpdated struct {
	NewValsetVersion *big.Int
	ValidatorCount   *big.Int
	TotalPower       *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogValsetUpdated is a free log retrieval operation binding the contract event 0x3a7ef0da3179668af8114719645585b5a37092ef2d66f187dcf63d83a221eaa6.
//
// Solidity: event LogValsetUpdated(uint256 _newValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_CosmosBridge *CosmosBridgeFilterer) FilterLogValsetUpdated(opts *bind.FilterOpts) (*CosmosBridgeLogValsetUpdatedIterator, error) {

	logs, sub, err := _CosmosBridge.contract.FilterLogs(opts, "LogValsetUpdated")
	if err != nil {
		return nil, err
	}
	return &CosmosBridgeLogValsetUpdatedIterator{contract: _CosmosBridge.contract, event: "LogValsetUpdated", logs: logs, sub: sub}, nil
}

// WatchLogValsetUpdated is a free log subscription operation binding the contract event 0x3a7ef0da3179668af8114719645585b5a37092ef2d66f187dcf63d83a221eaa6.
//
// Solidity: event LogValsetUpdated(uint256 _newValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_CosmosBridge *CosmosBridgeFilterer) WatchLogValsetUpdated(opts *bind.WatchOpts, sink chan<- *CosmosBridgeLogValsetUpdated) (event.Subscription, error) {

	logs, sub, err := _CosmosBridge.contract.WatchLogs(opts, "LogValsetUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CosmosBridgeLogValsetUpdated)
				if err := _CosmosBridge.contract.UnpackLog(event, "LogValsetUpdated", log); err != nil {
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
func (_CosmosBridge *CosmosBridgeFilterer) ParseLogValsetUpdated(log types.Log) (*CosmosBridgeLogValsetUpdated, error) {
	event := new(CosmosBridgeLogValsetUpdated)
	if err := _CosmosBridge.contract.UnpackLog(event, "LogValsetUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}
