package types

import (
	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BankKeeper interface {
	//SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	//SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

type CLPKeeper interface {
	GetPool(ctx sdk.Context, symbol string) (clptypes.Pool, error)
	SetPool(ctx sdk.Context, pool *clptypes.Pool) error
	GetNormalizationFactorForAsset(sdk.Context, string) (sdk.Dec, bool, error)

	ValidateZero(inputs []sdk.Uint) bool
	ReducePrecision(dec sdk.Dec, po int64) sdk.Dec
	IncreasePrecision(dec sdk.Dec, po int64) sdk.Dec
	GetMinLen(inputs []sdk.Uint) int64
}
