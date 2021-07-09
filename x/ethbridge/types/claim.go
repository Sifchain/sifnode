package types

import (
	"fmt"
	"strconv"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/ethereum/go-ethereum/crypto"
)

// NewEthBridgeClaim is a constructor function for NewEthBridgeClaim
func NewEthBridgeClaim(
	networkDescriptor oracletypes.NetworkDescriptor,
	bridgeContract EthereumAddress,
	nonce int64,
	symbol string,
	tokenContact EthereumAddress,
	ethereumSender EthereumAddress,
	cosmosReceiver sdk.AccAddress,
	validator sdk.ValAddress,
	amount sdk.Int,
	claimType ClaimType,
) *EthBridgeClaim {
	return &EthBridgeClaim{
		NetworkDescriptor:     networkDescriptor,
		BridgeContractAddress: bridgeContract.String(),
		Nonce:                 nonce,
		Symbol:                symbol,
		TokenContractAddress:  tokenContact.String(),
		EthereumSender:        ethereumSender.String(),
		CosmosReceiver:        cosmosReceiver.String(),
		ValidatorAddress:      validator.String(),
		Amount:                amount,
		ClaimType:             claimType,
	}
}

// GetProphecyID compute oracle id, get from keccak256 of the all content in claim
func (claim *EthBridgeClaim) GetProphecyID() []byte {
	allContentString := fmt.Sprintf("%s%s%s%s%s%s%s%s%s",
		claim.NetworkDescriptor.String(),
		claim.BridgeContractAddress,
		strconv.Itoa(int(claim.Nonce)),
		claim.Symbol,
		claim.TokenContractAddress,
		claim.EthereumSender,
		claim.CosmosReceiver,
		claim.Amount.String(),
		claim.ClaimType.String(),
	)
	claimBytes := []byte(allContentString)
	hashBytes := crypto.Keccak256(claimBytes)
	return hashBytes
}
