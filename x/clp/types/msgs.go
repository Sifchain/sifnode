package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strconv"
)

var (
	_ sdk.Msg = &MsgRemoveLiquidity{}
	_ sdk.Msg = &MsgCreatePool{}
	_ sdk.Msg = &MsgAddLiquidity{}
	_ sdk.Msg = &MsgSwap{}
)

type MsgSwap struct {
	Signer        sdk.AccAddress
	SentAsset     Asset
	ReceivedAsset Asset
	SentAmount    uint
}

func NewMsgSwap(signer sdk.AccAddress, sentAsset Asset, receivedAsset Asset, sentAmount uint) *MsgSwap {
	return &MsgSwap{Signer: signer, SentAsset: sentAsset, ReceivedAsset: receivedAsset, SentAmount: sentAmount}
}

func (m MsgSwap) Route() string {
	return RouterKey
}

func (m MsgSwap) Type() string {
	return "swap"
}

func (m MsgSwap) ValidateBasic() error {
	if m.Signer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer.String())
	}
	if !m.SentAsset.Validate() {
		return sdkerrors.Wrap(InValidAsset, m.SentAsset.Symbol)
	}
	if !m.ReceivedAsset.Validate() {
		return sdkerrors.Wrap(InValidAsset, m.SentAsset.Symbol)
	}
	return nil
}

func (m MsgSwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgSwap) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Signer}
}

type MsgRemoveLiquidity struct {
	Signer        sdk.AccAddress
	ExternalAsset Asset
	WBasisPoints  uint
	Asymmetry     uint
}

func NewMsgRemoveLiquidity(signer sdk.AccAddress, externalAsset Asset, wBasisPoints uint, asymmetry uint) MsgRemoveLiquidity {
	return MsgRemoveLiquidity{Signer: signer, ExternalAsset: externalAsset, WBasisPoints: wBasisPoints, Asymmetry: asymmetry}
}

func (m MsgRemoveLiquidity) Route() string {
	return RouterKey
}

func (m MsgRemoveLiquidity) Type() string {
	return "remove_liquidity"
}

func (m MsgRemoveLiquidity) ValidateBasic() error {
	if m.Signer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer.String())
	}
	if !m.ExternalAsset.Validate() {
		return sdkerrors.Wrap(InValidAsset, m.ExternalAsset.Symbol)
	}
	if !(m.WBasisPoints > 0) {
		return sdkerrors.Wrap(InValidAmount, strconv.Itoa(int(m.WBasisPoints)))
	}
	if !(m.Asymmetry > 0) {
		return sdkerrors.Wrap(InValidAmount, strconv.Itoa(int(m.Asymmetry)))
	}
	return nil
}

func (m MsgRemoveLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgRemoveLiquidity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Signer}
}

type MsgAddLiquidity struct {
	Signer              sdk.AccAddress
	ExternalAsset       Asset
	NativeAssetAmount   uint
	ExternalAssetAmount uint
}

func NewMsgAddLiquidity(signer sdk.AccAddress, externalAsset Asset, nativeAssetAmount uint, externalAssetAmount uint) MsgAddLiquidity {
	return MsgAddLiquidity{Signer: signer, ExternalAsset: externalAsset, NativeAssetAmount: nativeAssetAmount, ExternalAssetAmount: externalAssetAmount}
}

func (m MsgAddLiquidity) Route() string {
	return RouterKey
}

func (m MsgAddLiquidity) Type() string {
	return "add_liquidity"
}

func (m MsgAddLiquidity) ValidateBasic() error {
	if m.Signer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer.String())
	}
	if !m.ExternalAsset.Validate() {
		return sdkerrors.Wrap(InValidAsset, m.ExternalAsset.Symbol)
	}
	if !(m.NativeAssetAmount > 0) {
		return sdkerrors.Wrap(InValidAmount, strconv.Itoa(int(m.NativeAssetAmount)))
	}
	if !(m.ExternalAssetAmount > 0) {
		return sdkerrors.Wrap(InValidAmount, strconv.Itoa(int(m.NativeAssetAmount)))
	}
	return nil
}

func (m MsgAddLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgAddLiquidity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Signer}
}

type MsgCreatePool struct {
	Signer              sdk.AccAddress
	ExternalAsset       Asset
	NativeAssetAmount   uint
	ExternalAssetAmount uint
}

func NewMsgCreatePool(signer sdk.AccAddress, externalAsset Asset, nativeAssetAmount uint, externalAssetAmount uint) MsgCreatePool {
	return MsgCreatePool{Signer: signer, ExternalAsset: externalAsset, NativeAssetAmount: nativeAssetAmount, ExternalAssetAmount: externalAssetAmount}
}

func (m MsgCreatePool) Route() string {
	return RouterKey
}

func (m MsgCreatePool) Type() string {
	return "create_pool"
}

func (m MsgCreatePool) ValidateBasic() error {
	if m.Signer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer.String())
	}
	if !m.ExternalAsset.Validate() {
		return sdkerrors.Wrap(InValidAsset, m.ExternalAsset.Symbol)
	}
	if !(m.NativeAssetAmount > 0) {
		return sdkerrors.Wrap(InValidAmount, strconv.Itoa(int(m.NativeAssetAmount)))
	}
	if !(m.ExternalAssetAmount > 0) {
		return sdkerrors.Wrap(InValidAmount, strconv.Itoa(int(m.NativeAssetAmount)))
	}
	return nil
}

func (m MsgCreatePool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgCreatePool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Signer}
}
