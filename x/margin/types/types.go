package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	ModuleName = "margin"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// QuerierRoute is the querier route
	QuerierRoute = ModuleName

	// RouterKey is the msg router key
	RouterKey = ModuleName
)

func NewMTP(signer string, collateralAsset string, collateralAmount sdk.Uint, borrowAsset string) MTP {
	return MTP{
		Address:          signer,
		CollateralAsset:  collateralAsset,
		CollateralAmount: collateralAmount,
		LiabilitiesP:     sdk.Uint{},
		LiabilitiesI:     sdk.Uint{},
		CustodyAsset:     borrowAsset,
		CustodyAmount:    sdk.Uint{},
		Leverage:         sdk.Uint{},
		MtpHealth:        sdk.Dec{},
	}
}

func (mtp MTP) Validate() error {
	if mtp.CollateralAsset == "" {
		return sdkerrors.Wrap(ErrMTPInvalid, "no asset specified")
	}
	if mtp.Address == "" {
		return sdkerrors.Wrap(ErrMTPInvalid, "no address specified")
	}

	return nil
}

func GetSettlementAsset() string {
	return "rowan"
}
