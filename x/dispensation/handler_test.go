package dispensation_test

import (
	"github.com/Sifchain/sifnode/x/dispensation"
	"github.com/Sifchain/sifnode/x/dispensation/test"
	"github.com/Sifchain/sifnode/x/dispensation/types"
	dispensationUtils "github.com/Sifchain/sifnode/x/dispensation/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func TestNewHandler_CreateDistribution(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	handler := dispensation.NewHandler(keeper)
	recipients := 3000
	outputList := test.CreatOutputList(recipients, "10000000000000000000")
	distibutor := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	totalCoins := dispensationUtils.TotalOutput(outputList)
	_, err := keeper.GetBankKeeper().AddCoins(ctx, distibutor, totalCoins)
	assert.NoError(t, err)
	distributionname := "AR1"
	msgAirdrop := types.NewMsgDistribution(distibutor, distributionname, types.Airdrop, outputList)
	res, err := handler(ctx, msgAirdrop)
	require.NoError(t, err)
	require.NotNil(t, res)
	for _, e := range res.Events {
		if e.Type == "distribution_started" {
			assert.Len(t, e.Attributes, 3)
			assert.Contains(t, e.Attributes[1].String(), "distribution_name")
			assert.Contains(t, e.Attributes[1].String(), distributionname)
			assert.Contains(t, e.Attributes[2].String(), "distribution_type")
			assert.Contains(t, e.Attributes[2].String(), types.Airdrop.String())
		}

	}
	dr := keeper.GetRecordsForNameAll(ctx, "AR1")
	assert.Len(t, dr, recipients)
}

func TestNewHandler_CreateClaim(t *testing.T) {
	app, ctx := test.CreateTestApp(false)
	keeper := app.DispensationKeeper
	handler := dispensation.NewHandler(keeper)
	address := sdk.AccAddress(crypto.AddressHash([]byte("User1")))
	msgClaim := types.NewMsgCreateClaim(address, types.ValidatorSubsidy)
	res, err := handler(ctx, msgClaim)
	require.NoError(t, err)
	require.NotNil(t, res)

	cl, err := keeper.GetClaim(ctx, address.String(), types.ValidatorSubsidy)
	require.NoError(t, err)
	assert.False(t, cl.Locked)
	assert.Equal(t, cl.UserAddress.String(), address.String())
}
