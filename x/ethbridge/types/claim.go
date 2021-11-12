package types

import (
	"crypto/sha256"
	"encoding/hex"
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
	denomHash := GetDenomHash(networkDescriptor, tokenContract)
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
		DenomHash:                denomHash,
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
		claim.DenomHash,
	)
	claimBytes := []byte(allContentString)
	hashBytes := crypto.Keccak256(claimBytes)
	return hashBytes
}

/**
  Metadata Denom Naming Convention:
  For all pegged ERC20 assets, their respective token names on sifchain will be
  composed of the following two elements: network descriptor, ERC20 token address
  Fields will not be separated by any delimiter character.
  All characters will be made lower case before hashing.
  A pegged ERC20 asset with token address 0xbF45BFc92ebD305d4C0baf8395c4299bdFCE9EA2,
  a network descriptor of 2 has this output:

            20xbf45bfc92ebd305d4c0baf8395c4299bdfce9ea2

  Then, that data will be hashed with SHA256 and prefixed with the
  string ‘sif’ to produce the following hash:

           sife0d5240024941c95aa2ca714f4d798f81f36da2cb8ed0c2318970c12b4acca1f

**/
func GetDenomHash(
	networkDescriptor oracletypes.NetworkDescriptor,
	tokenContractAddress EthereumAddress,
) string {
	denomHashedString := fmt.Sprintf("%d%s",
		networkDescriptor,
		strings.ToLower(tokenContractAddress.String()),
	)

	rawDenomHash := sha256.Sum256([]byte(denomHashedString))
	// Cosmos SDK requires first character to be [a-zA-Z] so we prepend sif
	denomHash := "sif" + hex.EncodeToString(rawDenomHash[:])

	return denomHash
}
