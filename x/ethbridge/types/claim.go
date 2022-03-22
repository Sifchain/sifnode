package types

import (
	"fmt"
	"strconv"
	"strings"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/ethereum/go-ethereum/crypto"
)

// NewEthBridgeClaim is a constructor function for NewEthBridgeClaim
func NewEthBridgeClaim(
	networkDescriptor oracletypes.NetworkDescriptor,
	bridgeContract EthereumAddress,
	ethereumLockBurnSequence uint64,
	symbol string,
	tokenContract EthereumAddress,
	ethereumSender EthereumAddress,
	cosmosReceiver sdk.AccAddress,
	validator sdk.ValAddress,
	amount sdk.Int,
	claimType ClaimType,
	tokenName string,
	decimals uint8,
) *EthBridgeClaim {
	denom := GetDenom(networkDescriptor, tokenContract)
	return &EthBridgeClaim{
		NetworkDescriptor:        networkDescriptor,
		BridgeContractAddress:    bridgeContract.String(),
		EthereumLockBurnSequence: ethereumLockBurnSequence,
		Symbol:                   symbol,
		TokenContractAddress:     tokenContract.String(),
		EthereumSender:           ethereumSender.String(),
		CosmosReceiver:           cosmosReceiver.String(),
		ValidatorAddress:         validator.String(),
		Amount:                   amount,
		ClaimType:                claimType,
		TokenName:                tokenName,
		Decimals:                 int64(decimals),
		Denom:                    denom,
	}
}

// GetProphecyID compute oracle id, get from keccak256 of the all content in claim
func (claim *EthBridgeClaim) GetProphecyID() []byte {
	allContentString := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s%s",
		claim.NetworkDescriptor.String(),
		claim.BridgeContractAddress,
		strconv.Itoa(int(claim.EthereumLockBurnSequence)),
		claim.Symbol,
		claim.TokenContractAddress,
		claim.EthereumSender,
		claim.CosmosReceiver,
		claim.Amount.String(),
		claim.ClaimType.String(),
		claim.TokenName,
		strconv.Itoa(int(claim.Decimals)),
		claim.Denom,
	)
	claimBytes := []byte(allContentString)
	hashBytes := crypto.Keccak256(claimBytes)
	return hashBytes
}

/**
  Metadata Denom Naming Convention:
  For all pegged ERC20 assets, their respective token names on sifchain will be
  composed of the following three elements: the prefix sifBridge, the network descriptor
  (fixed as a four digit decimal number), and the ERC20 token address. Fields will not
  be separated by any delimiter character. All hexadecimal characters will be made
  lower case before hashing. A pegged ERC20 asset with token address
  0xbF45BFc92ebD305d4C0baf8395c4299bdFCE9EA2, a network descriptor of 2 has this denom:

            sifBridge0020xbf45bfc92ebd305d4c0baf8395c4299bdfce9ea2

  **WARNING**: This function will PANIC if a networkDescriptor is below 0 or larger then
  9999. This is because the networkDescriptor field is fixed to only be four digits and
  this function is used internally to calculate token Denoms.

**/
func GetDenom(
	networkDescriptor oracletypes.NetworkDescriptor,
	tokenContractAddress EthereumAddress,
) string {
	if (networkDescriptor < 0) || (networkDescriptor > 9999) {
		panic("Error: Network Descriptor must be between 0-9999")
	}
	// sifBridge + 0000 (four digit base ten number) + 0x (hex address)
	denomString := fmt.Sprintf("sifBridge%04d%s",
		networkDescriptor,
		strings.ToLower(tokenContractAddress.String()),
	)

	return denomString
}
