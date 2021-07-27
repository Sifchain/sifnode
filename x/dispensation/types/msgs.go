package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/pkg/errors"
)

func NewMsgCreateDistribution(distributor sdk.AccAddress, DistributionType DistributionType, output []types.Output, authorizedRunner string) MsgCreateDistribution {

	return MsgCreateDistribution{
		Distributor:      distributor.String(),
		AuthorizedRunner: authorizedRunner,
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
	_, ok := IsValidDistributionType(m.DistributionType.String())
	if !ok {
		return sdkerrors.Wrap(ErrInvalid, "Invalid Distribution Type")
	}
	// Validate length of output is not 0
	if len(m.Output) == 0 {
		return errors.Wrapf(ErrInvalid, "Outputlist cannot be empty")
	}
	_, err := sdk.AccAddressFromBech32(m.Distributor)
	if err != nil {
		return errors.Wrapf(ErrInvalid, "Invalid Distributor Address : %s", m.Distributor)
	}
	// Validator runner
	_, err = sdk.AccAddressFromBech32(m.AuthorizedRunner)
	if err != nil {
		return errors.Wrapf(ErrInvalid, "Invalid Authorized Address : %s", m.AuthorizedRunner)
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
	_, ok := IsValidClaimType(m.UserClaimType.String())
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

func NewMsgRunDistribution(runner string, distributionName string, distributionType DistributionType) MsgRunDistribution {
	return MsgRunDistribution{
		AuthorizedRunner: runner,
		DistributionName: distributionName,
		DistributionType: distributionType,
	}
}

func (m MsgRunDistribution) Route() string {
	return RouterKey
}

func (m MsgRunDistribution) Type() string {
	return MsgTypeRunDistribution
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
	_, err := sdk.AccAddressFromBech32(m.AuthorizedRunner)
	if err != nil {
		return errors.Wrapf(ErrInvalid, "Invalid Runner Address")
	}
	return nil
}

func (m MsgRunDistribution) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgRunDistribution) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.AuthorizedRunner)
	// Should never panic as ValidateBasic checks address validity
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}
