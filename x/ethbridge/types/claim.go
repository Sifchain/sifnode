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
	nonce int64,
	symbol string,
	tokenContract EthereumAddress,
	ethereumSender EthereumAddress,
	cosmosReceiver sdk.AccAddress,
	validator sdk.ValAddress,
	amount sdk.Int,
	claimType ClaimType,
	tokenName string,
	decimals int32,
) *EthBridgeClaim {
	denomHash := GetDenomHash(networkDescriptor, tokenContract.String(), decimals, tokenName, symbol)
	return &EthBridgeClaim{
		NetworkDescriptor:     networkDescriptor,
		BridgeContractAddress: bridgeContract.String(),
		Nonce:                 nonce,
		Symbol:                symbol,
		TokenContractAddress:  tokenContract.String(),
		EthereumSender:        ethereumSender.String(),
		CosmosReceiver:        cosmosReceiver.String(),
		ValidatorAddress:      validator.String(),
		Amount:                amount,
		ClaimType:             claimType,
		TokenName:             tokenName,
		Decimals:              decimals,
		DenomHash:             denomHash,
	}
}

// GetProphecyID compute oracle id, get from keccak256 of the all content in claim
func (claim *EthBridgeClaim) GetProphecyID() []byte {
	allContentString := fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s%s",
		claim.NetworkDescriptor.String(),
		claim.BridgeContractAddress,
		strconv.Itoa(int(claim.Nonce)),
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

func GetDenomHash(
	networkDescriptor oracletypes.NetworkDescriptor,
	tokenContractAddress string,
	decimals int32,
	tokenName string,
	symbol string,
) string {
	/**
	  * Metadata Denom Naming Convention:
	  * For all pegged ERC20 assets, their respective token names on sifchain will be composed of the
	  * following four elements: prefix to define the object type (coin, nft, etc), network descriptor,
	  * ERC20 token address, and the decimals of that token. Fields will not be separated by any delimiter
	  * character. A pegged ERC20 asset with token address 0xbF45BFc92ebD305d4C0baf8395c4299bdFCE9EA2, a
	  * network descriptor of 2, 9 decimals, a name of “wBTC” and symbol “WBTC” will add all of the strings
	  * together to get this output:
	  *    0xbF45BFc92ebD305d4C0baf8395c4299bdFCE9EA229wBTCWBTC
	  *
	  * Then, that data will be hashed with SHA256 to produce the following hash:
	  *    179e6a6f8ab6efb5fa1f3992aef69f855628cfd27868a1be0525f40b456494ff
	  *
	**/
	// No Prefix Yet....
	// "{Network Descriptor}{ERC20 Token Address}{Decimals}{Token Name}{Token Symbol}"
	denomHashedString := fmt.Sprintf("%d%s%d%s%s",
		networkDescriptor,
		strings.ToLower(tokenContractAddress),
		decimals,
		strings.ToLower(tokenName),
		strings.ToLower(symbol),
	)

	rawDenomHash := sha256.Sum256([]byte(denomHashedString))
	denomHash := hex.EncodeToString(rawDenomHash[:])

	return denomHash
}
