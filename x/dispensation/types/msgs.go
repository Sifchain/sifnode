package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/pkg/errors"
)

var (
	_ sdk.Msg = &MsgDistribution{}
)

// Basic message type to create a new distribution
// TODO modify this struct to keep adding more fields to identify different types of distributions
type MsgDistribution struct {
	Distributor      sdk.AccAddress   `json:"distributor"`
	DistributionName string           `json:"distribution_name"`
	DistributionType DistributionType `json:"distribution_type"`
	Output           []bank.Output    `json:"Output"`
}

func NewMsgDistribution(signer sdk.AccAddress, DistributionName string, DistributionType DistributionType, output []bank.Output) MsgDistribution {
	return MsgDistribution{Distributor: signer, DistributionName: DistributionName, DistributionType: DistributionType, Output: output}
}

func (m MsgDistribution) Route() string {
	return RouterKey
}

//TODO Replace with constant defined in keys.go with value CreateDispensation
func (m MsgDistribution) Type() string {
	return "airdrop"
}

func (m MsgDistribution) ValidateBasic() error {
	if m.DistributionName == "" {
		return sdkerrors.Wrap(ErrInvalid, "Name cannot be empty")
	}
	for _, out := range m.Output {
		if !out.Coins.IsValid() {
			return errors.Wrapf(ErrInvalid, "Invalid Coins")
		}
		if len(out.Coins) > 1 {
			return errors.Wrapf(ErrInvalid, "Invalid Coins Can only specify one coin type for an entry")
		}
		if out.Coins.GetDenomByIndex(0) != TokenSupported {
			return errors.Wrapf(ErrInvalid, "Invalid Coins Specified coin can only be rowan")
		}
	}
	return nil
}

func (m MsgDistribution) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgDistribution) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Distributor}
}

// Create a user claim
type MsgCreateClaim struct {
	UserClaimAddress sdk.AccAddress   `json:"user_claim_address"`
	UserClaimType    DistributionType `json:"user_claim_type"`
}

func NewMsgCreateClaim(Signer sdk.AccAddress, userClaimType DistributionType) MsgCreateClaim {
	return MsgCreateClaim{UserClaimAddress: Signer, UserClaimType: userClaimType}
}

func (m MsgCreateClaim) Route() string {
	return RouterKey
}

func (m MsgCreateClaim) Type() string {
	return "createClaim"
}

// Validation for claim type
func (m MsgCreateClaim) ValidateBasic() error {
	if m.UserClaimAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.UserClaimAddress.String())
	}
	_, ok := IsValidClaim(m.UserClaimType.String())
	if !ok {
		return sdkerrors.Wrap(ErrInvalid, m.UserClaimType.String())
	}
	return nil
}

func (m MsgCreateClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgCreateClaim) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.UserClaimAddress}
}
