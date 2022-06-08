package types

import (
	"strings"

	clptypes "github.com/Sifchain/sifnode/x/clp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"
)

var (
	_ sdk.Msg = &MsgOpen{}
	_ sdk.Msg = &MsgClose{}
	_ sdk.Msg = &MsgForceClose{}
	_ sdk.Msg = &MsgUpdateParams{}

	_ legacytx.LegacyMsg = &MsgOpen{}
	_ legacytx.LegacyMsg = &MsgClose{}
	_ legacytx.LegacyMsg = &MsgForceClose{}
	_ legacytx.LegacyMsg = &MsgUpdateParams{}
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

func (m MsgOpen) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgOpen) Route() string {
	return RouterKey
}

func (m MsgOpen) Type() string {
	return "open"
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

func (m MsgClose) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgClose) Route() string {
	return RouterKey
}

func (m MsgClose) Type() string {
	return "close"
}

func (m MsgClose) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	if m.Id == 0 {
		return sdkerrors.Wrap(ErrMTPDoesNotExist, "no id specified")
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

func (m MsgForceClose) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgForceClose) Route() string {
	return RouterKey
}

func (m MsgForceClose) Type() string {
	return "force_close"
}

func (m MsgForceClose) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	if len(m.MtpAddress) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.MtpAddress)
	}
	if m.Id == 0 {
		return sdkerrors.Wrap(ErrMTPDoesNotExist, "no id specified")
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

func (m MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgUpdateParams) Route() string {
	return RouterKey
}

func (m MsgUpdateParams) Type() string {
	return "update_params"
}

func (m MsgUpdateParams) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}

	return nil
}

func (m MsgUpdateParams) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}
