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
	signer        sdk.AccAddress
	sentAsset     Asset
	receivedAsset Asset
	sentAmount    uint
}

func NewMsgSwap(signer sdk.AccAddress, sentAsset Asset, receivedAsset Asset, sentAmount uint) *MsgSwap {
	return &MsgSwap{signer: signer, sentAsset: sentAsset, receivedAsset: receivedAsset, sentAmount: sentAmount}
}

func (m MsgSwap) Route() string {
	return RouterKey
}

func (m MsgSwap) Type() string {
	return "swap"
}

func (m MsgSwap) ValidateBasic() error {
	if m.signer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.signer.String())
	}
	if !m.sentAsset.Validate() {
		return sdkerrors.Wrap(InValidAsset, m.sentAsset.Symbol)
	}
	if !m.receivedAsset.Validate() {
		return sdkerrors.Wrap(InValidAsset, m.sentAsset.Symbol)
	}
	if m.sentAmount < 0 {
		return sdkerrors.Wrap(InValidAmount, strconv.Itoa(int(m.sentAmount)))
	}
	return nil
}

func (m MsgSwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgSwap) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.signer}
}

type MsgRemoveLiquidity struct {
	signer        sdk.AccAddress
	externalAsset Asset
	wBasisPoints  uint
	asymmetry     uint
}

func NewMsgRemoveLiquidity(signer sdk.AccAddress, externalAsset Asset, wBasisPoints uint, asymmetry uint) MsgRemoveLiquidity {
	return MsgRemoveLiquidity{signer: signer, externalAsset: externalAsset, wBasisPoints: wBasisPoints, asymmetry: asymmetry}
}

func (m MsgRemoveLiquidity) Route() string {
	return RouterKey
}

func (m MsgRemoveLiquidity) Type() string {
	return "remove_liquidity"
}

func (m MsgRemoveLiquidity) ValidateBasic() error {
	if m.signer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.signer.String())
	}
	if !m.externalAsset.Validate() {
		return sdkerrors.Wrap(InValidAsset, m.externalAsset.Symbol)
	}
	if m.wBasisPoints < 0 {
		return sdkerrors.Wrap(InValidAmount, strconv.Itoa(int(m.wBasisPoints)))
	}
	if m.asymmetry < 0 {
		return sdkerrors.Wrap(InValidAmount, strconv.Itoa(int(m.asymmetry)))
	}
	return nil
}

func (m MsgRemoveLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgRemoveLiquidity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.signer}
}

type MsgAddLiquidity struct {
	signer              sdk.AccAddress
	externalAsset       Asset
	nativeAssetAmount   uint
	externalAssetAmount uint
}

func NewMsgAddLiquidity(signer sdk.AccAddress, externalAsset Asset, nativeAssetAmount uint, externalAssetAmount uint) MsgAddLiquidity {
	return MsgAddLiquidity{signer: signer, externalAsset: externalAsset, nativeAssetAmount: nativeAssetAmount, externalAssetAmount: externalAssetAmount}
}

func (m MsgAddLiquidity) Route() string {
	return RouterKey
}

func (m MsgAddLiquidity) Type() string {
	return "add_liquidity"
}

func (m MsgAddLiquidity) ValidateBasic() error {
	if m.signer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.signer.String())
	}
	if !m.externalAsset.Validate() {
		return sdkerrors.Wrap(InValidAsset, m.externalAsset.Symbol)
	}
	if m.nativeAssetAmount < 0 {
		return sdkerrors.Wrap(InValidAmount, strconv.Itoa(int(m.nativeAssetAmount)))
	}
	if m.externalAssetAmount < 0 {
		return sdkerrors.Wrap(InValidAmount, strconv.Itoa(int(m.nativeAssetAmount)))
	}
	return nil
}

func (m MsgAddLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgAddLiquidity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.signer}
}

type MsgCreatePool struct {
	signer              sdk.AccAddress
	externalAsset       Asset
	nativeAssetAmount   uint
	externalAssetAmount uint
}

func NewMsgCreatePool(signer sdk.AccAddress, externalAsset Asset, nativeAssetAmount uint, externalAssetAmount uint) MsgCreatePool {
	return MsgCreatePool{signer: signer, externalAsset: externalAsset, nativeAssetAmount: nativeAssetAmount, externalAssetAmount: externalAssetAmount}
}

func (m MsgCreatePool) Route() string {
	return RouterKey
}

func (m MsgCreatePool) Type() string {
	return "create_pool"
}

func (m MsgCreatePool) ValidateBasic() error {
	if m.signer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.signer.String())
	}
	if !m.externalAsset.Validate() {
		return sdkerrors.Wrap(InValidAsset, m.externalAsset.Symbol)
	}
	if m.nativeAssetAmount < 0 {
		return sdkerrors.Wrap(InValidAmount, strconv.Itoa(int(m.nativeAssetAmount)))
	}
	if m.externalAssetAmount < 0 {
		return sdkerrors.Wrap(InValidAmount, strconv.Itoa(int(m.nativeAssetAmount)))
	}
	return nil
}

func (m MsgCreatePool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgCreatePool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.signer}
}
