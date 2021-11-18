// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Valset

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

// ValsetABI is the input ABI used to generate the binding from.
const ValsetABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_power\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_currentValsetVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_validatorCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_totalPower\",\"type\":\"uint256\"}],\"name\":\"LogValidatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_power\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_currentValsetVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_validatorCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_totalPower\",\"type\":\"uint256\"}],\"name\":\"LogValidatorPowerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_power\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_currentValsetVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_validatorCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_totalPower\",\"type\":\"uint256\"}],\"name\":\"LogValidatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newValsetVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_validatorCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_totalPower\",\"type\":\"uint256\"}],\"name\":\"LogValsetReset\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_newValsetVersion\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_validatorCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_totalPower\",\"type\":\"uint256\"}],\"name\":\"LogValsetUpdated\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_validatorPower\",\"type\":\"uint256\"}],\"name\":\"addValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentValsetVersion\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"}],\"name\":\"getValidatorPower\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"_initValidators\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_initPowers\",\"type\":\"uint256[]\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"}],\"name\":\"isActiveValidator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"operator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"powers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_message\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_signature\",\"type\":\"bytes\"}],\"name\":\"recover\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_valsetVersion\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"}],\"name\":\"recoverGas\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"}],\"name\":\"removeValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalPower\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validatorAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_newValidatorPower\",\"type\":\"uint256\"}],\"name\":\"updateValidatorPower\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_validators\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_powers\",\"type\":\"uint256[]\"}],\"name\":\"updateValset\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"validatorCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// ValsetBin is the compiled bytecode used for deploying new contracts.
var ValsetBin = "0x608060405234801561001057600080fd5b50611991806100206000396000f3fe608060405234801561001057600080fd5b50600436106100f55760003560e01c8063788cf92f11610097578063b872c52311610066578063b872c523146105ab578063db3ad22c146105ed578063fc6c1f021461060b578063ff1d575214610659576100f5565b8063788cf92f146103ad5780638d56c37d146104f95780639bdafcb314610517578063b5672be31461055d576100f5565b806340550a1c116100d357806340550a1c1461026b57806340a141ff146102c7578063473691a41461030b578063570ca73514610363576100f5565b80630f43a677146100fa57806319045a25146101185780632e75293b1461021d575b600080fd5b6101026107c5565b6040518082815260200191505060405180910390f35b6101db6004803603604081101561012e57600080fd5b81019080803590602001909291908035906020019064010000000081111561015557600080fd5b82018360208201111561016757600080fd5b8035906020019184600183028401116401000000008311171561018957600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506107cb565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6102696004803603604081101561023357600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506107e7565b005b6102ad6004803603602081101561028157600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610a79565b604051808215151515815260200191505060405180910390f35b610309600480360360208110156102dd57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610b07565b005b61034d6004803603602081101561032157600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610db6565b6040518082815260200191505060405180910390f35b61036b610e37565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6104f7600480360360408110156103c357600080fd5b81019080803590602001906401000000008111156103e057600080fd5b8201836020820111156103f257600080fd5b8035906020019184602083028401116401000000008311171561041457600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f8201169050808301925050505050505091929192908035906020019064010000000081111561047457600080fd5b82018360208201111561048657600080fd5b803590602001918460208302840111640100000000831117156104a857600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f820116905080830192505050505050509192919290505050610e5d565b005b61050161102e565b6040518082815260200191505060405180910390f35b6105436004803603602081101561052d57600080fd5b8101908080359060200190929190505050611034565b604051808215151515815260200191505060405180910390f35b6105a96004803603604081101561057357600080fd5b8101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611054565b005b6105d7600480360360208110156105c157600080fd5b8101908080359060200190929190505050611212565b6040518082815260200191505060405180910390f35b6105f561122a565b6040518082815260200191505060405180910390f35b6106576004803603604081101561062157600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050611230565b005b6107c36004803603606081101561066f57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001906401000000008111156106ac57600080fd5b8201836020820111156106be57600080fd5b803590602001918460208302840111640100000000831117156106e057600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f8201169050808301925050505050505091929192908035906020019064010000000081111561074057600080fd5b82018360208201111561075257600080fd5b8035906020019184602083028401116401000000008311171561077457600080fd5b919080806020026020016040519081016040528093929190818152602001838360200280828437600081840152601f19601f820116905080830192505050505050509192919290505050611301565b005b60025481565b60006107df6107d9846113f7565b8361144f565b905092915050565b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146108aa576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f4d75737420626520746865206f70657261746f722e000000000000000000000081525060200191505060405180910390fd5b600060015483604051602001808381526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b8152601401925050506040516020818303038152906040528051906020012090506003600082815260200190815260200160002060009054906101000a900460ff16610984576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602e815260200180611900602e913960400191505060405180910390fd5b6000600560008381526020019081526020016000205490506109b18160005461153190919063ffffffff16565b6000819055506109cc8360005461157b90919063ffffffff16565b6000819055508260056000848152602001908152602001600020819055507f335940ce4119f8aae891d73dba74510a3d51f6210134d058237f26e6a31d53408484600154600254600054604051808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018581526020018481526020018381526020018281526020019550505050505060405180910390a150505050565b60008060015483604051602001808381526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b8152601401925050506040516020818303038152906040528051906020012090506003600082815260200190815260200160002060009054906101000a900460ff16915050919050565b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610bca576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f4d75737420626520746865206f70657261746f722e000000000000000000000081525060200191505060405180910390fd5b600060015482604051602001808381526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b8152601401925050506040516020818303038152906040528051906020012090506003600082815260200190815260200160002060009054906101000a900460ff16610ca4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806118aa6021913960400191505060405180910390fd5b610cba600160025461153190919063ffffffff16565b600281905550610ce8600560008381526020019081526020016000205460005461153190919063ffffffff16565b6000819055506003600082815260200190815260200160002060006101000a81549060ff021916905560056000828152602001908152602001600020600090557f1241fb43a101ff98ab819a1882097d4ccada51ba60f326c1981cc48840f55b8c826000600154600254600054604051808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018581526020018481526020018381526020018281526020019550505050505060405180910390a15050565b60008060015483604051602001808381526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b8152601401925050506040516020818303038152906040528051906020012090506005600082815260200190815260200160002054915050919050565b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610f20576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f4d75737420626520746865206f70657261746f722e000000000000000000000081525060200191505060405180910390fd5b8051825114610f7a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602f81526020018061192e602f913960400191505060405180910390fd5b610f82611603565b60008090505b8251811015610fdc57610fc1838281518110610fa057fe5b6020026020010151838381518110610fb457fe5b602002602001015161167c565b610fd560018261157b90919063ffffffff16565b9050610f88565b507f3a7ef0da3179668af8114719645585b5a37092ef2d66f187dcf63d83a221eaa660015460025460005460405180848152602001838152602001828152602001935050505060405180910390a15050565b60015481565b60036020528060005260406000206000915054906101000a900460ff1681565b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614611117576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f4d75737420626520746865206f70657261746f722e000000000000000000000081525060200191505060405180910390fd5b6001548210611171576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260358152602001806118cb6035913960400191505060405180910390fd5b60008282604051602001808381526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b8152601401925050506040516020818303038152906040528051906020012090506003600082815260200190815260200160002060006101000a81549060ff02191690556005600082815260200190815260200160002060009055505050565b60056020528060005260406000206000915090505481565b60005481565b600460009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146112f3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260158152602001807f4d75737420626520746865206f70657261746f722e000000000000000000000081525060200191505060405180910390fd5b6112fd828261167c565b5050565b606a60009054906101000a900460ff1615611384576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600b8152602001807f496e697469616c697a656400000000000000000000000000000000000000000081525060200191505060405180910390fd5b82600460006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060006001819055506001606a60006101000a81548160ff0219169083151502179055506113f28282610e5d565b505050565b60008160405160200180807f19457468657265756d205369676e6564204d6573736167653a0a333200000000815250601c01828152602001915050604051602081830303815290604052805190602001209050919050565b600080600080604185511461146a576000935050505061152b565b6020850151925060408501519150606085015160001a9050601b8160ff16101561149557601b810190505b601b8160ff16141580156114ad5750601c8160ff1614155b156114be576000935050505061152b565b60018682858560405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa15801561151b573d6000803e3d6000fd5b5050506020604051035193505050505b92915050565b600061157383836040518060400160405280601e81526020017f536166654d6174683a207375627472616374696f6e206f766572666c6f7700008152506117e9565b905092915050565b6000808284019050838110156115f9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b6116186001805461157b90919063ffffffff16565b6001819055506000600281905550600080819055507fd870653e19f161500290fd0c4ca41bf5cf2bcb1ba66448f41c66c512dabd65f260015460025460005460405180848152602001838152602001828152602001935050505060405180910390a1565b600060015483604051602001808381526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1660601b8152601401925050506040516020818303038152906040528051906020012090506116f6600160025461157b90919063ffffffff16565b6002819055506117118260005461157b90919063ffffffff16565b60008190555060016003600083815260200190815260200160002060006101000a81548160ff0219169083151502179055508160056000838152602001908152602001600020819055507f1a396bcf647502e902dce665d58a0c1b25f982f193ab9a1d0f1500d8d927bf2a8383600154600254600054604051808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018581526020018481526020018381526020018281526020019550505050505060405180910390a1505050565b6000838311158290611896576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b8381101561185b578082015181840152602081019050611840565b50505050905090810190601f1680156118885780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b506000838503905080915050939250505056fe43616e206f6e6c792072656d6f7665206163746976652076616c646961746f7273476173207265636f76657279206f6e6c7920616c6c6f77656420666f722070726576696f75732076616c696461746f72207365747343616e206f6e6c79207570646174652074686520706f776572206f66206163746976652076616c646961746f727345766572792076616c696461746f72206d7573742068617665206120636f72726573706f6e64696e6720706f776572a265627a7a72315820d20fa526b20b6d12bfc264cf543fb4bb97d46e86ff133dc7bbdb6b4905a5b83c64736f6c63430005110032"

// DeployValset deploys a new Ethereum contract, binding an instance of Valset to it.
func DeployValset(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Valset, error) {
	parsed, err := abi.JSON(strings.NewReader(ValsetABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ValsetBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Valset{ValsetCaller: ValsetCaller{contract: contract}, ValsetTransactor: ValsetTransactor{contract: contract}, ValsetFilterer: ValsetFilterer{contract: contract}}, nil
}

// Valset is an auto generated Go binding around an Ethereum contract.
type Valset struct {
	ValsetCaller     // Read-only binding to the contract
	ValsetTransactor // Write-only binding to the contract
	ValsetFilterer   // Log filterer for contract events
}

// ValsetCaller is an auto generated read-only Go binding around an Ethereum contract.
type ValsetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValsetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ValsetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValsetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ValsetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValsetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ValsetSession struct {
	Contract     *Valset           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ValsetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ValsetCallerSession struct {
	Contract *ValsetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ValsetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ValsetTransactorSession struct {
	Contract     *ValsetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ValsetRaw is an auto generated low-level Go binding around an Ethereum contract.
type ValsetRaw struct {
	Contract *Valset // Generic contract binding to access the raw methods on
}

// ValsetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ValsetCallerRaw struct {
	Contract *ValsetCaller // Generic read-only contract binding to access the raw methods on
}

// ValsetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ValsetTransactorRaw struct {
	Contract *ValsetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewValset creates a new instance of Valset, bound to a specific deployed contract.
func NewValset(address common.Address, backend bind.ContractBackend) (*Valset, error) {
	contract, err := bindValset(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Valset{ValsetCaller: ValsetCaller{contract: contract}, ValsetTransactor: ValsetTransactor{contract: contract}, ValsetFilterer: ValsetFilterer{contract: contract}}, nil
}

// NewValsetCaller creates a new read-only instance of Valset, bound to a specific deployed contract.
func NewValsetCaller(address common.Address, caller bind.ContractCaller) (*ValsetCaller, error) {
	contract, err := bindValset(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ValsetCaller{contract: contract}, nil
}

// NewValsetTransactor creates a new write-only instance of Valset, bound to a specific deployed contract.
func NewValsetTransactor(address common.Address, transactor bind.ContractTransactor) (*ValsetTransactor, error) {
	contract, err := bindValset(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ValsetTransactor{contract: contract}, nil
}

// NewValsetFilterer creates a new log filterer instance of Valset, bound to a specific deployed contract.
func NewValsetFilterer(address common.Address, filterer bind.ContractFilterer) (*ValsetFilterer, error) {
	contract, err := bindValset(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ValsetFilterer{contract: contract}, nil
}

// bindValset binds a generic wrapper to an already deployed contract.
func bindValset(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ValsetABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Valset *ValsetRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Valset.Contract.ValsetCaller.contract.Call(opts, &[]interface{}{result}, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Valset *ValsetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Valset.Contract.ValsetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Valset *ValsetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Valset.Contract.ValsetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Valset *ValsetCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Valset.Contract.contract.Call(opts, &[]interface{}{result}, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Valset *ValsetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Valset.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Valset *ValsetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Valset.Contract.contract.Transact(opts, method, params...)
}

// CurrentValsetVersion is a free data retrieval call binding the contract method 0x8d56c37d.
//
// Solidity: function currentValsetVersion() view returns(uint256)
func (_Valset *ValsetCaller) CurrentValsetVersion(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Valset.contract.Call(opts, &[]interface{}{out}, "currentValsetVersion")
	return *ret0, err
}

// CurrentValsetVersion is a free data retrieval call binding the contract method 0x8d56c37d.
//
// Solidity: function currentValsetVersion() view returns(uint256)
func (_Valset *ValsetSession) CurrentValsetVersion() (*big.Int, error) {
	return _Valset.Contract.CurrentValsetVersion(&_Valset.CallOpts)
}

// CurrentValsetVersion is a free data retrieval call binding the contract method 0x8d56c37d.
//
// Solidity: function currentValsetVersion() view returns(uint256)
func (_Valset *ValsetCallerSession) CurrentValsetVersion() (*big.Int, error) {
	return _Valset.Contract.CurrentValsetVersion(&_Valset.CallOpts)
}

// GetValidatorPower is a free data retrieval call binding the contract method 0x473691a4.
//
// Solidity: function getValidatorPower(address _validatorAddress) view returns(uint256)
func (_Valset *ValsetCaller) GetValidatorPower(opts *bind.CallOpts, _validatorAddress common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Valset.contract.Call(opts, &[]interface{}{out}, "getValidatorPower", _validatorAddress)
	return *ret0, err
}

// GetValidatorPower is a free data retrieval call binding the contract method 0x473691a4.
//
// Solidity: function getValidatorPower(address _validatorAddress) view returns(uint256)
func (_Valset *ValsetSession) GetValidatorPower(_validatorAddress common.Address) (*big.Int, error) {
	return _Valset.Contract.GetValidatorPower(&_Valset.CallOpts, _validatorAddress)
}

// GetValidatorPower is a free data retrieval call binding the contract method 0x473691a4.
//
// Solidity: function getValidatorPower(address _validatorAddress) view returns(uint256)
func (_Valset *ValsetCallerSession) GetValidatorPower(_validatorAddress common.Address) (*big.Int, error) {
	return _Valset.Contract.GetValidatorPower(&_Valset.CallOpts, _validatorAddress)
}

// IsActiveValidator is a free data retrieval call binding the contract method 0x40550a1c.
//
// Solidity: function isActiveValidator(address _validatorAddress) view returns(bool)
func (_Valset *ValsetCaller) IsActiveValidator(opts *bind.CallOpts, _validatorAddress common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Valset.contract.Call(opts, &[]interface{}{out}, "isActiveValidator", _validatorAddress)
	return *ret0, err
}

// IsActiveValidator is a free data retrieval call binding the contract method 0x40550a1c.
//
// Solidity: function isActiveValidator(address _validatorAddress) view returns(bool)
func (_Valset *ValsetSession) IsActiveValidator(_validatorAddress common.Address) (bool, error) {
	return _Valset.Contract.IsActiveValidator(&_Valset.CallOpts, _validatorAddress)
}

// IsActiveValidator is a free data retrieval call binding the contract method 0x40550a1c.
//
// Solidity: function isActiveValidator(address _validatorAddress) view returns(bool)
func (_Valset *ValsetCallerSession) IsActiveValidator(_validatorAddress common.Address) (bool, error) {
	return _Valset.Contract.IsActiveValidator(&_Valset.CallOpts, _validatorAddress)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_Valset *ValsetCaller) Operator(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Valset.contract.Call(opts, &[]interface{}{out}, "operator")
	return *ret0, err
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_Valset *ValsetSession) Operator() (common.Address, error) {
	return _Valset.Contract.Operator(&_Valset.CallOpts)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_Valset *ValsetCallerSession) Operator() (common.Address, error) {
	return _Valset.Contract.Operator(&_Valset.CallOpts)
}

// Powers is a free data retrieval call binding the contract method 0xb872c523.
//
// Solidity: function powers(bytes32 ) view returns(uint256)
func (_Valset *ValsetCaller) Powers(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Valset.contract.Call(opts, &[]interface{}{out}, "powers", arg0)
	return *ret0, err
}

// Powers is a free data retrieval call binding the contract method 0xb872c523.
//
// Solidity: function powers(bytes32 ) view returns(uint256)
func (_Valset *ValsetSession) Powers(arg0 [32]byte) (*big.Int, error) {
	return _Valset.Contract.Powers(&_Valset.CallOpts, arg0)
}

// Powers is a free data retrieval call binding the contract method 0xb872c523.
//
// Solidity: function powers(bytes32 ) view returns(uint256)
func (_Valset *ValsetCallerSession) Powers(arg0 [32]byte) (*big.Int, error) {
	return _Valset.Contract.Powers(&_Valset.CallOpts, arg0)
}

// Recover is a free data retrieval call binding the contract method 0x19045a25.
//
// Solidity: function recover(bytes32 _message, bytes _signature) pure returns(address)
func (_Valset *ValsetCaller) Recover(opts *bind.CallOpts, _message [32]byte, _signature []byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Valset.contract.Call(opts, &[]interface{}{out}, "recover", _message, _signature)
	return *ret0, err
}

// Recover is a free data retrieval call binding the contract method 0x19045a25.
//
// Solidity: function recover(bytes32 _message, bytes _signature) pure returns(address)
func (_Valset *ValsetSession) Recover(_message [32]byte, _signature []byte) (common.Address, error) {
	return _Valset.Contract.Recover(&_Valset.CallOpts, _message, _signature)
}

// Recover is a free data retrieval call binding the contract method 0x19045a25.
//
// Solidity: function recover(bytes32 _message, bytes _signature) pure returns(address)
func (_Valset *ValsetCallerSession) Recover(_message [32]byte, _signature []byte) (common.Address, error) {
	return _Valset.Contract.Recover(&_Valset.CallOpts, _message, _signature)
}

// TotalPower is a free data retrieval call binding the contract method 0xdb3ad22c.
//
// Solidity: function totalPower() view returns(uint256)
func (_Valset *ValsetCaller) TotalPower(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Valset.contract.Call(opts, &[]interface{}{out}, "totalPower")
	return *ret0, err
}

// TotalPower is a free data retrieval call binding the contract method 0xdb3ad22c.
//
// Solidity: function totalPower() view returns(uint256)
func (_Valset *ValsetSession) TotalPower() (*big.Int, error) {
	return _Valset.Contract.TotalPower(&_Valset.CallOpts)
}

// TotalPower is a free data retrieval call binding the contract method 0xdb3ad22c.
//
// Solidity: function totalPower() view returns(uint256)
func (_Valset *ValsetCallerSession) TotalPower() (*big.Int, error) {
	return _Valset.Contract.TotalPower(&_Valset.CallOpts)
}

// ValidatorCount is a free data retrieval call binding the contract method 0x0f43a677.
//
// Solidity: function validatorCount() view returns(uint256)
func (_Valset *ValsetCaller) ValidatorCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Valset.contract.Call(opts, &[]interface{}{out}, "validatorCount")
	return *ret0, err
}

// ValidatorCount is a free data retrieval call binding the contract method 0x0f43a677.
//
// Solidity: function validatorCount() view returns(uint256)
func (_Valset *ValsetSession) ValidatorCount() (*big.Int, error) {
	return _Valset.Contract.ValidatorCount(&_Valset.CallOpts)
}

// ValidatorCount is a free data retrieval call binding the contract method 0x0f43a677.
//
// Solidity: function validatorCount() view returns(uint256)
func (_Valset *ValsetCallerSession) ValidatorCount() (*big.Int, error) {
	return _Valset.Contract.ValidatorCount(&_Valset.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0x9bdafcb3.
//
// Solidity: function validators(bytes32 ) view returns(bool)
func (_Valset *ValsetCaller) Validators(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Valset.contract.Call(opts, &[]interface{}{out}, "validators", arg0)
	return *ret0, err
}

// Validators is a free data retrieval call binding the contract method 0x9bdafcb3.
//
// Solidity: function validators(bytes32 ) view returns(bool)
func (_Valset *ValsetSession) Validators(arg0 [32]byte) (bool, error) {
	return _Valset.Contract.Validators(&_Valset.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0x9bdafcb3.
//
// Solidity: function validators(bytes32 ) view returns(bool)
func (_Valset *ValsetCallerSession) Validators(arg0 [32]byte) (bool, error) {
	return _Valset.Contract.Validators(&_Valset.CallOpts, arg0)
}

// AddValidator is a paid mutator transaction binding the contract method 0xfc6c1f02.
//
// Solidity: function addValidator(address _validatorAddress, uint256 _validatorPower) returns()
func (_Valset *ValsetTransactor) AddValidator(opts *bind.TransactOpts, _validatorAddress common.Address, _validatorPower *big.Int) (*types.Transaction, error) {
	return _Valset.contract.Transact(opts, "addValidator", _validatorAddress, _validatorPower)
}

// AddValidator is a paid mutator transaction binding the contract method 0xfc6c1f02.
//
// Solidity: function addValidator(address _validatorAddress, uint256 _validatorPower) returns()
func (_Valset *ValsetSession) AddValidator(_validatorAddress common.Address, _validatorPower *big.Int) (*types.Transaction, error) {
	return _Valset.Contract.AddValidator(&_Valset.TransactOpts, _validatorAddress, _validatorPower)
}

// AddValidator is a paid mutator transaction binding the contract method 0xfc6c1f02.
//
// Solidity: function addValidator(address _validatorAddress, uint256 _validatorPower) returns()
func (_Valset *ValsetTransactorSession) AddValidator(_validatorAddress common.Address, _validatorPower *big.Int) (*types.Transaction, error) {
	return _Valset.Contract.AddValidator(&_Valset.TransactOpts, _validatorAddress, _validatorPower)
}

// Initialize is a paid mutator transaction binding the contract method 0xff1d5752.
//
// Solidity: function initialize(address _operator, address[] _initValidators, uint256[] _initPowers) returns()
func (_Valset *ValsetTransactor) Initialize(opts *bind.TransactOpts, _operator common.Address, _initValidators []common.Address, _initPowers []*big.Int) (*types.Transaction, error) {
	return _Valset.contract.Transact(opts, "initialize", _operator, _initValidators, _initPowers)
}

// Initialize is a paid mutator transaction binding the contract method 0xff1d5752.
//
// Solidity: function initialize(address _operator, address[] _initValidators, uint256[] _initPowers) returns()
func (_Valset *ValsetSession) Initialize(_operator common.Address, _initValidators []common.Address, _initPowers []*big.Int) (*types.Transaction, error) {
	return _Valset.Contract.Initialize(&_Valset.TransactOpts, _operator, _initValidators, _initPowers)
}

// Initialize is a paid mutator transaction binding the contract method 0xff1d5752.
//
// Solidity: function initialize(address _operator, address[] _initValidators, uint256[] _initPowers) returns()
func (_Valset *ValsetTransactorSession) Initialize(_operator common.Address, _initValidators []common.Address, _initPowers []*big.Int) (*types.Transaction, error) {
	return _Valset.Contract.Initialize(&_Valset.TransactOpts, _operator, _initValidators, _initPowers)
}

// RecoverGas is a paid mutator transaction binding the contract method 0xb5672be3.
//
// Solidity: function recoverGas(uint256 _valsetVersion, address _validatorAddress) returns()
func (_Valset *ValsetTransactor) RecoverGas(opts *bind.TransactOpts, _valsetVersion *big.Int, _validatorAddress common.Address) (*types.Transaction, error) {
	return _Valset.contract.Transact(opts, "recoverGas", _valsetVersion, _validatorAddress)
}

// RecoverGas is a paid mutator transaction binding the contract method 0xb5672be3.
//
// Solidity: function recoverGas(uint256 _valsetVersion, address _validatorAddress) returns()
func (_Valset *ValsetSession) RecoverGas(_valsetVersion *big.Int, _validatorAddress common.Address) (*types.Transaction, error) {
	return _Valset.Contract.RecoverGas(&_Valset.TransactOpts, _valsetVersion, _validatorAddress)
}

// RecoverGas is a paid mutator transaction binding the contract method 0xb5672be3.
//
// Solidity: function recoverGas(uint256 _valsetVersion, address _validatorAddress) returns()
func (_Valset *ValsetTransactorSession) RecoverGas(_valsetVersion *big.Int, _validatorAddress common.Address) (*types.Transaction, error) {
	return _Valset.Contract.RecoverGas(&_Valset.TransactOpts, _valsetVersion, _validatorAddress)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(address _validatorAddress) returns()
func (_Valset *ValsetTransactor) RemoveValidator(opts *bind.TransactOpts, _validatorAddress common.Address) (*types.Transaction, error) {
	return _Valset.contract.Transact(opts, "removeValidator", _validatorAddress)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(address _validatorAddress) returns()
func (_Valset *ValsetSession) RemoveValidator(_validatorAddress common.Address) (*types.Transaction, error) {
	return _Valset.Contract.RemoveValidator(&_Valset.TransactOpts, _validatorAddress)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(address _validatorAddress) returns()
func (_Valset *ValsetTransactorSession) RemoveValidator(_validatorAddress common.Address) (*types.Transaction, error) {
	return _Valset.Contract.RemoveValidator(&_Valset.TransactOpts, _validatorAddress)
}

// UpdateValidatorPower is a paid mutator transaction binding the contract method 0x2e75293b.
//
// Solidity: function updateValidatorPower(address _validatorAddress, uint256 _newValidatorPower) returns()
func (_Valset *ValsetTransactor) UpdateValidatorPower(opts *bind.TransactOpts, _validatorAddress common.Address, _newValidatorPower *big.Int) (*types.Transaction, error) {
	return _Valset.contract.Transact(opts, "updateValidatorPower", _validatorAddress, _newValidatorPower)
}

// UpdateValidatorPower is a paid mutator transaction binding the contract method 0x2e75293b.
//
// Solidity: function updateValidatorPower(address _validatorAddress, uint256 _newValidatorPower) returns()
func (_Valset *ValsetSession) UpdateValidatorPower(_validatorAddress common.Address, _newValidatorPower *big.Int) (*types.Transaction, error) {
	return _Valset.Contract.UpdateValidatorPower(&_Valset.TransactOpts, _validatorAddress, _newValidatorPower)
}

// UpdateValidatorPower is a paid mutator transaction binding the contract method 0x2e75293b.
//
// Solidity: function updateValidatorPower(address _validatorAddress, uint256 _newValidatorPower) returns()
func (_Valset *ValsetTransactorSession) UpdateValidatorPower(_validatorAddress common.Address, _newValidatorPower *big.Int) (*types.Transaction, error) {
	return _Valset.Contract.UpdateValidatorPower(&_Valset.TransactOpts, _validatorAddress, _newValidatorPower)
}

// UpdateValset is a paid mutator transaction binding the contract method 0x788cf92f.
//
// Solidity: function updateValset(address[] _validators, uint256[] _powers) returns()
func (_Valset *ValsetTransactor) UpdateValset(opts *bind.TransactOpts, _validators []common.Address, _powers []*big.Int) (*types.Transaction, error) {
	return _Valset.contract.Transact(opts, "updateValset", _validators, _powers)
}

// UpdateValset is a paid mutator transaction binding the contract method 0x788cf92f.
//
// Solidity: function updateValset(address[] _validators, uint256[] _powers) returns()
func (_Valset *ValsetSession) UpdateValset(_validators []common.Address, _powers []*big.Int) (*types.Transaction, error) {
	return _Valset.Contract.UpdateValset(&_Valset.TransactOpts, _validators, _powers)
}

// UpdateValset is a paid mutator transaction binding the contract method 0x788cf92f.
//
// Solidity: function updateValset(address[] _validators, uint256[] _powers) returns()
func (_Valset *ValsetTransactorSession) UpdateValset(_validators []common.Address, _powers []*big.Int) (*types.Transaction, error) {
	return _Valset.Contract.UpdateValset(&_Valset.TransactOpts, _validators, _powers)
}

// ValsetLogValidatorAddedIterator is returned from FilterLogValidatorAdded and is used to iterate over the raw logs and unpacked data for LogValidatorAdded events raised by the Valset contract.
type ValsetLogValidatorAddedIterator struct {
	Event *ValsetLogValidatorAdded // Event containing the contract specifics and raw log

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
func (it *ValsetLogValidatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValsetLogValidatorAdded)
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
		it.Event = new(ValsetLogValidatorAdded)
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
func (it *ValsetLogValidatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValsetLogValidatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValsetLogValidatorAdded represents a LogValidatorAdded event raised by the Valset contract.
type ValsetLogValidatorAdded struct {
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
func (_Valset *ValsetFilterer) FilterLogValidatorAdded(opts *bind.FilterOpts) (*ValsetLogValidatorAddedIterator, error) {

	logs, sub, err := _Valset.contract.FilterLogs(opts, "LogValidatorAdded")
	if err != nil {
		return nil, err
	}
	return &ValsetLogValidatorAddedIterator{contract: _Valset.contract, event: "LogValidatorAdded", logs: logs, sub: sub}, nil
}

// WatchLogValidatorAdded is a free log subscription operation binding the contract event 0x1a396bcf647502e902dce665d58a0c1b25f982f193ab9a1d0f1500d8d927bf2a.
//
// Solidity: event LogValidatorAdded(address _validator, uint256 _power, uint256 _currentValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Valset *ValsetFilterer) WatchLogValidatorAdded(opts *bind.WatchOpts, sink chan<- *ValsetLogValidatorAdded) (event.Subscription, error) {

	logs, sub, err := _Valset.contract.WatchLogs(opts, "LogValidatorAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValsetLogValidatorAdded)
				if err := _Valset.contract.UnpackLog(event, "LogValidatorAdded", log); err != nil {
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
func (_Valset *ValsetFilterer) ParseLogValidatorAdded(log types.Log) (*ValsetLogValidatorAdded, error) {
	event := new(ValsetLogValidatorAdded)
	if err := _Valset.contract.UnpackLog(event, "LogValidatorAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ValsetLogValidatorPowerUpdatedIterator is returned from FilterLogValidatorPowerUpdated and is used to iterate over the raw logs and unpacked data for LogValidatorPowerUpdated events raised by the Valset contract.
type ValsetLogValidatorPowerUpdatedIterator struct {
	Event *ValsetLogValidatorPowerUpdated // Event containing the contract specifics and raw log

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
func (it *ValsetLogValidatorPowerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValsetLogValidatorPowerUpdated)
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
		it.Event = new(ValsetLogValidatorPowerUpdated)
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
func (it *ValsetLogValidatorPowerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValsetLogValidatorPowerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValsetLogValidatorPowerUpdated represents a LogValidatorPowerUpdated event raised by the Valset contract.
type ValsetLogValidatorPowerUpdated struct {
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
func (_Valset *ValsetFilterer) FilterLogValidatorPowerUpdated(opts *bind.FilterOpts) (*ValsetLogValidatorPowerUpdatedIterator, error) {

	logs, sub, err := _Valset.contract.FilterLogs(opts, "LogValidatorPowerUpdated")
	if err != nil {
		return nil, err
	}
	return &ValsetLogValidatorPowerUpdatedIterator{contract: _Valset.contract, event: "LogValidatorPowerUpdated", logs: logs, sub: sub}, nil
}

// WatchLogValidatorPowerUpdated is a free log subscription operation binding the contract event 0x335940ce4119f8aae891d73dba74510a3d51f6210134d058237f26e6a31d5340.
//
// Solidity: event LogValidatorPowerUpdated(address _validator, uint256 _power, uint256 _currentValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Valset *ValsetFilterer) WatchLogValidatorPowerUpdated(opts *bind.WatchOpts, sink chan<- *ValsetLogValidatorPowerUpdated) (event.Subscription, error) {

	logs, sub, err := _Valset.contract.WatchLogs(opts, "LogValidatorPowerUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValsetLogValidatorPowerUpdated)
				if err := _Valset.contract.UnpackLog(event, "LogValidatorPowerUpdated", log); err != nil {
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
func (_Valset *ValsetFilterer) ParseLogValidatorPowerUpdated(log types.Log) (*ValsetLogValidatorPowerUpdated, error) {
	event := new(ValsetLogValidatorPowerUpdated)
	if err := _Valset.contract.UnpackLog(event, "LogValidatorPowerUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ValsetLogValidatorRemovedIterator is returned from FilterLogValidatorRemoved and is used to iterate over the raw logs and unpacked data for LogValidatorRemoved events raised by the Valset contract.
type ValsetLogValidatorRemovedIterator struct {
	Event *ValsetLogValidatorRemoved // Event containing the contract specifics and raw log

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
func (it *ValsetLogValidatorRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValsetLogValidatorRemoved)
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
		it.Event = new(ValsetLogValidatorRemoved)
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
func (it *ValsetLogValidatorRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValsetLogValidatorRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValsetLogValidatorRemoved represents a LogValidatorRemoved event raised by the Valset contract.
type ValsetLogValidatorRemoved struct {
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
func (_Valset *ValsetFilterer) FilterLogValidatorRemoved(opts *bind.FilterOpts) (*ValsetLogValidatorRemovedIterator, error) {

	logs, sub, err := _Valset.contract.FilterLogs(opts, "LogValidatorRemoved")
	if err != nil {
		return nil, err
	}
	return &ValsetLogValidatorRemovedIterator{contract: _Valset.contract, event: "LogValidatorRemoved", logs: logs, sub: sub}, nil
}

// WatchLogValidatorRemoved is a free log subscription operation binding the contract event 0x1241fb43a101ff98ab819a1882097d4ccada51ba60f326c1981cc48840f55b8c.
//
// Solidity: event LogValidatorRemoved(address _validator, uint256 _power, uint256 _currentValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Valset *ValsetFilterer) WatchLogValidatorRemoved(opts *bind.WatchOpts, sink chan<- *ValsetLogValidatorRemoved) (event.Subscription, error) {

	logs, sub, err := _Valset.contract.WatchLogs(opts, "LogValidatorRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValsetLogValidatorRemoved)
				if err := _Valset.contract.UnpackLog(event, "LogValidatorRemoved", log); err != nil {
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
func (_Valset *ValsetFilterer) ParseLogValidatorRemoved(log types.Log) (*ValsetLogValidatorRemoved, error) {
	event := new(ValsetLogValidatorRemoved)
	if err := _Valset.contract.UnpackLog(event, "LogValidatorRemoved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ValsetLogValsetResetIterator is returned from FilterLogValsetReset and is used to iterate over the raw logs and unpacked data for LogValsetReset events raised by the Valset contract.
type ValsetLogValsetResetIterator struct {
	Event *ValsetLogValsetReset // Event containing the contract specifics and raw log

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
func (it *ValsetLogValsetResetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValsetLogValsetReset)
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
		it.Event = new(ValsetLogValsetReset)
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
func (it *ValsetLogValsetResetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValsetLogValsetResetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValsetLogValsetReset represents a LogValsetReset event raised by the Valset contract.
type ValsetLogValsetReset struct {
	NewValsetVersion *big.Int
	ValidatorCount   *big.Int
	TotalPower       *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogValsetReset is a free log retrieval operation binding the contract event 0xd870653e19f161500290fd0c4ca41bf5cf2bcb1ba66448f41c66c512dabd65f2.
//
// Solidity: event LogValsetReset(uint256 _newValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Valset *ValsetFilterer) FilterLogValsetReset(opts *bind.FilterOpts) (*ValsetLogValsetResetIterator, error) {

	logs, sub, err := _Valset.contract.FilterLogs(opts, "LogValsetReset")
	if err != nil {
		return nil, err
	}
	return &ValsetLogValsetResetIterator{contract: _Valset.contract, event: "LogValsetReset", logs: logs, sub: sub}, nil
}

// WatchLogValsetReset is a free log subscription operation binding the contract event 0xd870653e19f161500290fd0c4ca41bf5cf2bcb1ba66448f41c66c512dabd65f2.
//
// Solidity: event LogValsetReset(uint256 _newValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Valset *ValsetFilterer) WatchLogValsetReset(opts *bind.WatchOpts, sink chan<- *ValsetLogValsetReset) (event.Subscription, error) {

	logs, sub, err := _Valset.contract.WatchLogs(opts, "LogValsetReset")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValsetLogValsetReset)
				if err := _Valset.contract.UnpackLog(event, "LogValsetReset", log); err != nil {
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
func (_Valset *ValsetFilterer) ParseLogValsetReset(log types.Log) (*ValsetLogValsetReset, error) {
	event := new(ValsetLogValsetReset)
	if err := _Valset.contract.UnpackLog(event, "LogValsetReset", log); err != nil {
		return nil, err
	}
	return event, nil
}

// ValsetLogValsetUpdatedIterator is returned from FilterLogValsetUpdated and is used to iterate over the raw logs and unpacked data for LogValsetUpdated events raised by the Valset contract.
type ValsetLogValsetUpdatedIterator struct {
	Event *ValsetLogValsetUpdated // Event containing the contract specifics and raw log

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
func (it *ValsetLogValsetUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValsetLogValsetUpdated)
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
		it.Event = new(ValsetLogValsetUpdated)
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
func (it *ValsetLogValsetUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValsetLogValsetUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValsetLogValsetUpdated represents a LogValsetUpdated event raised by the Valset contract.
type ValsetLogValsetUpdated struct {
	NewValsetVersion *big.Int
	ValidatorCount   *big.Int
	TotalPower       *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterLogValsetUpdated is a free log retrieval operation binding the contract event 0x3a7ef0da3179668af8114719645585b5a37092ef2d66f187dcf63d83a221eaa6.
//
// Solidity: event LogValsetUpdated(uint256 _newValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Valset *ValsetFilterer) FilterLogValsetUpdated(opts *bind.FilterOpts) (*ValsetLogValsetUpdatedIterator, error) {

	logs, sub, err := _Valset.contract.FilterLogs(opts, "LogValsetUpdated")
	if err != nil {
		return nil, err
	}
	return &ValsetLogValsetUpdatedIterator{contract: _Valset.contract, event: "LogValsetUpdated", logs: logs, sub: sub}, nil
}

// WatchLogValsetUpdated is a free log subscription operation binding the contract event 0x3a7ef0da3179668af8114719645585b5a37092ef2d66f187dcf63d83a221eaa6.
//
// Solidity: event LogValsetUpdated(uint256 _newValsetVersion, uint256 _validatorCount, uint256 _totalPower)
func (_Valset *ValsetFilterer) WatchLogValsetUpdated(opts *bind.WatchOpts, sink chan<- *ValsetLogValsetUpdated) (event.Subscription, error) {

	logs, sub, err := _Valset.contract.WatchLogs(opts, "LogValsetUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValsetLogValsetUpdated)
				if err := _Valset.contract.UnpackLog(event, "LogValsetUpdated", log); err != nil {
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
func (_Valset *ValsetFilterer) ParseLogValsetUpdated(log types.Log) (*ValsetLogValsetUpdated, error) {
	event := new(ValsetLogValsetUpdated)
	if err := _Valset.contract.UnpackLog(event, "LogValsetUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}
