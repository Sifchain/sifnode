package types

import (
	admintypes "github.com/Sifchain/sifnode/x/admin/types"
	tokenregistryTypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// ParamSubspace defines the expected Subspace interfacace
type ParamSubspace interface {
	WithKeyTable(table paramtypes.KeyTable) paramtypes.Subspace
	Get(ctx sdk.Context, key []byte, ptr interface{})
	GetParamSet(ctx sdk.Context, ps paramtypes.ParamSet)
	SetParamSet(ctx sdk.Context, ps paramtypes.ParamSet)
}

//go:generate mockery --srcpkg . --name BankKeeper --structname BankKeeper --filename bank_keeper.go --with-expecter
type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool
}

type AuthKeeper interface {
	SetModuleAccount(sdk.Context, authtypes.ModuleAccountI)
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
}

type TokenRegistryKeeper interface {
	GetEntry(registry tokenregistryTypes.Registry, denom string) (*tokenregistryTypes.RegistryEntry, error)
	CheckEntryPermissions(entry *tokenregistryTypes.RegistryEntry, permissions []tokenregistryTypes.Permission) bool
	GetRegistry(ctx sdk.Context) tokenregistryTypes.Registry
}

type AdminKeeper interface {
	IsAdminAccount(ctx sdk.Context, moduleName admintypes.AdminType, adminAccount sdk.AccAddress) bool
}
