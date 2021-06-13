package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/pkg/errors"
)

var (
	_ sdk.Msg = &MsgDistribution{}
	_ sdk.Msg = &MsgRunDistribution{}
	_ sdk.Msg = &MsgCreateClaim{}
)

// Basic message type to create a new distribution

type MsgDistribution struct {
	Distributor      sdk.AccAddress   `json:"distributor"`
	Runner           sdk.AccAddress   `json:"runner"`
	DistributionType DistributionType `json:"distribution_type"`
	Output           []bank.Output    `json:"output"`
}

func NewMsgDistribution(distributor sdk.AccAddress, DistributionType DistributionType, output []bank.Output, runner sdk.AccAddress) MsgDistribution {
	return MsgDistribution{Distributor: distributor, DistributionType: DistributionType, Output: output, Runner: runner}
}

func (m MsgDistribution) Route() string {
	return RouterKey
}

func (m MsgDistribution) Type() string {
	return "airdrop"
}

func (m MsgDistribution) ValidateBasic() error {
	// Validate distribution Type
	_, ok := IsValidDistributionType(m.DistributionType.String())
	if !ok {
		return sdkerrors.Wrap(ErrInvalid, "Invalid Distribution Type")
	}
	// Validate length of output is not 0
	if len(m.Output) == 0 {
		return errors.Wrapf(ErrInvalid, "Outputlist cannot be empty")
	}
	// Validator distributor
	_, err := sdk.AccAddressFromBech32(m.Distributor.String())
	if err != nil {
		return errors.Wrapf(ErrInvalid, "Invalid Distributor Address")
	}
	// Validator Runner
	_, err = sdk.AccAddressFromBech32(m.Runner.String())
	if err != nil {
		return errors.Wrapf(ErrInvalid, "Invalid Runner Address")
	}
	// Validate individual out records
	for _, out := range m.Output {
		_, err := sdk.AccAddressFromBech32(out.Address.String())
		if err != nil {
			return errors.Wrapf(ErrInvalid, "Invalid Recipient Address")
		}
		if !out.Coins.IsValid() {
			return errors.Wrapf(ErrInvalid, "Invalid Coins")
		}
		if len(out.Coins) > 1 {
			return errors.Wrapf(ErrInvalid, "Invalid Coins Can only specify one coin type for an entry")
		}
		if out.Coins.GetDenomByIndex(0) != TokenSupported {
			return errors.Wrapf(ErrInvalid, "Invalid Coins Specified coin can only be %s", TokenSupported)
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

// Run distribution

type MsgRunDistribution struct {
	DistributionRunner sdk.AccAddress   `json:"distribution_runner"`
	DistributionName   string           `json:"distribution_name"`
	DistributionType   DistributionType `json:"distribution_type"`
}

func NewMsgRunDistribution(runner sdk.AccAddress, distributionName string, distributionType DistributionType) MsgRunDistribution {
	return MsgRunDistribution{
		DistributionRunner: runner,
		DistributionName:   distributionName,
		DistributionType:   distributionType,
	}
}

func (m MsgRunDistribution) Route() string {
	return RouterKey
}

func (m MsgRunDistribution) Type() string {
	return "run_distribution"
}

func (m MsgRunDistribution) ValidateBasic() error {
	//Validate DistributionType
	_, ok := IsValidDistributionType(m.DistributionType.String())
	if !ok {
		return sdkerrors.Wrap(ErrInvalid, "Invalid Distribution Type")
	}
	// Validate distribution Name
	if m.DistributionName == "" {
		return sdkerrors.Wrap(ErrInvalid, m.DistributionName)
	}
	// Validator runner
	_, err := sdk.AccAddressFromBech32(m.DistributionRunner.String())
	if err != nil {
		return errors.Wrapf(ErrInvalid, "Invalid Runner Address")
	}
	return nil
}

func (m MsgRunDistribution) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgRunDistribution) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.DistributionRunner}
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
