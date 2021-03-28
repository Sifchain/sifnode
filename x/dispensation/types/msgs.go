package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

var (
	_ sdk.Msg = &MsgDistribution{}
)

type MsgDistribution struct {
	Signer           sdk.AccAddress `json:"Signer"`
	DistributionName string         `json:"distribution_name"`
	Input            []bank.Input   `json:"Input"`
	Output           []bank.Output  `json:"Output"`
}

func NewMsgDistribution(signer sdk.AccAddress, name string, input []bank.Input, output []bank.Output) MsgDistribution {
	return MsgDistribution{Signer: signer, DistributionName: name, Input: input, Output: output}
}

func (m MsgDistribution) Route() string {
	return RouterKey
}

func (m MsgDistribution) Type() string {
	return "airdrop"
}

func (m MsgDistribution) ValidateBasic() error {
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
