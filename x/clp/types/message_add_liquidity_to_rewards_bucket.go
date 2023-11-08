package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddLiquidityToRewardsBucketRequest = "add_liquidity_to_rewards_bucket"

var _ sdk.Msg = &MsgAddLiquidityToRewardsBucketRequest{}

func NewMsgAddLiquidityToRewardsBucketRequest(signer string, amount sdk.Coins) *MsgAddLiquidityToRewardsBucketRequest {
	return &MsgAddLiquidityToRewardsBucketRequest{
		Signer: signer,
		Amount: amount,
	}
}

func (msg *MsgAddLiquidityToRewardsBucketRequest) Route() string {
	return RouterKey
}

func (msg *MsgAddLiquidityToRewardsBucketRequest) Type() string {
	return TypeMsgAddLiquidityToRewardsBucketRequest
}

func (msg *MsgAddLiquidityToRewardsBucketRequest) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{signer}
}

func (msg *MsgAddLiquidityToRewardsBucketRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddLiquidityToRewardsBucketRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid signer address (%s)", err)
	}
	return nil
}
