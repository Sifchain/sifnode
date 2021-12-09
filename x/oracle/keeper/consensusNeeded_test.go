package keeper_test

import (
	"math"
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

	app.OracleKeeper.SetConsensusNeeded(ctx, networkDescriptor, float32(0.1))

	// case for well set the ConsensusNeeded
	consensusNeeded, err := app.OracleKeeper.GetConsensusNeeded(ctx, networkDescriptor)
	assert.NoError(t, err)
	diff := math.Abs(consensusNeeded - float64(0.1))
	assert.Equal(t, diff < 0.0001, true)
}

func TestKeeper_SetConsensusNeeded(t *testing.T) {
	app := sifapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	networkDescriptor := types.NewNetworkIdentity(types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM)

	app.OracleKeeper.SetConsensusNeeded(ctx, networkDescriptor, float32(0.1))

	consensusNeeded, err := app.OracleKeeper.GetConsensusNeeded(ctx, networkDescriptor)
	assert.NoError(t, err)
	diff := math.Abs(consensusNeeded - float64(0.1))
	assert.Equal(t, diff < 0.0001, true)
}
