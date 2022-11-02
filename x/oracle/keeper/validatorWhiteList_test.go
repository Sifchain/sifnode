package keeper_test

import (
	"bytes"
	"testing"

	"github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/stretchr/testify/assert"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sifapp "github.com/Sifchain/sifnode/app"
)

func TestKeeper_SetValidatorWhiteList(t *testing.T) {
	app := sifapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addresses := sifapp.CreateRandomAccounts(2)
	valAddresses := sifapp.ConvertAddrsToValAddrs(addresses)
	networkDescriptor := types.NewNetworkIdentity(types.NetworkDescriptor(0))

	for _, address := range valAddresses {
		err := app.OracleKeeper.UpdateOracleWhiteList(ctx, types.NetworkDescriptor(0), address, 10)
		if err != nil {
			panic(err)
		}
	}

	vList := app.OracleKeeper.GetOracleWhiteList(ctx, networkDescriptor)
	assert.Equal(t, len(vList.ValidatorPower), 2)
	assert.True(t, app.OracleKeeper.ExistsOracleWhiteList(ctx, networkDescriptor))
}

func TestKeeper_ValidateAddress(t *testing.T) {
	app := sifapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addresses := sifapp.CreateRandomAccounts(2)
	valAddresses := sifapp.ConvertAddrsToValAddrs(addresses)
	networkDescriptor := types.NewNetworkIdentity(types.NetworkDescriptor(0))

	for _, address := range valAddresses {
		err := app.OracleKeeper.UpdateOracleWhiteList(ctx, types.NetworkDescriptor(0), address, 10)
		if err != nil {
			panic(err)
		}
	}

	assert.True(t, app.OracleKeeper.ValidateAddress(ctx, networkDescriptor, valAddresses[0]))
	assert.True(t, app.OracleKeeper.ValidateAddress(ctx, networkDescriptor, valAddresses[1]))
	addresses = sifapp.CreateRandomAccounts(3)
	valAddresses = sifapp.ConvertAddrsToValAddrs(addresses)
	assert.False(t, app.OracleKeeper.ValidateAddress(ctx, networkDescriptor, valAddresses[2]))
}

func TestKeeper_GetAllWhiteList(t *testing.T) {
	app := sifapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	addresses := sifapp.CreateRandomAccounts(2)
	valAddresses := sifapp.ConvertAddrsToValAddrs(addresses)
	whitelist := make([]*types.ValidatorPower, 0)
	for _, address := range valAddresses {
		whitelist = append(whitelist, &types.ValidatorPower{
			ValidatorAddress: address,
			VotingPower:      100,
		})
		err := app.OracleKeeper.UpdateOracleWhiteList(ctx, types.NetworkDescriptor(0), address, 100)
		assert.NoError(t, err)
	}

	allWhiteList := app.OracleKeeper.GetAllWhiteList(ctx)

	expectedWhitelist := make([]*types.ValidatorPower, 0)
	found := false

	for _, value := range allWhiteList {
		if value.NetworkDescriptor == types.NetworkDescriptor(0) {
			found = true
			expectedWhitelist = value.ValidatorWhitelist.ValidatorPower
		}
	}
	assert.Equal(t, found, true)

	for _, value := range whitelist {
		found := false
		for _, expected := range expectedWhitelist {
			if bytes.Compare(value.ValidatorAddress, expected.ValidatorAddress) == 0 {
				found = true
				assert.Equal(t, value.VotingPower, expected.VotingPower)
			}
		}
		assert.Equal(t, found, true)
	}
}
