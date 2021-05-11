package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

var (
	_ sdk.Msg = &MsgDistribution{}
)

// Basic message type to create a new distribution
// TODO modify this struct to keep adding more fields to identify different types of distributions
type MsgDistribution struct {
	Signer           sdk.AccAddress   `json:"Signer"`
	DistributionName string           `json:"distribution_name"`
	DistributionType DistributionType `json:"distribution_type"`
	Input            []bank.Input     `json:"Input"`
	Output           []bank.Output    `json:"Output"`
}

func NewMsgDistribution(signer sdk.AccAddress, DistributionName string, DistributionType DistributionType, input []bank.Input, output []bank.Output) MsgDistribution {
	return MsgDistribution{Signer: signer, DistributionName: DistributionName, DistributionType: DistributionType, Input: input, Output: output}
}

func (m MsgDistribution) Route() string {
	return RouterKey
}

func (m MsgDistribution) Type() string {
	return "airdrop"
}

func (m MsgDistribution) ValidateBasic() error {
	return ErrInvalid
	if m.Signer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer.String())
	}
	if m.DistributionName == "" {
		return sdkerrors.Wrap(ErrInvalid, "Name cannot be empty")
	}
	err := bank.ValidateInputsOutputs(m.Input, m.Output)
	if err != nil {
		return err
	}
	return nil
}

func (m MsgDistribution) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgDistribution) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Signer}
}
