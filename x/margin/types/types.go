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

func NewMTP(signer string, collateralAsset string, borrowAsset string, position Position, leverage sdk.Dec) *MTP {
	return &MTP{
		Address:                  signer,
		CollateralAsset:          collateralAsset,
		CollateralAmount:         sdk.ZeroUint(),
		Liabilities:              sdk.ZeroUint(),
		InterestPaidCollateral:   sdk.ZeroUint(),
		InterestPaidCustody:      sdk.ZeroUint(),
		InterestUnpaidCollateral: sdk.ZeroUint(),
		CustodyAsset:             borrowAsset,
		CustodyAmount:            sdk.ZeroUint(),
		Leverage:                 leverage,
		MtpHealth:                sdk.ZeroDec(),
		Position:                 position,
	}
}

func (mtp MTP) Validate() error {
	if mtp.CollateralAsset == "" {
		return sdkerrors.Wrap(ErrMTPInvalid, "no asset specified")
	}
	if mtp.Address == "" {
		return sdkerrors.Wrap(ErrMTPInvalid, "no address specified")
	}
	if mtp.Position == Position_UNSPECIFIED {
		return sdkerrors.Wrap(ErrMTPInvalid, "no position specified")
	}
	if mtp.Id == 0 {
		return sdkerrors.Wrap(ErrMTPInvalid, "no id specified")
	}

	return nil
}

func GetSettlementAsset() string {
	return "rowan"
}

func GetPositionFromString(s string) Position {
	switch s {
	case "long":
		return Position_LONG
	case "short":
		return Position_SHORT
	default:
		return Position_UNSPECIFIED
	}
}

func StringCompare(a, b string) bool {
	return a == b
}
