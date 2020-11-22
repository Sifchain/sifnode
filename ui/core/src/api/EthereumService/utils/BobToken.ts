// @ts-ignore
module.exports =  {
  "contractName": "BobToken",
  "abi": [
    {
      "inputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "constructor"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "spender",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "Approval",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "Transfer",
      "type": "event"
    },
    {
      "constant": true,
      "inputs": [
        {
          "internalType": "address",
          "name": "owner",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        }
      ],
      "name": "allowance",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "approve",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "internalType": "address",
          "name": "account",
          "type": "address"
        }
      ],
      "name": "balanceOf",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "decimals",
      "outputs": [
        {
          "internalType": "uint8",
          "name": "",
          "type": "uint8"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "subtractedValue",
          "type": "uint256"
        }
      ],
      "name": "decreaseAllowance",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "address",
          "name": "spender",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "addedValue",
          "type": "uint256"
        }
      ],
      "name": "increaseAllowance",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "name",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "symbol",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "totalSupply",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "address",
          "name": "recipient",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "transfer",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "address",
          "name": "sender",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "recipient",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "transferFrom",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ],
  "metadata": "{\"compiler\":{\"version\":\"0.5.17+commit.d19bba13\"},\"language\":\"Solidity\",\"output\":{\"abi\":[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}],\"devdoc\":{\"methods\":{\"allowance(address,address)\":{\"details\":\"See {IERC20-allowance}.\"},\"approve(address,uint256)\":{\"details\":\"See {IERC20-approve}.     * Requirements:     * - `spender` cannot be the zero address.\"},\"balanceOf(address)\":{\"details\":\"See {IERC20-balanceOf}.\"},\"decimals()\":{\"details\":\"Returns the number of decimals used to get its user representation. For example, if `decimals` equals `2`, a balance of `505` tokens should be displayed to a user as `5,05` (`505 / 10 ** 2`).     * Tokens usually opt for a value of 18, imitating the relationship between Ether and Wei.     * NOTE: This information is only used for _display_ purposes: it in no way affects any of the arithmetic of the contract, including {IERC20-balanceOf} and {IERC20-transfer}.\"},\"decreaseAllowance(address,uint256)\":{\"details\":\"Atomically decreases the allowance granted to `spender` by the caller.     * This is an alternative to {approve} that can be used as a mitigation for problems described in {IERC20-approve}.     * Emits an {Approval} event indicating the updated allowance.     * Requirements:     * - `spender` cannot be the zero address. - `spender` must have allowance for the caller of at least `subtractedValue`.\"},\"increaseAllowance(address,uint256)\":{\"details\":\"Atomically increases the allowance granted to `spender` by the caller.     * This is an alternative to {approve} that can be used as a mitigation for problems described in {IERC20-approve}.     * Emits an {Approval} event indicating the updated allowance.     * Requirements:     * - `spender` cannot be the zero address.\"},\"name()\":{\"details\":\"Returns the name of the token.\"},\"symbol()\":{\"details\":\"Returns the symbol of the token, usually a shorter version of the name.\"},\"totalSupply()\":{\"details\":\"See {IERC20-totalSupply}.\"},\"transfer(address,uint256)\":{\"details\":\"See {IERC20-transfer}.     * Requirements:     * - `recipient` cannot be the zero address. - the caller must have a balance of at least `amount`.\"},\"transferFrom(address,address,uint256)\":{\"details\":\"See {IERC20-transferFrom}.     * Emits an {Approval} event indicating the updated allowance. This is not required by the EIP. See the note at the beginning of {ERC20};     * Requirements: - `sender` and `recipient` cannot be the zero address. - `sender` must have a balance of at least `amount`. - the caller must have allowance for `sender`'s tokens of at least `amount`.\"}}},\"userdoc\":{\"methods\":{}}},\"settings\":{\"compilationTarget\":{\"/Users/mpierce/sif/sifnode/ui/chains/ethereum/contracts/BobToken.sol\":\"BobToken\"},\"evmVersion\":\"istanbul\",\"libraries\":{},\"optimizer\":{\"enabled\":false,\"runs\":200},\"remappings\":[]},\"sources\":{\"/Users/mpierce/sif/sifnode/ui/chains/ethereum/contracts/BobToken.sol\":{\"keccak256\":\"0x55b9ad223d8d621ee64d02f0915310671eeca574a6bddd8b1b3e6f2631dfe554\",\"urls\":[\"bzz-raw://56b8e525a766f063b95874d76e7fb924f9f386c4aa8023d499ead2c4b17be00e\",\"dweb:/ipfs/QmWq64cBiZv9UMGgs1Bas63frXHBQixMizjLjXWc3by5Vt\"]},\"@openzeppelin/contracts/GSN/Context.sol\":{\"keccak256\":\"0x90a3995645af7562d84b9d69363ffa5ae7217714ab61e951bf7bc450f40e4061\",\"urls\":[\"bzz-raw://216ef9d6b614db4eb46970b4e84903f2534a45572dd30a79f0041f1a5830f436\",\"dweb:/ipfs/QmNPrJ4MWKUAWzKXpUqeyKRUfosaoANZAqXgvepdrCwZAG\"]},\"@openzeppelin/contracts/math/SafeMath.sol\":{\"keccak256\":\"0x640b6dee7a4b830bdfd52b5031a07fc2b12209f5b2e29e5d364a7d37f69d8076\",\"urls\":[\"bzz-raw://31113152e1ddb78fe7a4197f247591ca894e93f916867beb708d8e747b6cc74f\",\"dweb:/ipfs/QmbZaJyXdpsYGykVhHH9qpVGQg9DGCxE2QufbCUy3daTgq\"]},\"@openzeppelin/contracts/token/ERC20/ERC20.sol\":{\"keccak256\":\"0xb15af804e2bc97db51e4e103f13de9fe13f87e6b835d7a88c897966c0e58506e\",\"urls\":[\"bzz-raw://1e8cff8437557fc915a3bed968fcd8f2df9809599e665ef69c2c9ce628548055\",\"dweb:/ipfs/QmP5spYP8vs2jvLF8zNrXUbqB79hMsoEvMHiLcBxerWKcm\"]},\"@openzeppelin/contracts/token/ERC20/ERC20Detailed.sol\":{\"keccak256\":\"0x4a3a810b7ebe742e897e1fd428b3eeed2196d3acea58eaf9c566ed10d545d2ed\",\"urls\":[\"bzz-raw://28d03f328e19e4099d5013de0649afaeaabac1a614e130767ab1cb4aca9775b6\",\"dweb:/ipfs/Qmb9iW7yNuYehB2NfhRMs9TakqbLVQhBhmGMkaDZ5g1Eb4\"]},\"@openzeppelin/contracts/token/ERC20/IERC20.sol\":{\"keccak256\":\"0xe5bb0f57cff3e299f360052ba50f1ea0fff046df2be070b6943e0e3c3fdad8a9\",\"urls\":[\"bzz-raw://59fd025151435da35faa8093a5c7a17de02de9d08ad27275c5cdf05050820d91\",\"dweb:/ipfs/QmQMvwEcPhoRXzbXyrdoeRtvLoifUW9Qh7Luho7bmUPRkc\"]}},\"version\":1}",
  "bytecode": "0x60806040523480156200001157600080fd5b506040518060400160405280600881526020017f426f62546f6b656e0000000000000000000000000000000000000000000000008152506040518060400160405280600381526020017f42544b000000000000000000000000000000000000000000000000000000000081525060068260039080519060200190620000989291906200036b565b508160049080519060200190620000b19291906200036b565b5080600560006101000a81548160ff021916908360ff160217905550505050620000fb33620000e56200010160201b60201c565b60ff16600a0a612710026200011860201b60201c565b6200041a565b6000600560009054906101000a900460ff16905090565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415620001bc576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601f8152602001807f45524332303a206d696e7420746f20746865207a65726f20616464726573730081525060200191505060405180910390fd5b620001d881600254620002e260201b62000f221790919060201c565b60028190555062000236816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054620002e260201b62000f221790919060201c565b6000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040518082815260200191505060405180910390a35050565b60008082840190508381101562000361576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b8091505092915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620003ae57805160ff1916838001178555620003df565b82800160010185558215620003df579182015b82811115620003de578251825591602001919060010190620003c1565b5b509050620003ee9190620003f2565b5090565b6200041791905b8082111562000413576000816000905550600101620003f9565b5090565b90565b6110e0806200042a6000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c80633950935111610071578063395093511461025f57806370a08231146102c557806395d89b411461031d578063a457c2d7146103a0578063a9059cbb14610406578063dd62ed3e1461046c576100a9565b806306fdde03146100ae578063095ea7b31461013157806318160ddd1461019757806323b872dd146101b5578063313ce5671461023b575b600080fd5b6100b66104e4565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100f65780820151818401526020810190506100db565b50505050905090810190601f1680156101235780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61017d6004803603604081101561014757600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610586565b604051808215151515815260200191505060405180910390f35b61019f6105a4565b6040518082815260200191505060405180910390f35b610221600480360360608110156101cb57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506105ae565b604051808215151515815260200191505060405180910390f35b610243610687565b604051808260ff1660ff16815260200191505060405180910390f35b6102ab6004803603604081101561027557600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919050505061069e565b604051808215151515815260200191505060405180910390f35b610307600480360360208110156102db57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610751565b6040518082815260200191505060405180910390f35b610325610799565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561036557808201518184015260208101905061034a565b50505050905090810190601f1680156103925780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6103ec600480360360408110156103b657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919050505061083b565b604051808215151515815260200191505060405180910390f35b6104526004803603604081101561041c57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610908565b604051808215151515815260200191505060405180910390f35b6104ce6004803603604081101561048257600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610926565b6040518082815260200191505060405180910390f35b606060038054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561057c5780601f106105515761010080835404028352916020019161057c565b820191906000526020600020905b81548152906001019060200180831161055f57829003601f168201915b5050505050905090565b600061059a6105936109ad565b84846109b5565b6001905092915050565b6000600254905090565b60006105bb848484610bac565b61067c846105c76109ad565b6106778560405180606001604052806028815260200161101660289139600160008b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600061062d6109ad565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610e629092919063ffffffff16565b6109b5565b600190509392505050565b6000600560009054906101000a900460ff16905090565b60006107476106ab6109ad565b8461074285600160006106bc6109ad565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610f2290919063ffffffff16565b6109b5565b6001905092915050565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b606060048054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156108315780601f1061080657610100808354040283529160200191610831565b820191906000526020600020905b81548152906001019060200180831161081457829003601f168201915b5050505050905090565b60006108fe6108486109ad565b846108f98560405180606001604052806025815260200161108760259139600160006108726109ad565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610e629092919063ffffffff16565b6109b5565b6001905092915050565b600061091c6109156109ad565b8484610bac565b6001905092915050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610a3b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260248152602001806110636024913960400191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610ac1576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526022815260200180610fce6022913960400191505060405180910390fd5b80600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925836040518082815260200191505060405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610c32576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602581526020018061103e6025913960400191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610cb8576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526023815260200180610fab6023913960400191505060405180910390fd5b610d2381604051806060016040528060268152602001610ff0602691396000808773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610e629092919063ffffffff16565b6000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610db6816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610f2290919063ffffffff16565b6000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040518082815260200191505060405180910390a3505050565b6000838311158290610f0f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b83811015610ed4578082015181840152602081019050610eb9565b50505050905090810190601f168015610f015780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b5060008385039050809150509392505050565b600080828401905083811015610fa0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b809150509291505056fe45524332303a207472616e7366657220746f20746865207a65726f206164647265737345524332303a20617070726f766520746f20746865207a65726f206164647265737345524332303a207472616e7366657220616d6f756e7420657863656564732062616c616e636545524332303a207472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636545524332303a207472616e736665722066726f6d20746865207a65726f206164647265737345524332303a20617070726f76652066726f6d20746865207a65726f206164647265737345524332303a2064656372656173656420616c6c6f77616e63652062656c6f77207a65726fa265627a7a72315820ac471f30ed2abf806e248c9209741abb03b630552969dd83dc46f78fcc4e103264736f6c63430005110032",
  "deployedBytecode": "0x608060405234801561001057600080fd5b50600436106100a95760003560e01c80633950935111610071578063395093511461025f57806370a08231146102c557806395d89b411461031d578063a457c2d7146103a0578063a9059cbb14610406578063dd62ed3e1461046c576100a9565b806306fdde03146100ae578063095ea7b31461013157806318160ddd1461019757806323b872dd146101b5578063313ce5671461023b575b600080fd5b6100b66104e4565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100f65780820151818401526020810190506100db565b50505050905090810190601f1680156101235780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b61017d6004803603604081101561014757600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610586565b604051808215151515815260200191505060405180910390f35b61019f6105a4565b6040518082815260200191505060405180910390f35b610221600480360360608110156101cb57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291905050506105ae565b604051808215151515815260200191505060405180910390f35b610243610687565b604051808260ff1660ff16815260200191505060405180910390f35b6102ab6004803603604081101561027557600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919050505061069e565b604051808215151515815260200191505060405180910390f35b610307600480360360208110156102db57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610751565b6040518082815260200191505060405180910390f35b610325610799565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561036557808201518184015260208101905061034a565b50505050905090810190601f1680156103925780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6103ec600480360360408110156103b657600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919050505061083b565b604051808215151515815260200191505060405180910390f35b6104526004803603604081101561041c57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050610908565b604051808215151515815260200191505060405180910390f35b6104ce6004803603604081101561048257600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610926565b6040518082815260200191505060405180910390f35b606060038054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561057c5780601f106105515761010080835404028352916020019161057c565b820191906000526020600020905b81548152906001019060200180831161055f57829003601f168201915b5050505050905090565b600061059a6105936109ad565b84846109b5565b6001905092915050565b6000600254905090565b60006105bb848484610bac565b61067c846105c76109ad565b6106778560405180606001604052806028815260200161101660289139600160008b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600061062d6109ad565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610e629092919063ffffffff16565b6109b5565b600190509392505050565b6000600560009054906101000a900460ff16905090565b60006107476106ab6109ad565b8461074285600160006106bc6109ad565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008973ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610f2290919063ffffffff16565b6109b5565b6001905092915050565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b606060048054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156108315780601f1061080657610100808354040283529160200191610831565b820191906000526020600020905b81548152906001019060200180831161081457829003601f168201915b5050505050905090565b60006108fe6108486109ad565b846108f98560405180606001604052806025815260200161108760259139600160006108726109ad565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610e629092919063ffffffff16565b6109b5565b6001905092915050565b600061091c6109156109ad565b8484610bac565b6001905092915050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610a3b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260248152602001806110636024913960400191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610ac1576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526022815260200180610fce6022913960400191505060405180910390fd5b80600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925836040518082815260200191505060405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff161415610c32576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602581526020018061103e6025913960400191505060405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff161415610cb8576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401808060200182810382526023815260200180610fab6023913960400191505060405180910390fd5b610d2381604051806060016040528060268152602001610ff0602691396000808773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610e629092919063ffffffff16565b6000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550610db6816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054610f2290919063ffffffff16565b6000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef836040518082815260200191505060405180910390a3505050565b6000838311158290610f0f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825283818151815260200191508051906020019080838360005b83811015610ed4578082015181840152602081019050610eb9565b50505050905090810190601f168015610f015780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b5060008385039050809150509392505050565b600080828401905083811015610fa0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601b8152602001807f536166654d6174683a206164646974696f6e206f766572666c6f77000000000081525060200191505060405180910390fd5b809150509291505056fe45524332303a207472616e7366657220746f20746865207a65726f206164647265737345524332303a20617070726f766520746f20746865207a65726f206164647265737345524332303a207472616e7366657220616d6f756e7420657863656564732062616c616e636545524332303a207472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636545524332303a207472616e736665722066726f6d20746865207a65726f206164647265737345524332303a20617070726f76652066726f6d20746865207a65726f206164647265737345524332303a2064656372656173656420616c6c6f77616e63652062656c6f77207a65726fa265627a7a72315820ac471f30ed2abf806e248c9209741abb03b630552969dd83dc46f78fcc4e103264736f6c63430005110032",
  "sourceMap": "155:176:1:-;;;203:126;8:9:-1;5:2;;;30:1;27;20:12;5:2;203:126:1;416:163:6;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;257:1:1;512:4:6;504:5;:12;;;;;;;;;;;;:::i;:::-;;536:6;526:7;:16;;;;;;;;;;;;:::i;:::-;;564:8;552:9;;:20;;;;;;;;;;;;;;;;;;416:163;;;270:52:1;276:10;309;:8;;;:10;;:::i;:::-;301:19;;297:2;:23;288:5;:33;270:5;;;:52;;:::i;:::-;155:176;;1472:81:6;1513:5;1537:9;;;;;;;;;;;1530:16;;1472:81;:::o;5962:302:5:-;6056:1;6037:21;;:7;:21;;;;6029:65;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;6120:24;6137:6;6120:12;;:16;;;;;;:24;;;;:::i;:::-;6105:12;:39;;;;6175:30;6198:6;6175:9;:18;6185:7;6175:18;;;;;;;;;;;;;;;;:22;;;;;;:30;;;;:::i;:::-;6154:9;:18;6164:7;6154:18;;;;;;;;;;;;;;;:51;;;;6241:7;6220:37;;6237:1;6220:37;;;6250:6;6220:37;;;;;;;;;;;;;;;;;;5962:302;;:::o;834:176:4:-;892:7;911:9;927:1;923;:5;911:17;;951:1;946;:6;;938:46;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;1002:1;995:8;;;834:176;;;;:::o;155::1:-;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;:::i;:::-;;;:::o;:::-;;;;;;;;;;;;;;;;;;;;;;;;;;;:::o;:::-;;;;;;;",
  "deployedSourceMap": "155:176:1:-;;;;8:9:-1;5:2;;;30:1;27;20:12;5:2;155:176:1;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;644:81:6;;;:::i;:::-;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;23:1:-1;8:100;33:3;30:1;27:10;8:100;;;99:1;94:3;90:11;84:18;80:1;75:3;71:11;64:39;52:2;49:1;45:10;40:15;;8:100;;;12:14;644:81:6;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;2500:149:5;;;;;;13:2:-1;8:3;5:11;2:2;;;29:1;26;19:12;2:2;2500:149:5;;;;;;;;;;;;;;;;;;;;;;;;;;;;:::i;:::-;;;;;;;;;;;;;;;;;;;;;;;1559:89;;;:::i;:::-;;;;;;;;;;;;;;;;;;;3107:300;;;;;;13:2:-1;8:3;5:11;2:2;;;29:1;26;19:12;2:2;3107:300:5;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;:::i;:::-;;;;;;;;;;;;;;;;;;;;;;;1472:81:6;;;:::i;:::-;;;;;;;;;;;;;;;;;;;;;;;3802:207:5;;;;;;13:2:-1;8:3;5:11;2:2;;;29:1;26;19:12;2:2;3802:207:5;;;;;;;;;;;;;;;;;;;;;;;;;;;;:::i;:::-;;;;;;;;;;;;;;;;;;;;;;;1706:108;;;;;;13:2:-1;8:3;5:11;2:2;;;29:1;26;19:12;2:2;1706:108:5;;;;;;;;;;;;;;;;;;;:::i;:::-;;;;;;;;;;;;;;;;;;;838:85:6;;;:::i;:::-;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;23:1:-1;8:100;33:3;30:1;27:10;8:100;;;99:1;94:3;90:11;84:18;80:1;75:3;71:11;64:39;52:2;49:1;45:10;40:15;;8:100;;;12:14;838:85:6;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;4496:258:5;;;;;;13:2:-1;8:3;5:11;2:2;;;29:1;26;19:12;2:2;4496:258:5;;;;;;;;;;;;;;;;;;;;;;;;;;;;:::i;:::-;;;;;;;;;;;;;;;;;;;;;;;2017:155;;;;;;13:2:-1;8:3;5:11;2:2;;;29:1;26;19:12;2:2;2017:155:5;;;;;;;;;;;;;;;;;;;;;;;;;;;;:::i;:::-;;;;;;;;;;;;;;;;;;;;;;;2230:132;;;;;;13:2:-1;8:3;5:11;2:2;;;29:1;26;19:12;2:2;2230:132:5;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;:::i;:::-;;;;;;;;;;;;;;;;;;;644:81:6;681:13;713:5;706:12;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;644:81;:::o;2500:149:5:-;2566:4;2582:39;2591:12;:10;:12::i;:::-;2605:7;2614:6;2582:8;:39::i;:::-;2638:4;2631:11;;2500:149;;;;:::o;1559:89::-;1603:7;1629:12;;1622:19;;1559:89;:::o;3107:300::-;3196:4;3212:36;3222:6;3230:9;3241:6;3212:9;:36::i;:::-;3258:121;3267:6;3275:12;:10;:12::i;:::-;3289:89;3327:6;3289:89;;;;;;;;;;;;;;;;;:11;:19;3301:6;3289:19;;;;;;;;;;;;;;;:33;3309:12;:10;:12::i;:::-;3289:33;;;;;;;;;;;;;;;;:37;;:89;;;;;:::i;:::-;3258:8;:121::i;:::-;3396:4;3389:11;;3107:300;;;;;:::o;1472:81:6:-;1513:5;1537:9;;;;;;;;;;;1530:16;;1472:81;:::o;3802:207:5:-;3882:4;3898:83;3907:12;:10;:12::i;:::-;3921:7;3930:50;3969:10;3930:11;:25;3942:12;:10;:12::i;:::-;3930:25;;;;;;;;;;;;;;;:34;3956:7;3930:34;;;;;;;;;;;;;;;;:38;;:50;;;;:::i;:::-;3898:8;:83::i;:::-;3998:4;3991:11;;3802:207;;;;:::o;1706:108::-;1763:7;1789:9;:18;1799:7;1789:18;;;;;;;;;;;;;;;;1782:25;;1706:108;;;:::o;838:85:6:-;877:13;909:7;902:14;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;838:85;:::o;4496:258:5:-;4581:4;4597:129;4606:12;:10;:12::i;:::-;4620:7;4629:96;4668:15;4629:96;;;;;;;;;;;;;;;;;:11;:25;4641:12;:10;:12::i;:::-;4629:25;;;;;;;;;;;;;;;:34;4655:7;4629:34;;;;;;;;;;;;;;;;:38;;:96;;;;;:::i;:::-;4597:8;:129::i;:::-;4743:4;4736:11;;4496:258;;;;:::o;2017:155::-;2086:4;2102:42;2112:12;:10;:12::i;:::-;2126:9;2137:6;2102:9;:42::i;:::-;2161:4;2154:11;;2017:155;;;;:::o;2230:132::-;2302:7;2328:11;:18;2340:5;2328:18;;;;;;;;;;;;;;;:27;2347:7;2328:27;;;;;;;;;;;;;;;;2321:34;;2230:132;;;;:::o;788:96:3:-;833:15;867:10;860:17;;788:96;:::o;7350:332:5:-;7460:1;7443:19;;:5;:19;;;;7435:68;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;7540:1;7521:21;;:7;:21;;;;7513:68;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;7622:6;7592:11;:18;7604:5;7592:18;;;;;;;;;;;;;;;:27;7611:7;7592:27;;;;;;;;;;;;;;;:36;;;;7659:7;7643:32;;7652:5;7643:32;;;7668:6;7643:32;;;;;;;;;;;;;;;;;;7350:332;;;:::o;5228:464::-;5343:1;5325:20;;:6;:20;;;;5317:70;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;5426:1;5405:23;;:9;:23;;;;5397:71;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;5499;5521:6;5499:71;;;;;;;;;;;;;;;;;:9;:17;5509:6;5499:17;;;;;;;;;;;;;;;;:21;;:71;;;;;:::i;:::-;5479:9;:17;5489:6;5479:17;;;;;;;;;;;;;;;:91;;;;5603:32;5628:6;5603:9;:20;5613:9;5603:20;;;;;;;;;;;;;;;;:24;;:32;;;;:::i;:::-;5580:9;:20;5590:9;5580:20;;;;;;;;;;;;;;;:55;;;;5667:9;5650:35;;5659:6;5650:35;;;5678:6;5650:35;;;;;;;;;;;;;;;;;;5228:464;;;:::o;1732:187:4:-;1818:7;1850:1;1845;:6;;1853:12;1837:29;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;23:1:-1;8:100;33:3;30:1;27:10;8:100;;;99:1;94:3;90:11;84:18;80:1;75:3;71:11;64:39;52:2;49:1;45:10;40:15;;8:100;;;12:14;1837:29:4;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;1876:9;1892:1;1888;:5;1876:17;;1911:1;1904:8;;;1732:187;;;;;:::o;834:176::-;892:7;911:9;927:1;923;:5;911:17;;951:1;946;:6;;938:46;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;1002:1;995:8;;;834:176;;;;:::o",
  "source": "pragma solidity >=0.4.22 <0.8.0;\n\nimport \"@openzeppelin/contracts/token/ERC20/ERC20.sol\";\nimport \"@openzeppelin/contracts/token/ERC20/ERC20Detailed.sol\";\n\ncontract BobToken is ERC20, ERC20Detailed {\n    constructor() public ERC20Detailed(\"BobToken\", \"BTK\", 6) {\n        _mint(msg.sender, 10000 * (10**uint256(decimals())));\n    }\n}\n",
  "sourcePath": "/Users/mpierce/sif/sifnode/ui/chains/ethereum/contracts/BobToken.sol",
  "ast": {
    "absolutePath": "/Users/mpierce/sif/sifnode/ui/chains/ethereum/contracts/BobToken.sol",
    "exportedSymbols": {
      "BobToken": [
        63
      ]
    },
    "id": 64,
    "nodeType": "SourceUnit",
    "nodes": [
      {
        "id": 33,
        "literals": [
          "solidity",
          ">=",
          "0.4",
          ".22",
          "<",
          "0.8",
          ".0"
        ],
        "nodeType": "PragmaDirective",
        "src": "0:32:1"
      },
      {
        "absolutePath": "@openzeppelin/contracts/token/ERC20/ERC20.sol",
        "file": "@openzeppelin/contracts/token/ERC20/ERC20.sol",
        "id": 34,
        "nodeType": "ImportDirective",
        "scope": 64,
        "sourceUnit": 716,
        "src": "34:55:1",
        "symbolAliases": [],
        "unitAlias": ""
      },
      {
        "absolutePath": "@openzeppelin/contracts/token/ERC20/ERC20Detailed.sol",
        "file": "@openzeppelin/contracts/token/ERC20/ERC20Detailed.sol",
        "id": 35,
        "nodeType": "ImportDirective",
        "scope": 64,
        "sourceUnit": 774,
        "src": "90:63:1",
        "symbolAliases": [],
        "unitAlias": ""
      },
      {
        "baseContracts": [
          {
            "arguments": null,
            "baseName": {
              "contractScope": null,
              "id": 36,
              "name": "ERC20",
              "nodeType": "UserDefinedTypeName",
              "referencedDeclaration": 715,
              "src": "176:5:1",
              "typeDescriptions": {
                "typeIdentifier": "t_contract$_ERC20_$715",
                "typeString": "contract ERC20"
              }
            },
            "id": 37,
            "nodeType": "InheritanceSpecifier",
            "src": "176:5:1"
          },
          {
            "arguments": null,
            "baseName": {
              "contractScope": null,
              "id": 38,
              "name": "ERC20Detailed",
              "nodeType": "UserDefinedTypeName",
              "referencedDeclaration": 773,
              "src": "183:13:1",
              "typeDescriptions": {
                "typeIdentifier": "t_contract$_ERC20Detailed_$773",
                "typeString": "contract ERC20Detailed"
              }
            },
            "id": 39,
            "nodeType": "InheritanceSpecifier",
            "src": "183:13:1"
          }
        ],
        "contractDependencies": [
          123,
          715,
          773,
          842
        ],
        "contractKind": "contract",
        "documentation": null,
        "fullyImplemented": true,
        "id": 63,
        "linearizedBaseContracts": [
          63,
          773,
          715,
          842,
          123
        ],
        "name": "BobToken",
        "nodeType": "ContractDefinition",
        "nodes": [
          {
            "body": {
              "id": 61,
              "nodeType": "Block",
              "src": "260:69:1",
              "statements": [
                {
                  "expression": {
                    "argumentTypes": null,
                    "arguments": [
                      {
                        "argumentTypes": null,
                        "expression": {
                          "argumentTypes": null,
                          "id": 48,
                          "name": "msg",
                          "nodeType": "Identifier",
                          "overloadedDeclarations": [],
                          "referencedDeclaration": 857,
                          "src": "276:3:1",
                          "typeDescriptions": {
                            "typeIdentifier": "t_magic_message",
                            "typeString": "msg"
                          }
                        },
                        "id": 49,
                        "isConstant": false,
                        "isLValue": false,
                        "isPure": false,
                        "lValueRequested": false,
                        "memberName": "sender",
                        "nodeType": "MemberAccess",
                        "referencedDeclaration": null,
                        "src": "276:10:1",
                        "typeDescriptions": {
                          "typeIdentifier": "t_address_payable",
                          "typeString": "address payable"
                        }
                      },
                      {
                        "argumentTypes": null,
                        "commonType": {
                          "typeIdentifier": "t_uint256",
                          "typeString": "uint256"
                        },
                        "id": 58,
                        "isConstant": false,
                        "isLValue": false,
                        "isPure": false,
                        "lValueRequested": false,
                        "leftExpression": {
                          "argumentTypes": null,
                          "hexValue": "3130303030",
                          "id": 50,
                          "isConstant": false,
                          "isLValue": false,
                          "isPure": true,
                          "kind": "number",
                          "lValueRequested": false,
                          "nodeType": "Literal",
                          "src": "288:5:1",
                          "subdenomination": null,
                          "typeDescriptions": {
                            "typeIdentifier": "t_rational_10000_by_1",
                            "typeString": "int_const 10000"
                          },
                          "value": "10000"
                        },
                        "nodeType": "BinaryOperation",
                        "operator": "*",
                        "rightExpression": {
                          "argumentTypes": null,
                          "components": [
                            {
                              "argumentTypes": null,
                              "commonType": {
                                "typeIdentifier": "t_uint256",
                                "typeString": "uint256"
                              },
                              "id": 56,
                              "isConstant": false,
                              "isLValue": false,
                              "isPure": false,
                              "lValueRequested": false,
                              "leftExpression": {
                                "argumentTypes": null,
                                "hexValue": "3130",
                                "id": 51,
                                "isConstant": false,
                                "isLValue": false,
                                "isPure": true,
                                "kind": "number",
                                "lValueRequested": false,
                                "nodeType": "Literal",
                                "src": "297:2:1",
                                "subdenomination": null,
                                "typeDescriptions": {
                                  "typeIdentifier": "t_rational_10_by_1",
                                  "typeString": "int_const 10"
                                },
                                "value": "10"
                              },
                              "nodeType": "BinaryOperation",
                              "operator": "**",
                              "rightExpression": {
                                "argumentTypes": null,
                                "arguments": [
                                  {
                                    "argumentTypes": null,
                                    "arguments": [],
                                    "expression": {
                                      "argumentTypes": [],
                                      "id": 53,
                                      "name": "decimals",
                                      "nodeType": "Identifier",
                                      "overloadedDeclarations": [],
                                      "referencedDeclaration": 772,
                                      "src": "309:8:1",
                                      "typeDescriptions": {
                                        "typeIdentifier": "t_function_internal_view$__$returns$_t_uint8_$",
                                        "typeString": "function () view returns (uint8)"
                                      }
                                    },
                                    "id": 54,
                                    "isConstant": false,
                                    "isLValue": false,
                                    "isPure": false,
                                    "kind": "functionCall",
                                    "lValueRequested": false,
                                    "names": [],
                                    "nodeType": "FunctionCall",
                                    "src": "309:10:1",
                                    "typeDescriptions": {
                                      "typeIdentifier": "t_uint8",
                                      "typeString": "uint8"
                                    }
                                  }
                                ],
                                "expression": {
                                  "argumentTypes": [
                                    {
                                      "typeIdentifier": "t_uint8",
                                      "typeString": "uint8"
                                    }
                                  ],
                                  "id": 52,
                                  "isConstant": false,
                                  "isLValue": false,
                                  "isPure": true,
                                  "lValueRequested": false,
                                  "nodeType": "ElementaryTypeNameExpression",
                                  "src": "301:7:1",
                                  "typeDescriptions": {
                                    "typeIdentifier": "t_type$_t_uint256_$",
                                    "typeString": "type(uint256)"
                                  },
                                  "typeName": "uint256"
                                },
                                "id": 55,
                                "isConstant": false,
                                "isLValue": false,
                                "isPure": false,
                                "kind": "typeConversion",
                                "lValueRequested": false,
                                "names": [],
                                "nodeType": "FunctionCall",
                                "src": "301:19:1",
                                "typeDescriptions": {
                                  "typeIdentifier": "t_uint256",
                                  "typeString": "uint256"
                                }
                              },
                              "src": "297:23:1",
                              "typeDescriptions": {
                                "typeIdentifier": "t_uint256",
                                "typeString": "uint256"
                              }
                            }
                          ],
                          "id": 57,
                          "isConstant": false,
                          "isInlineArray": false,
                          "isLValue": false,
                          "isPure": false,
                          "lValueRequested": false,
                          "nodeType": "TupleExpression",
                          "src": "296:25:1",
                          "typeDescriptions": {
                            "typeIdentifier": "t_uint256",
                            "typeString": "uint256"
                          }
                        },
                        "src": "288:33:1",
                        "typeDescriptions": {
                          "typeIdentifier": "t_uint256",
                          "typeString": "uint256"
                        }
                      }
                    ],
                    "expression": {
                      "argumentTypes": [
                        {
                          "typeIdentifier": "t_address_payable",
                          "typeString": "address payable"
                        },
                        {
                          "typeIdentifier": "t_uint256",
                          "typeString": "uint256"
                        }
                      ],
                      "id": 47,
                      "name": "_mint",
                      "nodeType": "Identifier",
                      "overloadedDeclarations": [],
                      "referencedDeclaration": 599,
                      "src": "270:5:1",
                      "typeDescriptions": {
                        "typeIdentifier": "t_function_internal_nonpayable$_t_address_$_t_uint256_$returns$__$",
                        "typeString": "function (address,uint256)"
                      }
                    },
                    "id": 59,
                    "isConstant": false,
                    "isLValue": false,
                    "isPure": false,
                    "kind": "functionCall",
                    "lValueRequested": false,
                    "names": [],
                    "nodeType": "FunctionCall",
                    "src": "270:52:1",
                    "typeDescriptions": {
                      "typeIdentifier": "t_tuple$__$",
                      "typeString": "tuple()"
                    }
                  },
                  "id": 60,
                  "nodeType": "ExpressionStatement",
                  "src": "270:52:1"
                }
              ]
            },
            "documentation": null,
            "id": 62,
            "implemented": true,
            "kind": "constructor",
            "modifiers": [
              {
                "arguments": [
                  {
                    "argumentTypes": null,
                    "hexValue": "426f62546f6b656e",
                    "id": 42,
                    "isConstant": false,
                    "isLValue": false,
                    "isPure": true,
                    "kind": "string",
                    "lValueRequested": false,
                    "nodeType": "Literal",
                    "src": "238:10:1",
                    "subdenomination": null,
                    "typeDescriptions": {
                      "typeIdentifier": "t_stringliteral_6961ba82ab463124347186eaacd7907445b2abce40f118cc3374a5f8fd621cd7",
                      "typeString": "literal_string \"BobToken\""
                    },
                    "value": "BobToken"
                  },
                  {
                    "argumentTypes": null,
                    "hexValue": "42544b",
                    "id": 43,
                    "isConstant": false,
                    "isLValue": false,
                    "isPure": true,
                    "kind": "string",
                    "lValueRequested": false,
                    "nodeType": "Literal",
                    "src": "250:5:1",
                    "subdenomination": null,
                    "typeDescriptions": {
                      "typeIdentifier": "t_stringliteral_cbd09d845af28ee5b0642a9230e90a4426144cbbda10cea126e4301fa7e1a7f0",
                      "typeString": "literal_string \"BTK\""
                    },
                    "value": "BTK"
                  },
                  {
                    "argumentTypes": null,
                    "hexValue": "36",
                    "id": 44,
                    "isConstant": false,
                    "isLValue": false,
                    "isPure": true,
                    "kind": "number",
                    "lValueRequested": false,
                    "nodeType": "Literal",
                    "src": "257:1:1",
                    "subdenomination": null,
                    "typeDescriptions": {
                      "typeIdentifier": "t_rational_6_by_1",
                      "typeString": "int_const 6"
                    },
                    "value": "6"
                  }
                ],
                "id": 45,
                "modifierName": {
                  "argumentTypes": null,
                  "id": 41,
                  "name": "ERC20Detailed",
                  "nodeType": "Identifier",
                  "overloadedDeclarations": [],
                  "referencedDeclaration": 773,
                  "src": "224:13:1",
                  "typeDescriptions": {
                    "typeIdentifier": "t_type$_t_contract$_ERC20Detailed_$773_$",
                    "typeString": "type(contract ERC20Detailed)"
                  }
                },
                "nodeType": "ModifierInvocation",
                "src": "224:35:1"
              }
            ],
            "name": "",
            "nodeType": "FunctionDefinition",
            "parameters": {
              "id": 40,
              "nodeType": "ParameterList",
              "parameters": [],
              "src": "214:2:1"
            },
            "returnParameters": {
              "id": 46,
              "nodeType": "ParameterList",
              "parameters": [],
              "src": "260:0:1"
            },
            "scope": 63,
            "src": "203:126:1",
            "stateMutability": "nonpayable",
            "superFunction": null,
            "visibility": "public"
          }
        ],
        "scope": 64,
        "src": "155:176:1"
      }
    ],
    "src": "0:332:1"
  },
  "legacyAST": {
    "absolutePath": "/Users/mpierce/sif/sifnode/ui/chains/ethereum/contracts/BobToken.sol",
    "exportedSymbols": {
      "BobToken": [
        63
      ]
    },
    "id": 64,
    "nodeType": "SourceUnit",
    "nodes": [
      {
        "id": 33,
        "literals": [
          "solidity",
          ">=",
          "0.4",
          ".22",
          "<",
          "0.8",
          ".0"
        ],
        "nodeType": "PragmaDirective",
        "src": "0:32:1"
      },
      {
        "absolutePath": "@openzeppelin/contracts/token/ERC20/ERC20.sol",
        "file": "@openzeppelin/contracts/token/ERC20/ERC20.sol",
        "id": 34,
        "nodeType": "ImportDirective",
        "scope": 64,
        "sourceUnit": 716,
        "src": "34:55:1",
        "symbolAliases": [],
        "unitAlias": ""
      },
      {
        "absolutePath": "@openzeppelin/contracts/token/ERC20/ERC20Detailed.sol",
        "file": "@openzeppelin/contracts/token/ERC20/ERC20Detailed.sol",
        "id": 35,
        "nodeType": "ImportDirective",
        "scope": 64,
        "sourceUnit": 774,
        "src": "90:63:1",
        "symbolAliases": [],
        "unitAlias": ""
      },
      {
        "baseContracts": [
          {
            "arguments": null,
            "baseName": {
              "contractScope": null,
              "id": 36,
              "name": "ERC20",
              "nodeType": "UserDefinedTypeName",
              "referencedDeclaration": 715,
              "src": "176:5:1",
              "typeDescriptions": {
                "typeIdentifier": "t_contract$_ERC20_$715",
                "typeString": "contract ERC20"
              }
            },
            "id": 37,
            "nodeType": "InheritanceSpecifier",
            "src": "176:5:1"
          },
          {
            "arguments": null,
            "baseName": {
              "contractScope": null,
              "id": 38,
              "name": "ERC20Detailed",
              "nodeType": "UserDefinedTypeName",
              "referencedDeclaration": 773,
              "src": "183:13:1",
              "typeDescriptions": {
                "typeIdentifier": "t_contract$_ERC20Detailed_$773",
                "typeString": "contract ERC20Detailed"
              }
            },
            "id": 39,
            "nodeType": "InheritanceSpecifier",
            "src": "183:13:1"
          }
        ],
        "contractDependencies": [
          123,
          715,
          773,
          842
        ],
        "contractKind": "contract",
        "documentation": null,
        "fullyImplemented": true,
        "id": 63,
        "linearizedBaseContracts": [
          63,
          773,
          715,
          842,
          123
        ],
        "name": "BobToken",
        "nodeType": "ContractDefinition",
        "nodes": [
          {
            "body": {
              "id": 61,
              "nodeType": "Block",
              "src": "260:69:1",
              "statements": [
                {
                  "expression": {
                    "argumentTypes": null,
                    "arguments": [
                      {
                        "argumentTypes": null,
                        "expression": {
                          "argumentTypes": null,
                          "id": 48,
                          "name": "msg",
                          "nodeType": "Identifier",
                          "overloadedDeclarations": [],
                          "referencedDeclaration": 857,
                          "src": "276:3:1",
                          "typeDescriptions": {
                            "typeIdentifier": "t_magic_message",
                            "typeString": "msg"
                          }
                        },
                        "id": 49,
                        "isConstant": false,
                        "isLValue": false,
                        "isPure": false,
                        "lValueRequested": false,
                        "memberName": "sender",
                        "nodeType": "MemberAccess",
                        "referencedDeclaration": null,
                        "src": "276:10:1",
                        "typeDescriptions": {
                          "typeIdentifier": "t_address_payable",
                          "typeString": "address payable"
                        }
                      },
                      {
                        "argumentTypes": null,
                        "commonType": {
                          "typeIdentifier": "t_uint256",
                          "typeString": "uint256"
                        },
                        "id": 58,
                        "isConstant": false,
                        "isLValue": false,
                        "isPure": false,
                        "lValueRequested": false,
                        "leftExpression": {
                          "argumentTypes": null,
                          "hexValue": "3130303030",
                          "id": 50,
                          "isConstant": false,
                          "isLValue": false,
                          "isPure": true,
                          "kind": "number",
                          "lValueRequested": false,
                          "nodeType": "Literal",
                          "src": "288:5:1",
                          "subdenomination": null,
                          "typeDescriptions": {
                            "typeIdentifier": "t_rational_10000_by_1",
                            "typeString": "int_const 10000"
                          },
                          "value": "10000"
                        },
                        "nodeType": "BinaryOperation",
                        "operator": "*",
                        "rightExpression": {
                          "argumentTypes": null,
                          "components": [
                            {
                              "argumentTypes": null,
                              "commonType": {
                                "typeIdentifier": "t_uint256",
                                "typeString": "uint256"
                              },
                              "id": 56,
                              "isConstant": false,
                              "isLValue": false,
                              "isPure": false,
                              "lValueRequested": false,
                              "leftExpression": {
                                "argumentTypes": null,
                                "hexValue": "3130",
                                "id": 51,
                                "isConstant": false,
                                "isLValue": false,
                                "isPure": true,
                                "kind": "number",
                                "lValueRequested": false,
                                "nodeType": "Literal",
                                "src": "297:2:1",
                                "subdenomination": null,
                                "typeDescriptions": {
                                  "typeIdentifier": "t_rational_10_by_1",
                                  "typeString": "int_const 10"
                                },
                                "value": "10"
                              },
                              "nodeType": "BinaryOperation",
                              "operator": "**",
                              "rightExpression": {
                                "argumentTypes": null,
                                "arguments": [
                                  {
                                    "argumentTypes": null,
                                    "arguments": [],
                                    "expression": {
                                      "argumentTypes": [],
                                      "id": 53,
                                      "name": "decimals",
                                      "nodeType": "Identifier",
                                      "overloadedDeclarations": [],
                                      "referencedDeclaration": 772,
                                      "src": "309:8:1",
                                      "typeDescriptions": {
                                        "typeIdentifier": "t_function_internal_view$__$returns$_t_uint8_$",
                                        "typeString": "function () view returns (uint8)"
                                      }
                                    },
                                    "id": 54,
                                    "isConstant": false,
                                    "isLValue": false,
                                    "isPure": false,
                                    "kind": "functionCall",
                                    "lValueRequested": false,
                                    "names": [],
                                    "nodeType": "FunctionCall",
                                    "src": "309:10:1",
                                    "typeDescriptions": {
                                      "typeIdentifier": "t_uint8",
                                      "typeString": "uint8"
                                    }
                                  }
                                ],
                                "expression": {
                                  "argumentTypes": [
                                    {
                                      "typeIdentifier": "t_uint8",
                                      "typeString": "uint8"
                                    }
                                  ],
                                  "id": 52,
                                  "isConstant": false,
                                  "isLValue": false,
                                  "isPure": true,
                                  "lValueRequested": false,
                                  "nodeType": "ElementaryTypeNameExpression",
                                  "src": "301:7:1",
                                  "typeDescriptions": {
                                    "typeIdentifier": "t_type$_t_uint256_$",
                                    "typeString": "type(uint256)"
                                  },
                                  "typeName": "uint256"
                                },
                                "id": 55,
                                "isConstant": false,
                                "isLValue": false,
                                "isPure": false,
                                "kind": "typeConversion",
                                "lValueRequested": false,
                                "names": [],
                                "nodeType": "FunctionCall",
                                "src": "301:19:1",
                                "typeDescriptions": {
                                  "typeIdentifier": "t_uint256",
                                  "typeString": "uint256"
                                }
                              },
                              "src": "297:23:1",
                              "typeDescriptions": {
                                "typeIdentifier": "t_uint256",
                                "typeString": "uint256"
                              }
                            }
                          ],
                          "id": 57,
                          "isConstant": false,
                          "isInlineArray": false,
                          "isLValue": false,
                          "isPure": false,
                          "lValueRequested": false,
                          "nodeType": "TupleExpression",
                          "src": "296:25:1",
                          "typeDescriptions": {
                            "typeIdentifier": "t_uint256",
                            "typeString": "uint256"
                          }
                        },
                        "src": "288:33:1",
                        "typeDescriptions": {
                          "typeIdentifier": "t_uint256",
                          "typeString": "uint256"
                        }
                      }
                    ],
                    "expression": {
                      "argumentTypes": [
                        {
                          "typeIdentifier": "t_address_payable",
                          "typeString": "address payable"
                        },
                        {
                          "typeIdentifier": "t_uint256",
                          "typeString": "uint256"
                        }
                      ],
                      "id": 47,
                      "name": "_mint",
                      "nodeType": "Identifier",
                      "overloadedDeclarations": [],
                      "referencedDeclaration": 599,
                      "src": "270:5:1",
                      "typeDescriptions": {
                        "typeIdentifier": "t_function_internal_nonpayable$_t_address_$_t_uint256_$returns$__$",
                        "typeString": "function (address,uint256)"
                      }
                    },
                    "id": 59,
                    "isConstant": false,
                    "isLValue": false,
                    "isPure": false,
                    "kind": "functionCall",
                    "lValueRequested": false,
                    "names": [],
                    "nodeType": "FunctionCall",
                    "src": "270:52:1",
                    "typeDescriptions": {
                      "typeIdentifier": "t_tuple$__$",
                      "typeString": "tuple()"
                    }
                  },
                  "id": 60,
                  "nodeType": "ExpressionStatement",
                  "src": "270:52:1"
                }
              ]
            },
            "documentation": null,
            "id": 62,
            "implemented": true,
            "kind": "constructor",
            "modifiers": [
              {
                "arguments": [
                  {
                    "argumentTypes": null,
                    "hexValue": "426f62546f6b656e",
                    "id": 42,
                    "isConstant": false,
                    "isLValue": false,
                    "isPure": true,
                    "kind": "string",
                    "lValueRequested": false,
                    "nodeType": "Literal",
                    "src": "238:10:1",
                    "subdenomination": null,
                    "typeDescriptions": {
                      "typeIdentifier": "t_stringliteral_6961ba82ab463124347186eaacd7907445b2abce40f118cc3374a5f8fd621cd7",
                      "typeString": "literal_string \"BobToken\""
                    },
                    "value": "BobToken"
                  },
                  {
                    "argumentTypes": null,
                    "hexValue": "42544b",
                    "id": 43,
                    "isConstant": false,
                    "isLValue": false,
                    "isPure": true,
                    "kind": "string",
                    "lValueRequested": false,
                    "nodeType": "Literal",
                    "src": "250:5:1",
                    "subdenomination": null,
                    "typeDescriptions": {
                      "typeIdentifier": "t_stringliteral_cbd09d845af28ee5b0642a9230e90a4426144cbbda10cea126e4301fa7e1a7f0",
                      "typeString": "literal_string \"BTK\""
                    },
                    "value": "BTK"
                  },
                  {
                    "argumentTypes": null,
                    "hexValue": "36",
                    "id": 44,
                    "isConstant": false,
                    "isLValue": false,
                    "isPure": true,
                    "kind": "number",
                    "lValueRequested": false,
                    "nodeType": "Literal",
                    "src": "257:1:1",
                    "subdenomination": null,
                    "typeDescriptions": {
                      "typeIdentifier": "t_rational_6_by_1",
                      "typeString": "int_const 6"
                    },
                    "value": "6"
                  }
                ],
                "id": 45,
                "modifierName": {
                  "argumentTypes": null,
                  "id": 41,
                  "name": "ERC20Detailed",
                  "nodeType": "Identifier",
                  "overloadedDeclarations": [],
                  "referencedDeclaration": 773,
                  "src": "224:13:1",
                  "typeDescriptions": {
                    "typeIdentifier": "t_type$_t_contract$_ERC20Detailed_$773_$",
                    "typeString": "type(contract ERC20Detailed)"
                  }
                },
                "nodeType": "ModifierInvocation",
                "src": "224:35:1"
              }
            ],
            "name": "",
            "nodeType": "FunctionDefinition",
            "parameters": {
              "id": 40,
              "nodeType": "ParameterList",
              "parameters": [],
              "src": "214:2:1"
            },
            "returnParameters": {
              "id": 46,
              "nodeType": "ParameterList",
              "parameters": [],
              "src": "260:0:1"
            },
            "scope": 63,
            "src": "203:126:1",
            "stateMutability": "nonpayable",
            "superFunction": null,
            "visibility": "public"
          }
        ],
        "scope": 64,
        "src": "155:176:1"
      }
    ],
    "src": "0:332:1"
  },
  "compiler": {
    "name": "solc",
    "version": "0.5.17+commit.d19bba13.Emscripten.clang"
  },
  "networks": {
    "1605659474353": {
      "events": {
        "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925": {
          "anonymous": false,
          "inputs": [
            {
              "indexed": true,
              "internalType": "address",
              "name": "owner",
              "type": "address"
            },
            {
              "indexed": true,
              "internalType": "address",
              "name": "spender",
              "type": "address"
            },
            {
              "indexed": false,
              "internalType": "uint256",
              "name": "value",
              "type": "uint256"
            }
          ],
          "name": "Approval",
          "type": "event"
        },
        "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef": {
          "anonymous": false,
          "inputs": [
            {
              "indexed": true,
              "internalType": "address",
              "name": "from",
              "type": "address"
            },
            {
              "indexed": true,
              "internalType": "address",
              "name": "to",
              "type": "address"
            },
            {
              "indexed": false,
              "internalType": "uint256",
              "name": "value",
              "type": "uint256"
            }
          ],
          "name": "Transfer",
          "type": "event"
        }
      },
      "links": {},
      "address": "0x88E775C90031CedC975B7568D4fDCf4E3Dd06C9e",
      "transactionHash": "0x2cb1f0932a2a087a16effc772a6cdfc441f9ee85c1cee0fc9279d3b2a38378b3"
    }
  },
  "schemaVersion": "3.3.0",
  "updatedAt": "2020-11-18T00:31:27.935Z",
  "networkType": "ethereum",
  "devdoc": {
    "methods": {
      "allowance(address,address)": {
        "details": "See {IERC20-allowance}."
      },
      "approve(address,uint256)": {
        "details": "See {IERC20-approve}.     * Requirements:     * - `spender` cannot be the zero address."
      },
      "balanceOf(address)": {
        "details": "See {IERC20-balanceOf}."
      },
      "decimals()": {
        "details": "Returns the number of decimals used to get its user representation. For example, if `decimals` equals `2`, a balance of `505` tokens should be displayed to a user as `5,05` (`505 / 10 ** 2`).     * Tokens usually opt for a value of 18, imitating the relationship between Ether and Wei.     * NOTE: This information is only used for _display_ purposes: it in no way affects any of the arithmetic of the contract, including {IERC20-balanceOf} and {IERC20-transfer}."
      },
      "decreaseAllowance(address,uint256)": {
        "details": "Atomically decreases the allowance granted to `spender` by the caller.     * This is an alternative to {approve} that can be used as a mitigation for problems described in {IERC20-approve}.     * Emits an {Approval} event indicating the updated allowance.     * Requirements:     * - `spender` cannot be the zero address. - `spender` must have allowance for the caller of at least `subtractedValue`."
      },
      "increaseAllowance(address,uint256)": {
        "details": "Atomically increases the allowance granted to `spender` by the caller.     * This is an alternative to {approve} that can be used as a mitigation for problems described in {IERC20-approve}.     * Emits an {Approval} event indicating the updated allowance.     * Requirements:     * - `spender` cannot be the zero address."
      },
      "name()": {
        "details": "Returns the name of the token."
      },
      "symbol()": {
        "details": "Returns the symbol of the token, usually a shorter version of the name."
      },
      "totalSupply()": {
        "details": "See {IERC20-totalSupply}."
      },
      "transfer(address,uint256)": {
        "details": "See {IERC20-transfer}.     * Requirements:     * - `recipient` cannot be the zero address. - the caller must have a balance of at least `amount`."
      },
      "transferFrom(address,address,uint256)": {
        "details": "See {IERC20-transferFrom}.     * Emits an {Approval} event indicating the updated allowance. This is not required by the EIP. See the note at the beginning of {ERC20};     * Requirements: - `sender` and `recipient` cannot be the zero address. - `sender` must have a balance of at least `amount`. - the caller must have allowance for `sender`'s tokens of at least `amount`."
      }
    }
  },
  "userdoc": {
    "methods": {}
  }
}