package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/pkg/errors"
)

func NewMsgCreateDistribution(distributor sdk.AccAddress, DistributionName string, DistributionType DistributionType, output []types.Output) MsgCreateDistribution {

	return MsgCreateDistribution{
		Distributor: distributor.String(),
		Distribution: &Distribution{
			DistributionName: DistributionName,
			DistributionType: DistributionType,
		},
		Output: output,
	}
}

func (m MsgCreateDistribution) Route() string {
	return RouterKey
}

func (m MsgCreateDistribution) Type() string {
	return MsgTypeCreateDistribution
}

func (m MsgCreateDistribution) ValidateBasic() error {
	if m.Distribution.DistributionName == "" {
		return sdkerrors.Wrap(ErrInvalid, "Name cannot be empty")
	}
	if len(m.Output) == 0 {
		return errors.Wrapf(ErrInvalid, "Outputlist cannot be empty")
	}
	for _, out := range m.Output {
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
