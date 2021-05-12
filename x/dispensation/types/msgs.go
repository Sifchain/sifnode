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

//TODO Replace with constant defined in keys.go with value CreateDispensation
func (m MsgDistribution) Type() string {
	return "airdrop"
}

func (m MsgDistribution) ValidateBasic() error {
	return ErrInvalid
}

func (m MsgDistribution) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgDistribution) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Signer}
}

// Create a user claim
type MsgCreateClaim struct {
	Signer        sdk.AccAddress   `json:"signer"`
	UserClaimType DistributionType `json:"user_claim_type"`
}

func NewMsgCreateClaim(Signer sdk.AccAddress, userClaimType DistributionType) MsgCreateClaim {
	return MsgCreateClaim{Signer: Signer, UserClaimType: userClaimType}
}

func (m MsgCreateClaim) Route() string {
	return RouterKey
}

func (m MsgCreateClaim) Type() string {
	return "createClaim"
}

func (m MsgCreateClaim) ValidateBasic() error {
	if m.Signer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Signer.String())
	}
	return nil
}

func (m MsgCreateClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgCreateClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Signer}
}
