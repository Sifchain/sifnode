package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
	"go.uber.org/zap"

	"github.com/Sifchain/sifnode/x/oracle"
)

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetAccount(sdk.Context, sdk.AccAddress) authexported.Account
}

// SupplyKeeper defines the expected supply keeper
type SupplyKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	SetModuleAccount(sdk.Context, supplyexported.ModuleAccountI)
}

// OracleKeeper defines the expected oracle keeper
type OracleKeeper interface {
	ProcessClaim(ctx sdk.Context, claim oracle.Claim, sugaredLogger *zap.SugaredLogger) (oracle.Status, error)
	GetProphecy(ctx sdk.Context, id string) (oracle.Prophecy, bool)
	ProcessUpdateWhiteListValidator(ctx sdk.Context, cosmosSender sdk.AccAddress, validator sdk.ValAddress, operationtype string, sugaredLogger *zap.SugaredLogger) error
}
