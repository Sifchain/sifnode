package types

import (

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

)

var (
	_ sdk.Msg = &MsgAirdrop{}
)


type MsgAirdrop struct {
	Signer              sdk.AccAddress
	NativeAssetAmount   sdk.Uint
	ExternalAssetAmount sdk.Uint
}

func NewMsgAirdrop(signer sdk.AccAddress, nativeAssetAmount sdk.Uint, externalAssetAmount sdk.Uint) MsgAirdrop {
	return MsgAirdrop{Signer: signer, NativeAssetAmount: nativeAssetAmount, ExternalAssetAmount: externalAssetAmount}
}

func (m MsgAirdrop) Route() string {
	return RouterKey
}

func (m MsgAirdrop) Type() string {
	return "create_pool"
}

func (m MsgAirdrop) ValidateBasic() error {
	if m.Signer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer.String())
	}
	return nil
}

func (m MsgAirdrop) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgAirdrop) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Signer}
}
