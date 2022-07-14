package types_test

import (
	"strings"
	"testing"

	sifapp "github.com/Sifchain/sifnode/app"
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

func TestMsgCreateDistribution(t *testing.T) {
	distributor := sdk.AccAddress("addr1_______________")
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_AIRDROP
	output := test.CreatOutputList(2000, "1")
	authorizedRunner := sdk.AccAddress("addr2_______________")
	result := types.NewMsgCreateDistribution(distributor, distributionType, output, authorizedRunner.String())

	assert.Equal(t, distributor.String(), result.Distributor)
	assert.Equal(t, distributionType, result.DistributionType)
	assert.Equal(t, output, result.Output)
	assert.Equal(t, authorizedRunner.String(), result.GetAuthorizedRunner())

}

func TestMsgCreateDistribution_Route(t *testing.T) {
	distributor := sdk.AccAddress("addr1_______________")
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_AIRDROP
	output := test.CreatOutputList(2000, "1")
	authorizedRunner := sdk.AccAddress("addr2_______________")
	result := types.NewMsgCreateDistribution(distributor, distributionType, output, authorizedRunner.String())

	key := result.Route()
	key1 := types.RouterKey
	assert.Equal(t, key, key1)
}

func TestMsgCreateDistribution_Type(t *testing.T) {
	distributor := sdk.AccAddress("addr1_______________")
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_AIRDROP
	output := test.CreatOutputList(2000, "1")
	authorizedRunner := sdk.AccAddress("addr2_______________")
	result := types.NewMsgCreateDistribution(distributor, distributionType, output, authorizedRunner.String())

	typ := result.Type()
	typ1 := types.MsgTypeCreateDistribution
	assert.Equal(t, typ, typ1)
}

func TestMsgCreateDistribution_GetSigners(t *testing.T) {
	sifapp.SetConfig(false)
	distributor := sdk.AccAddress("addr1_______________")
	authorizedRunner := sdk.AccAddress("addr2_______________")
	msg := types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		Output:           test.CreatOutputList(2000, "1"),
		AuthorizedRunner: authorizedRunner.String(),
	}
	result := msg.GetSigners()
	addr, err := sdk.AccAddressFromBech32(msg.Distributor)
	output := []sdk.AccAddress{addr}
	assert.Equal(t, result, output)
	assert.NoError(t, err)

}

func TestMsgCreateDistribution_GetSignBytes(t *testing.T) {
	distributor := sdk.AccAddress("addr1_______________")
	authorizedRunner := sdk.AccAddress("addr2_______________")
	msg := types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		Output:           test.CreatOutputList(2000, "1"),
		AuthorizedRunner: authorizedRunner.String(),
	}
	byt := msg.GetSignBytes()
	byt1 := sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(&msg))
	assert.Equal(t, byt, byt1)

}

func TestMsgCreateDistribution_ValidateBasic(t *testing.T) {
	sifapp.SetConfig(false)
	distributor := sdk.AccAddress("addr1_______________")
	authorizedRunner := sdk.AccAddress("addr2_______________")
	msg := types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		Output:           test.CreatOutputList(2000, "1"),
		AuthorizedRunner: authorizedRunner.String(),
	}
	err := msg.ValidateBasic()
	assert.NoError(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_InvalidDistributionType(t *testing.T) {
	sifapp.SetConfig(false)
	distributor := sdk.AccAddress("addr1_______________")
	authorizedRunner := sdk.AccAddress("addr2_______________")
	msg := types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED,
		Output:           test.CreatOutputList(2000, "1"),
		AuthorizedRunner: authorizedRunner.String(),
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_Order(t *testing.T) {
	sifapp.SetConfig(false)
	distributor := sdk.AccAddress("addr1_______________")
	authorizedRunner := sdk.AccAddress("addr2_______________")
	validAddress := sdk.AccAddress("addr3_______________")

	coin := []sdk.Coin{sdk.NewCoin("rowan", sdk.NewInt(1000000)), sdk.NewCoin("cusdt", sdk.NewInt(1000000))}

	output := banktypes.NewOutput(validAddress, sdk.NewCoins(coin...))
	outputList := []banktypes.Output{output}
	msg := types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		Output:           outputList,
		AuthorizedRunner: authorizedRunner.String(),
	}
	err := msg.ValidateBasic()
	assert.NoError(t, err)
	for i, j := 0, len(output.Coins)-1; i < j; i, j = i+1, j-1 {
		output.Coins[i], output.Coins[j] = output.Coins[j], output.Coins[i]
	}
	revOutputList := []banktypes.Output{output}
	msg = types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		Output:           revOutputList,
		AuthorizedRunner: authorizedRunner.String(),
	}
	err = msg.ValidateBasic()
	assert.NoError(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_WrongAddress(t *testing.T) {
	distributor := sdk.AccAddress("addr1_______________")
	outputList := test.CreatOutputList(10, "1")
	validAddress := sdk.AccAddress("addr2_______________")
	authorizedRunner := sdk.AccAddress("addr3_______________")
	// Address is valid as long as its length is between 1-255 bytes.
	invalidAddress := sdk.AccAddress("")
	outputList = append(outputList, banktypes.NewOutput(invalidAddress, sdk.NewCoins(sdk.NewCoin("rowan", sdk.NewInt(1000000)))))
	msg := types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		Output:           outputList,
		AuthorizedRunner: authorizedRunner.String(),
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
	invalidAddress2 := sdk.AccAddress(strings.Repeat(validAddress.String(), 7))
	outputList = append(outputList, banktypes.NewOutput(invalidAddress2, sdk.NewCoins(sdk.NewCoin("rowan", sdk.NewInt(1000000)))))
	msg = types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		Output:           outputList,
		AuthorizedRunner: authorizedRunner.String(),
	}
	err = msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_InvalidDistributionAddress(t *testing.T) {

	distributor := sdk.AccAddress("")
	authorizedRunner := sdk.AccAddress("addr2_______________")
	msg := types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		Output:           test.CreatOutputList(2000, "1"),
		AuthorizedRunner: authorizedRunner.String(),
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_InvalidauthorizedAddress(t *testing.T) {

	distributor := sdk.AccAddress("addr1_______________")
	authorizedRunner := sdk.AccAddress("")
	msg := types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		Output:           test.CreatOutputList(2000, "1"),
		AuthorizedRunner: authorizedRunner.String(),
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_NonRowan(t *testing.T) {
	distributor := sdk.AccAddress([]byte("addr1_______________"))
	outputlist := test.CreatOutputList(2000, "1")
	authorizedRunner := sdk.AccAddress("addr2_______________")
	outputlist = append(outputlist, banktypes.NewOutput(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()),
		sdk.NewCoins(sdk.NewCoin("dash", sdk.NewInt(10)))))
	msg := types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		Output:           outputlist,
		AuthorizedRunner: authorizedRunner.String(),
	}
	err := msg.ValidateBasic()
	assert.NoError(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_MultipleCoins(t *testing.T) {
	distributor := sdk.AccAddress("addr1_______________")
	outputlist := test.CreatOutputList(2000, "1")
	authorizedRunner := sdk.AccAddress("addr2_______________")
	outputlist = append(outputlist, banktypes.NewOutput(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()),
		sdk.NewCoins(sdk.NewCoin("rowan", sdk.NewInt(10)), sdk.NewCoin("dash", sdk.NewInt(100)))))
	msg := types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		Output:           outputlist,
		AuthorizedRunner: authorizedRunner.String(),
	}
	err := msg.ValidateBasic()
	assert.NoError(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_ZeroCoins(t *testing.T) {
	distributor := sdk.AccAddress("addr1_______________")
	authorizedRunner := sdk.AccAddress("addr2_______________")
	var outputlist []banktypes.Output
	msg := types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		Output:           outputlist,
		AuthorizedRunner: authorizedRunner.String(),
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateDistribution_ValidateBasic_invalidCoins(t *testing.T) {
	distributor := sdk.AccAddress("addr1_______________")
	outputlist := test.CreatOutputList(2000, "1")
	authorizedRunner := sdk.AccAddress("addr2_______________")
	outputlist = append(outputlist, banktypes.NewOutput(sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address()),
		sdk.Coins{sdk.Coin{Denom: sdk.ErrInvalidLengthCoin.Error(),
			Amount: sdk.NewInt(20)}}))
	msg := types.MsgCreateDistribution{
		Distributor:      distributor.String(),
		DistributionType: types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
		Output:           outputlist,
		AuthorizedRunner: authorizedRunner.String(),
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateClaim_ValidateBasic_WrongType(t *testing.T) {
	claimer := sdk.AccAddress("addr1_______________")
	msg := types.MsgCreateUserClaim{
		UserClaimAddress: claimer.String(),
		UserClaimType:    types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgCreateClaim_ValidateBasic_InvalidAddress(t *testing.T) {
	claimer := sdk.AccAddress("")
	msg := types.MsgCreateUserClaim{
		UserClaimAddress: claimer.String(),
		UserClaimType:    types.DistributionType_DISTRIBUTION_TYPE_AIRDROP,
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}
func TestMsgCreateClaim_ValidateBasic(t *testing.T) {
	claimer := sdk.AccAddress("addr1_______________")
	msg := types.MsgCreateUserClaim{
		UserClaimAddress: claimer.String(),
		UserClaimType:    types.DistributionType_DISTRIBUTION_TYPE_VALIDATOR_SUBSIDY,
	}
	err := msg.ValidateBasic()
	assert.NoError(t, err)
}

func TestMsgCreateUserClaim(t *testing.T) {
	userClaimAddress := sdk.AccAddress("addr2_______________")
	claimType := types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED
	result := types.NewMsgCreateUserClaim(userClaimAddress, claimType)

	assert.Equal(t, userClaimAddress.String(), result.UserClaimAddress)
	assert.Equal(t, claimType, result.UserClaimType)
}

func TestMsgCreateUserClaim_Route(t *testing.T) {
	userClaimAddress := sdk.AccAddress("addr2_______________")
	claimType := types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED
	result := types.NewMsgCreateUserClaim(userClaimAddress, claimType)

	key := result.Route()
	key1 := types.RouterKey
	assert.Equal(t, key, key1)
}

func TestMsgCreateUserClaim_Type(t *testing.T) {
	userClaimAddress := sdk.AccAddress("addr2_______________")
	claimType := types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED
	result := types.NewMsgCreateUserClaim(userClaimAddress, claimType)

	typ := result.Type()
	typ1 := types.MsgTypeCreateUserClaim
	assert.Equal(t, typ, typ1)
}

func TestMsgCreateUserClaim_GetSignBytes(t *testing.T) {
	addr := sdk.AccAddress("addr2_______________")
	msg := types.MsgCreateUserClaim{
		UserClaimAddress: addr.String(),
		UserClaimType:    types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED,
	}
	byt := msg.GetSignBytes()
	byt1 := sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(&msg))
	assert.Equal(t, byt, byt1)

}

func TestMsgCreateUserClaim_GetSigners(t *testing.T) {
	addr := sdk.AccAddress("addr2_______________")
	msg := types.MsgCreateUserClaim{
		UserClaimAddress: addr.String(),
		UserClaimType:    types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED,
	}
	result := msg.GetSigners()
	addr, err := sdk.AccAddressFromBech32(msg.UserClaimAddress)
	output := []sdk.AccAddress{addr}
	assert.Equal(t, result, output)
	assert.NoError(t, err)

}

func TestMsgRunDistribution(t *testing.T) {
	runner := sdk.AccAddress("addr2_______________")
	distributionName := types.AttributeKeyDistributionName
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_AIRDROP
	result := types.NewMsgRunDistribution(runner.String(), distributionName, distributionType, 10)

	assert.Equal(t, runner.String(), result.GetAuthorizedRunner())
	assert.Equal(t, distributionName, result.DistributionName)
	assert.Equal(t, distributionType, result.DistributionType)

}

func TestMsgRunDistribution_Route(t *testing.T) {
	runner := sdk.AccAddress("addr2_______________")
	distributionName := types.AttributeKeyDistributionName
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_AIRDROP
	result := types.NewMsgRunDistribution(runner.String(), distributionName, distributionType, 10)

	key := result.Route()
	key1 := types.RouterKey
	assert.Equal(t, key, key1)

}

func TestMsgRunDistribution_Type(t *testing.T) {
	runner := sdk.AccAddress("addr2_______________")
	distributionName := types.AttributeKeyDistributionName
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_AIRDROP
	result := types.NewMsgRunDistribution(runner.String(), distributionName, distributionType, 10)

	typ := result.Type()
	typ1 := types.MsgTypeRunDistribution
	assert.Equal(t, typ, typ1)
}

func TestMsgRunDistribution_validateBasic(t *testing.T) {
	runner := sdk.AccAddress("addr2_______________")
	distributionName := types.AttributeKeyDistributionName
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_AIRDROP

	msg := types.MsgRunDistribution{
		AuthorizedRunner:  runner.String(),
		DistributionName:  distributionName,
		DistributionType:  distributionType,
		DistributionCount: 10,
	}
	err := msg.ValidateBasic()
	assert.NoError(t, err)
}

func TestMsgRunDistribution_dispensationCount(t *testing.T) {
	runner := sdk.AccAddress("addr2_______________")
	distributionName := types.AttributeKeyDistributionName
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_AIRDROP

	msg := types.MsgRunDistribution{
		AuthorizedRunner:  runner.String(),
		DistributionName:  distributionName,
		DistributionType:  distributionType,
		DistributionCount: -1,
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgRunDistribution_validateBasic_InvalidDistributiontype(t *testing.T) {
	runner := sdk.AccAddress("addr2_______________")
	distributionName := types.AttributeKeyDistributionName
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_UNSPECIFIED

	msg := types.MsgRunDistribution{
		AuthorizedRunner:  runner.String(),
		DistributionName:  distributionName,
		DistributionType:  distributionType,
		DistributionCount: 10,
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgRunDistribution_validateBasic_EmptyDistributionName(t *testing.T) {
	runner := sdk.AccAddress("addr2_______________")
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_AIRDROP

	msg := types.MsgRunDistribution{
		AuthorizedRunner:  runner.String(),
		DistributionName:  "",
		DistributionType:  distributionType,
		DistributionCount: 10,
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgRunDistribution_validateBasic_InvalidRunnerAddress(t *testing.T) {
	distributionName := types.AttributeKeyDistributionName
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_AIRDROP

	msg := types.MsgRunDistribution{
		AuthorizedRunner:  types.AttributeKeyDistributionRunner,
		DistributionName:  distributionName,
		DistributionType:  distributionType,
		DistributionCount: 10,
	}
	err := msg.ValidateBasic()
	assert.Error(t, err)
}

func TestMsgRunDistribution_GetSignBytes(t *testing.T) {
	runner := sdk.AccAddress("addr2_______________")
	distributionName := types.AttributeKeyDistributionName
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_AIRDROP

	msg := types.MsgRunDistribution{
		AuthorizedRunner:  runner.String(),
		DistributionName:  distributionName,
		DistributionType:  distributionType,
		DistributionCount: 10,
	}
	byt := msg.GetSignBytes()
	byt1 := sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(&msg))
	assert.Equal(t, byt, byt1)

}

func TestMsgRunDistribution_GetSigners(t *testing.T) {
	sifapp.SetConfig(false)
	runner := sdk.AccAddress("addr2_______________")
	distributionName := types.AttributeKeyDistributionName
	distributionType := types.DistributionType_DISTRIBUTION_TYPE_AIRDROP

	msg := types.MsgRunDistribution{
		AuthorizedRunner:  runner.String(),
		DistributionName:  distributionName,
		DistributionType:  distributionType,
		DistributionCount: 10,
	}
	result := msg.GetSigners()
	addr, err := sdk.AccAddressFromBech32(msg.AuthorizedRunner)
	output := []sdk.AccAddress{addr}
	assert.Equal(t, result, output)
	assert.NoError(t, err)
}
