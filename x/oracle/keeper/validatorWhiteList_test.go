package keeper_test

import (
	"testing"

	oracleKeeper "github.com/Sifchain/sifnode/x/oracle/keeper"
	"github.com/Sifchain/sifnode/x/oracle/types"
	"github.com/stretchr/testify/assert"
)

const networkID = 1

func TestKeeper_SetValidatorWhiteList(t *testing.T) {
	powers := []int64{3, 7}
	ctx, keeper, _, _, _, validateAddress, _ := oracleKeeper.CreateTestKeepers(t, 0.7, powers, "")
	networkDescriptor := types.NewNetworkDescriptor(networkID)
	whitelist := types.NewValidatorWhitelist()
	for index, address := range validateAddress {
		whitelist.UpdateValidator(address, uint32(powers[index]))
	}
	keeper.SetOracleWhiteList(ctx, networkDescriptor, whitelist)
	vList := keeper.GetOracleWhiteList(ctx, networkDescriptor)
	assert.Equal(t, len(vList.Whitelist), 2)
	assert.True(t, keeper.ExistsOracleWhiteList(ctx, networkDescriptor))
}

func TestKeeper_ValidateAddress(t *testing.T) {
	powers := []int64{3, 7}
	ctx, keeper, _, _, _, validateAddress, _ := oracleKeeper.CreateTestKeepers(t, 0.7, powers, "")
	networkDescriptor := types.NewNetworkDescriptor(networkID)
	whitelist := types.NewValidatorWhitelist()
	for index, address := range validateAddress {
		whitelist.UpdateValidator(address, uint32(powers[index]))
	}

	keeper.SetOracleWhiteList(ctx, networkDescriptor, whitelist)
	assert.True(t, keeper.ValidateAddress(ctx, networkDescriptor, validateAddress[0]))
	assert.True(t, keeper.ValidateAddress(ctx, networkDescriptor, validateAddress[1]))
	_, validateAddress = oracleKeeper.CreateTestAddrs(3)
	assert.False(t, keeper.ValidateAddress(ctx, networkDescriptor, validateAddress[2]))
}
