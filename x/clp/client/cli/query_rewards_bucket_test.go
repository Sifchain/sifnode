package cli_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Sifchain/sifnode/testutil/network"
	"github.com/Sifchain/sifnode/testutil/nullify"
	"github.com/Sifchain/sifnode/x/clp/client/cli"
	"github.com/Sifchain/sifnode/x/clp/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func networkWithRewardsBucketObjects(t *testing.T, n int) (*network.Network, []types.RewardsBucket) {
	t.Helper()
	cfg := network.DefaultConfig()

	state := types.DefaultGenesisState()
	for i := 0; i < n; i++ {
		rewardsBucket := types.RewardsBucket{
			Denom:  strconv.Itoa(i),
			Amount: sdk.NewInt(int64(i)),
		}
		nullify.Fill(&rewardsBucket)
		state.RewardsBucketList = append(state.RewardsBucketList, rewardsBucket)
	}
	buf, err := cfg.Codec.MarshalJSON(state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.RewardsBucketList
}

func TestShowRewardsBucket(t *testing.T) {
	net, objs := networkWithRewardsBucketObjects(t, 2)

	ctx := net.Validators[0].ClientCtx
	common := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	tests := []struct {
		desc    string
		idDenom string

		args []string
		err  error
		obj  types.RewardsBucket
	}{
		{
			desc:    "found",
			idDenom: objs[0].Denom,

			args: common,
			obj:  objs[0],
		},
		{
			desc:    "not found",
			idDenom: strconv.Itoa(100000),

			args: common,
			err:  status.Error(codes.NotFound, "not found"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{
				tc.idDenom,
			}
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdShowRewardsBucket(), args)
			if tc.err != nil {
				stat, ok := status.FromError(tc.err)
				require.True(t, ok)
				require.ErrorIs(t, stat.Err(), tc.err)
			} else {
				require.NoError(t, err)
				var resp types.RewardsBucketRes
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.RewardsBucket)
				require.Equal(t,
					nullify.Fill(&tc.obj), //nolint
					nullify.Fill(&resp.RewardsBucket),
				)
			}
		})
	}
}

func TestListRewardsBucket(t *testing.T) {
	net, objs := networkWithRewardsBucketObjects(t, 5)

	ctx := net.Validators[0].ClientCtx
	request := func(next []byte, offset, limit uint64, total bool) []string {
		args := []string{
			fmt.Sprintf("--%s=json", tmcli.OutputFlag),
		}
		if next == nil {
			args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
		} else {
			args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
		}
		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
		if total {
			args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
		}
		return args
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			// print args
			fmt.Println(args)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdListRewardsBucket(), args)
			require.NoError(t, err)
			var resp types.AllRewardsBucketRes
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.RewardsBucket), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.RewardsBucket),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			// print args
			fmt.Println(args)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdListRewardsBucket(), args)
			require.NoError(t, err)
			var resp types.AllRewardsBucketRes
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.RewardsBucket), step)
			require.Subset(t,
				nullify.Fill(objs),
				nullify.Fill(resp.RewardsBucket),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.GetCmdListRewardsBucket(), args)
		require.NoError(t, err)
		var resp types.AllRewardsBucketRes
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(objs),
			nullify.Fill(resp.RewardsBucket),
		)
	})
}
