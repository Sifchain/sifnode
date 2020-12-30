package types

import (
	"io"
	"log"

	sdkContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/sethvargo/go-password/password"
	"github.com/tendermint/go-amino"
	tmLog "github.com/tendermint/tendermint/libs/log"
)

// CosmosContext wrapper all variable needed to interact with sifchain
type CosmosContext struct {
	Cdc              *codec.Codec
	ValidatorName    string
	ValidatorAddress sdk.ValAddress
	CliCtx           sdkContext.CLIContext
	TxBldr           authtypes.TxBuilder
	TempPassword     string
	Logger           tmLog.Logger
}

// NewKeybase create a new keybase instance
func NewKeybase(validatorMoniker, mnemonic, password string) (keys.Keybase, keys.Info, error) {
	keybase := keys.NewInMemory()
	hdpath := *hd.NewFundraiserParams(0, sdk.CoinType, 0)
	info, err := keybase.CreateAccount(validatorMoniker, mnemonic, "", password, hdpath.String(), keys.Secp256k1)
	if err != nil {
		return nil, nil, err
	}

	return keybase, info, nil
}

// LoadTendermintCLIContext : loads CLI context for tendermint txs
func LoadTendermintCLIContext(appCodec *amino.Codec, validatorAddress sdk.ValAddress, validatorName string,
	rpcURL string, chainID string) (sdkContext.CLIContext, error) {
	// Create the new CLI context
	cliCtx := sdkContext.NewCLIContext().
		WithCodec(appCodec).
		WithFromAddress(sdk.AccAddress(validatorAddress)).
		WithFromName(validatorName)

	if rpcURL != "" {
		cliCtx = cliCtx.WithNodeURI(rpcURL)
	}
	cliCtx.SkipConfirm = true

	// Confirm that the validator's address exists
	accountRetriever := authtypes.NewAccountRetriever(cliCtx)
	err := accountRetriever.EnsureExists(sdk.AccAddress(validatorAddress))
	if err != nil {
		log.Println(err)
		return sdkContext.CLIContext{}, err
	}
	return cliCtx, nil
}

// NewCosmosContext get a new instance of CosmosContext
func NewCosmosContext(inBuf io.Reader, cdc *codec.Codec, rpcURL, validatorMoniker, mnemonic, chainID string,
	logger tmLog.Logger) (*CosmosContext, error) {
	tempPassword, _ := password.Generate(32, 5, 0, false, false)
	keybase, info, err := NewKeybase(validatorMoniker, mnemonic, tempPassword)
	if err != nil {
		return nil, err
	}

	validatorAddress := sdk.ValAddress(info.GetAddress())

	// Load CLI context and Tx builder
	cliCtx, err := LoadTendermintCLIContext(cdc, validatorAddress, validatorMoniker, rpcURL, chainID)
	if err != nil {
		return nil, err
	}

	txBldr := authtypes.NewTxBuilderFromCLI(inBuf).
		WithTxEncoder(utils.GetTxEncoder(cdc)).
		WithChainID(chainID).
		WithKeybase(keybase)

	return &CosmosContext{
		Cdc:              cdc,
		ValidatorName:    validatorMoniker,
		ValidatorAddress: validatorAddress,
		CliCtx:           cliCtx,
		TxBldr:           txBldr,
		// PrivateKey       *ecdsa.PrivateKey
		TempPassword: tempPassword,
		Logger:       logger,
	}, nil
}
