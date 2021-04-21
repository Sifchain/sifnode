package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
)


func NewMsgDistribution(
	signer sdk.AccAddress,
	DistributionName string,
	DistributionType DistributionType,
	input []types.Input,
	output []types.Output) MsgDistribution {

	return MsgDistribution{
		Signer: signer.String(),
		DistributionName: DistributionName,
		DistributionType: DistributionType,
		Input: input,
		Output: output,
	}
}

func (m MsgDistribution) Route() string {
	return RouterKey
}

func (m MsgDistribution) Type() string {
	return "airdrop"
}

func (m MsgDistribution) ValidateBasic() error {
	if m.Signer == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer)
	}

	if m.DistributionName == "" {
		return sdkerrors.Wrap(ErrInvalid, "Name cannot be empty")
	}

	err := types.ValidateInputsOutputs(m.Input, m.Output)
	if err != nil {
		return err
	}

	return nil
}

func (m MsgDistribution) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgDistribution) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{addr}
}
