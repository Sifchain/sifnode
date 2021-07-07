package types

import (
	"encoding/json"
	"fmt"
	"strconv"

	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

// OracleClaimContent is the details of how the content of the claim for each validator will be stored in the oracle
type OracleClaimContent struct {
	CosmosReceiver       sdk.AccAddress  `json:"cosmos_receiver" yaml:"cosmos_receiver"`
	Amount               sdk.Int         `json:"amount" yaml:"amount"`
	Symbol               string          `json:"symbol" yaml:"symbol"`
	TokenContractAddress EthereumAddress `json:"token_contract_address" yaml:"token_contract_address"`
	ClaimType            ClaimType       `json:"claim_type" yaml:"claim_type"`
}

// NewOracleClaimContent is a constructor function for OracleClaim
func NewOracleClaimContent(
	cosmosReceiver sdk.AccAddress, amount sdk.Int, symbol string, tokenContractAddress EthereumAddress, claimType ClaimType,
) OracleClaimContent {
	return OracleClaimContent{
		CosmosReceiver:       cosmosReceiver,
		Amount:               amount,
		Symbol:               symbol,
		TokenContractAddress: tokenContractAddress,
		ClaimType:            claimType,
	}
}

// CreateOracleClaimFromEthClaim converts a specific ethereum bridge claim to a general oracle claim to be used by
// the oracle module. The oracle module expects every claim for a particular prophecy to have the same id, so this id
// must be created in a deterministic way that all validators can follow.
// For this, we use the Nonce an Ethereum Sender provided,
// as all validators will see this same data from the smart contract.
func CreateOracleClaimFromEthClaim(ethClaim *EthBridgeClaim) (oracletypes.Claim, error) {
	oracleID := strconv.FormatInt(int64(ethClaim.NetworkDescriptor), 10) + strconv.FormatInt(ethClaim.Nonce, 10) +
		ethClaim.EthereumSender

	cosmosReceiver, err := sdk.AccAddressFromBech32(ethClaim.CosmosReceiver)
	if err != nil {
		return oracletypes.Claim{}, err
	}

	claimContent := NewOracleClaimContent(cosmosReceiver, ethClaim.Amount,
		ethClaim.Symbol, NewEthereumAddress(ethClaim.TokenContractAddress), ethClaim.ClaimType)
	claimBytes, err := json.Marshal(claimContent)
	if err != nil {
		return oracletypes.Claim{}, err
	}
	claimString := string(claimBytes)
	claim := oracletypes.NewClaim(oracleID, ethClaim.ValidatorAddress, claimString)
	return claim, nil
}

// CreateEthClaimFromOracleString converts a string
// from any generic claim from the oracle module into an ethereum bridge specific claim.
func CreateEthClaimFromOracleString(
	networkDescriptor oracletypes.NetworkDescriptor,
	bridgeContract EthereumAddress,
	nonce int64,
	ethereumAddress EthereumAddress,
	validator sdk.ValAddress,
	oracleClaimString string,
) (*EthBridgeClaim, error) {
	oracleClaim, err := CreateOracleClaimFromOracleString(oracleClaimString)
	if err != nil {
		return nil, err
	}

	return NewEthBridgeClaim(
		networkDescriptor,
		bridgeContract,
		nonce,
		oracleClaim.Symbol,
		oracleClaim.TokenContractAddress,
		ethereumAddress,
		oracleClaim.CosmosReceiver,
		validator,
		oracleClaim.Amount,
		oracleClaim.ClaimType,
	), nil
}

// CreateOracleClaimFromOracleString converts a JSON string into an OracleClaimContent struct used by this module.
// In general, it is expected that the oracle module will store claims in this JSON format
// and so this should be used to convert oracle claims.
func CreateOracleClaimFromOracleString(oracleClaimString string) (OracleClaimContent, error) {
	var oracleClaimContent OracleClaimContent

	bz := []byte(oracleClaimString)
	if err := json.Unmarshal(bz, &oracleClaimContent); err != nil {
		return OracleClaimContent{}, sdkerrors.Wrap(ErrJSONMarshalling, fmt.Sprintf("failed to parse claim: %s", err.Error()))
	}

	return oracleClaimContent, nil
}
