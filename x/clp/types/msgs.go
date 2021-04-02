package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgRemoveLiquidity{}
	_ sdk.Msg = &MsgCreatePool{}
	_ sdk.Msg = &MsgAddLiquidity{}
	_ sdk.Msg = &MsgSwap{}
	_ sdk.Msg = &MsgDecommissionPool{}
)

// type MsgDecommissionPool struct {
// 	Signer sdk.AccAddress `json:"signer"`
// 	Symbol string         `json:"symbol"`
// }

func NewMsgDecommissionPool(signer sdk.AccAddress, symbol string) MsgDecommissionPool {
	return MsgDecommissionPool{Signer: signer.String(), Symbol: symbol}
}

func (m MsgDecommissionPool) Route() string {
	return RouterKey
}

func (m MsgDecommissionPool) Type() string {
	return "decommission_pool"
}

func (m MsgDecommissionPool) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	if !VerifyRange(len(strings.TrimSpace(m.Symbol)), 0, MaxSymbolLength) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, m.Symbol)
	}
	return nil
}

func (m MsgDecommissionPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgDecommissionPool) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgSwap(signer sdk.AccAddress, sentAsset Asset, receivedAsset Asset, sentAmount sdk.Uint, minReceivingAmount sdk.Uint) MsgSwap {
	return MsgSwap{Signer: signer.String(), SentAsset: &sentAsset, ReceivedAsset: &receivedAsset, SentAmount: sentAmount, MinReceivingAmount: minReceivingAmount}
}

func (m MsgSwap) Route() string {
	return RouterKey
}

func (m MsgSwap) Type() string {
	return "swap"
}

func (m MsgSwap) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	if !m.SentAsset.Validate() {
		return sdkerrors.Wrap(ErrInValidAsset, m.SentAsset.Symbol)
	}
	if !m.ReceivedAsset.Validate() {
		return sdkerrors.Wrap(ErrInValidAsset, m.ReceivedAsset.Symbol)
	}

	if m.SentAsset == m.ReceivedAsset {
		return sdkerrors.Wrap(ErrInValidAsset, "Sent And Received asset cannot be the same")
	}
	if m.SentAmount.IsZero() {
		return sdkerrors.Wrap(ErrInValidAmount, m.SentAmount.String())
	}
	return nil
}

func (m MsgSwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgSwap) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgRemoveLiquidity(signer sdk.AccAddress, externalAsset Asset, wBasisPoints sdk.Int, asymmetry sdk.Int) MsgRemoveLiquidity {
	return MsgRemoveLiquidity{Signer: signer.String(), ExternalAsset: &externalAsset, WBasisPoints: wBasisPoints, Asymmetry: asymmetry}
}

func (m MsgRemoveLiquidity) Route() string {
	return RouterKey
}

func (m MsgRemoveLiquidity) Type() string {
	return "remove_liquidity"
}

func (m MsgRemoveLiquidity) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	if !m.ExternalAsset.Validate() {
		return sdkerrors.Wrap(ErrInValidAsset, m.ExternalAsset.Symbol)
	}
	if !(m.WBasisPoints.IsPositive()) || m.WBasisPoints.GT(sdk.NewInt(MaxWbasis)) {
		return sdkerrors.Wrap(ErrInvalidWBasis, m.WBasisPoints.String())
	}
	if m.Asymmetry.GT(sdk.NewInt(10000)) || m.Asymmetry.LT(sdk.NewInt(-10000)) {
		return sdkerrors.Wrap(ErrInvalidAsymmetry, m.Asymmetry.String())
	}
	return nil
}

func (m MsgRemoveLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgRemoveLiquidity) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func NewMsgAddLiquidity(signer sdk.AccAddress, externalAsset Asset, nativeAssetAmount sdk.Uint, externalAssetAmount sdk.Uint) MsgAddLiquidity {
	return MsgAddLiquidity{Signer: signer.String(), ExternalAsset: &externalAsset, NativeAssetAmount: nativeAssetAmount, ExternalAssetAmount: externalAssetAmount}
}

func (m MsgAddLiquidity) Route() string {
	return RouterKey
}

func (m MsgAddLiquidity) Type() string {
	return "add_liquidity"
}

func (m MsgAddLiquidity) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	if !m.ExternalAsset.Validate() {
		return sdkerrors.Wrap(ErrInValidAsset, m.ExternalAsset.Symbol)
	}
	if m.ExternalAsset.Equals(GetSettlementAsset()) {
		return sdkerrors.Wrap(ErrInValidAsset, "External asset cannot be rowan")
	}
	if !(m.NativeAssetAmount.GTE(sdk.ZeroUint())) && (m.ExternalAssetAmount.GTE(sdk.ZeroUint())) {
		return sdkerrors.Wrap(ErrInValidAmount, fmt.Sprintf("Both asset ammounts cannot be 0 %s / %s", m.NativeAssetAmount.String(), m.ExternalAssetAmount.String()))
	}

	return nil
}

func (m MsgAddLiquidity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgAddLiquidity) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

// type MsgCreatePool struct {
// 	Signer              sdk.AccAddress
// 	ExternalAsset       Asset
// 	NativeAssetAmount   sdk.Uint
// 	ExternalAssetAmount sdk.Uint
// }

func NewMsgCreatePool(signer sdk.AccAddress, externalAsset Asset, nativeAssetAmount sdk.Uint, externalAssetAmount sdk.Uint) MsgCreatePool {
	return MsgCreatePool{Signer: signer.String(), ExternalAsset: &externalAsset, NativeAssetAmount: nativeAssetAmount, ExternalAssetAmount: externalAssetAmount}
}

func (m MsgCreatePool) Route() string {
	return RouterKey
}

func (m MsgCreatePool) Type() string {
	return "create_pool"
}

func (m MsgCreatePool) ValidateBasic() error {
	if len(m.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}
	if !m.ExternalAsset.Validate() {
		return sdkerrors.Wrap(ErrInValidAsset, m.ExternalAsset.Symbol)
	}
	if m.ExternalAsset.Equals(GetSettlementAsset()) {
		return sdkerrors.Wrap(ErrInValidAsset, "External Asset cannot be rowan")
	}
	if !(m.NativeAssetAmount.GT(sdk.ZeroUint())) {
		return sdkerrors.Wrap(ErrInValidAmount, m.NativeAssetAmount.String())
	}
	if !(m.ExternalAssetAmount.GT(sdk.ZeroUint())) {
		return sdkerrors.Wrap(ErrInValidAmount, m.NativeAssetAmount.String())
	}
	return nil
}

func (m MsgCreatePool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgCreatePool) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
