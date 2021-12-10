package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/stretchr/testify/assert"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sifapp "github.com/Sifchain/sifnode/app"
)

func TestKeeper_GetConsensusNeeded(t *testing.T) {
	app := sifapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	// test against wrong network identity
	networkDescriptor := types.NewNetworkIdentity(types.NetworkDescriptor_NETWORK_DESCRIPTOR_UNSPECIFIED)
	_, err := app.OracleKeeper.GetCrossChainFee(ctx, networkDescriptor)
	assert.Error(t, err)

	// test if ConsensusNeeded not set for ethereum
	networkDescriptor = types.NewNetworkIdentity(types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM)
	_, err = app.OracleKeeper.GetConsensusNeeded(ctx, networkDescriptor)
	assert.Error(t, err)

	app.OracleKeeper.SetConsensusNeeded(ctx, networkDescriptor, 10)

	// case for well set the ConsensusNeeded
	consensusNeeded, err := app.OracleKeeper.GetConsensusNeeded(ctx, networkDescriptor)
	assert.NoError(t, err)
	assert.Equal(t, consensusNeeded, uint32(10))
}

func TestKeeper_SetConsensusNeeded(t *testing.T) {
	app := sifapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	networkDescriptor := types.NewNetworkIdentity(types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM)

	app.OracleKeeper.SetConsensusNeeded(ctx, networkDescriptor, 10)

	consensusNeeded, err := app.OracleKeeper.GetConsensusNeeded(ctx, networkDescriptor)
	assert.NoError(t, err)
	assert.Equal(t, uint32(10), consensusNeeded)
}

// test if the value is too large
func TestKeeper_SetConsensusNeededWithLargeNumber(t *testing.T) {
	app := sifapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	networkDescriptor := types.NewNetworkIdentity(types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM)

	app.OracleKeeper.SetConsensusNeeded(ctx, networkDescriptor, 1000)

	_, err := app.OracleKeeper.GetConsensusNeeded(ctx, networkDescriptor)
	assert.Error(t, err)
}
