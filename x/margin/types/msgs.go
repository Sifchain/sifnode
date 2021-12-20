package types

import (
	"strings"

	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgOpenLong{}
	_ sdk.Msg = &MsgCloseLong{}
)

func Validate(asset string) bool {
	if !clptypes.VerifyRange(len(strings.TrimSpace(asset)), 0, clptypes.MaxSymbolLength) {
		return false
	}
	coin := sdk.NewCoin(asset, sdk.OneInt())
	return coin.IsValid()
}

func (m MsgOpenLong) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}

	if !Validate(m.CollateralAsset) {
		return sdkerrors.Wrap(clptypes.ErrInValidAsset, m.CollateralAsset)
	}
	if !Validate(m.BorrowAsset) {
		return sdkerrors.Wrap(clptypes.ErrInValidAsset, m.BorrowAsset)
	}

	if m.CollateralAmount.IsZero() {
		return sdkerrors.Wrap(clptypes.ErrInValidAmount, m.CollateralAmount.String())
	}
	return nil
}

func (m MsgOpenLong) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (m MsgCloseLong) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	if !Validate(m.CollateralAsset) {
		return sdkerrors.Wrap(clptypes.ErrInValidAsset, m.CollateralAsset)
	}
	if !Validate(m.BorrowAsset) {
		return sdkerrors.Wrap(clptypes.ErrInValidAsset, m.BorrowAsset)
	}

	return nil
}

func (m MsgCloseLong) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}
