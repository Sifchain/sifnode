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
	_ sdk.Msg = &MsgUpdatePools{}
	_ sdk.Msg = &MsgUpdateRowanCollateral{}
	_ sdk.Msg = &MsgDewhitelist{}
	_ sdk.Msg = &MsgWhitelist{}
	_ sdk.Msg = &MsgAdminClose{}
	_ sdk.Msg = &MsgAdminCloseAll{}

	_ legacytx.LegacyMsg = &MsgOpen{}
	_ legacytx.LegacyMsg = &MsgClose{}
	_ legacytx.LegacyMsg = &MsgForceClose{}
	_ legacytx.LegacyMsg = &MsgUpdateParams{}
	_ legacytx.LegacyMsg = &MsgUpdatePools{}
	_ legacytx.LegacyMsg = &MsgUpdateRowanCollateral{}
	_ legacytx.LegacyMsg = &MsgWhitelist{}
	_ legacytx.LegacyMsg = &MsgDewhitelist{}
	_ legacytx.LegacyMsg = &MsgAdminClose{}
	_ legacytx.LegacyMsg = &MsgAdminCloseAll{}
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
	// Avoid eta (leverage - 1) going below zero.
	if m.Leverage.IsNil() || m.Leverage.LT(sdk.NewDec(1)) {
		return sdkerrors.Wrap(clptypes.ErrInValidAmount, m.Leverage.String())
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

func (m MsgUpdatePools) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgUpdatePools) Route() string {
	return RouterKey
}

func (m MsgUpdatePools) Type() string {
	return "update_pools"
}

func (m MsgUpdatePools) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}

	return nil
}

func (m MsgUpdateRowanCollateral) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}
func (m MsgUpdateRowanCollateral) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgUpdateRowanCollateral) Route() string {
	return RouterKey
}

func (m MsgUpdateRowanCollateral) Type() string {
	return "update_rowan_collateral"
}

func (m MsgUpdateRowanCollateral) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}

	return nil
}

func (m MsgUpdatePools) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (m MsgWhitelist) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgWhitelist) Route() string {
	return RouterKey
}

func (m MsgWhitelist) Type() string {
	return "whitelist"
}

func (m MsgWhitelist) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	_, err := sdk.AccAddressFromBech32(m.WhitelistedAddress)
	if err != nil {
		return err
	}

	return nil
}

func (m MsgWhitelist) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (m MsgDewhitelist) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgDewhitelist) Route() string {
	return RouterKey
}

func (m MsgDewhitelist) Type() string {
	return "dewhitelist"
}

func (m MsgDewhitelist) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	_, err := sdk.AccAddressFromBech32(m.WhitelistedAddress)
	if err != nil {
		return err
	}

	return nil
}

func (m MsgDewhitelist) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (m MsgAdminClose) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgAdminClose) Route() string {
	return RouterKey
}

func (m MsgAdminClose) Type() string {
	return "admin_close"
}

func (m MsgAdminClose) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}

	return nil
}

func (m MsgAdminClose) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (m MsgAdminCloseAll) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgAdminCloseAll) Route() string {
	return RouterKey
}

func (m MsgAdminCloseAll) Type() string {
	return "admin_close_all"
}

func (m MsgAdminCloseAll) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}

	return nil
}

func (m MsgAdminCloseAll) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}
