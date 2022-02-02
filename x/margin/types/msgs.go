package types

import (
	"strings"

	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgOpen{}
	_ sdk.Msg = &MsgClose{}
	_ sdk.Msg = &MsgForceClose{}
)

func Validate(asset string) bool {
	if !clptypes.VerifyRange(len(strings.TrimSpace(asset)), 0, clptypes.MaxSymbolLength) {
		return false
	}
	coin := sdk.NewCoin(asset, sdk.OneInt())
	return coin.IsValid()
}

func IsValidPosition(position Position) bool {
	switch position {
	case Position_LONG:
		return true
	case Position_SHORT:
		return true
	default:
		return false
	}
}

func (m MsgOpen) ValidateBasic() error {
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

	ok := IsValidPosition(m.Position)
	if !ok {
		return sdkerrors.Wrap(ErrInvalidPosition, m.Position.String())
	}

	return nil
}

func (m MsgOpen) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (m MsgClose) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	if !Validate(m.CollateralAsset) {
		return sdkerrors.Wrap(clptypes.ErrInValidAsset, m.CollateralAsset)
	}
	if !Validate(m.BorrowAsset) {
		return sdkerrors.Wrap(clptypes.ErrInValidAsset, m.BorrowAsset)
	}
	ok := IsValidPosition(m.Position)
	if !ok {
		return sdkerrors.Wrap(ErrInvalidPosition, m.Position.String())
	}

	return nil
}

func (m MsgClose) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (m MsgForceClose) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	if len(m.MtpAddress) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.MtpAddress)
	}
	if !Validate(m.CollateralAsset) {
		return sdkerrors.Wrap(clptypes.ErrInValidAsset, m.CollateralAsset)
	}
	if !Validate(m.BorrowAsset) {
		return sdkerrors.Wrap(clptypes.ErrInValidAsset, m.BorrowAsset)
	}
	ok := IsValidPosition(m.Position)
	if !ok {
		return sdkerrors.Wrap(ErrInvalidPosition, m.Position.String())
	}

	return nil
}

func (m MsgForceClose) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}
