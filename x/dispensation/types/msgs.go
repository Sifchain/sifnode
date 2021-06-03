package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/pkg/errors"
)

func NewMsgCreateDistribution(distributor sdk.AccAddress, DistributionType DistributionType, output []types.Output) MsgCreateDistribution {

	return MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: DistributionType,
		Output:           output,
	}
}

func (m MsgCreateDistribution) Route() string {
	return RouterKey
}

func (m MsgCreateDistribution) Type() string {
	return MsgTypeCreateDistribution
}

func (m MsgCreateDistribution) ValidateBasic() error {
	// Validate distribution Type
	_, ok := IsValidDistribution(m.DistributionType.String())
	if !ok {
		return sdkerrors.Wrap(ErrInvalid, "Invalid Distribution Type")
	}
	// Validate length of output is not 0
	if len(m.Output) == 0 {
		return errors.Wrapf(ErrInvalid, "Outputlist cannot be empty")
	}
	// Validator distributor
	_, err := sdk.AccAddressFromBech32(m.Distributor)
	if err != nil {
		return errors.Wrapf(ErrInvalid, "Invalid Distributor Address")
	}
	// Validate individual out records
	for _, out := range m.Output {
		_, err := sdk.AccAddressFromBech32(out.Address)
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

func NewMsgCreateUserClaim(userClaimAddress sdk.AccAddress, claimType DistributionType) MsgCreateUserClaim {
	return MsgCreateUserClaim{
		UserClaimAddress: userClaimAddress.String(),
		UserClaimType:    claimType,
	}
}
func (m MsgCreateDistribution) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgCreateDistribution) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Distributor)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{addr}
}

func (m MsgCreateUserClaim) Route() string {
	return RouterKey
}

func (m MsgCreateUserClaim) Type() string {
	return MsgTypeCreateUserClaim
}

func (m MsgCreateUserClaim) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.UserClaimAddress)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.UserClaimAddress)
	}
	_, ok := IsValidClaim(m.UserClaimType.String())
	if !ok {
		return sdkerrors.Wrap(ErrInvalid, m.UserClaimType.String())
	}
	return nil
}

func (m MsgCreateUserClaim) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgCreateUserClaim) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.UserClaimAddress)
	// Should never panic as ValidateBasic checks address validity
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
