package keeper_test

import (
	"testing"

	"github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/stretchr/testify/assert"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sifapp "github.com/Sifchain/sifnode/app"
)

func TestKeeper_SetNativeToken(t *testing.T) {
	app := sifapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	networkDescriptor := types.NewNetworkIdentity(types.NetworkDescriptor_NETWORK_DESCRIPTOR_ETHEREUM)
	token := "ceth"

	app.OracleKeeper.SetNativeToken(ctx, networkDescriptor, token)

	tokenStored, err := app.OracleKeeper.GetNativeToken(ctx, networkDescriptor)
	assert.NoError(t, err)
	assert.Equal(t, token, tokenStored)
}
