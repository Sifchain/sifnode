package events

import (
	"database/sql"

	"github.com/Sifchain/sifnode/tools/siflisten/chain"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/pflag"
)

func Sync(clientCtx sdkclient.Context, _ *sql.DB, flagSet *pflag.FlagSet) {
	chain.QueryMarginParams(clientCtx)
	chain.QueryPools(clientCtx, flagSet)

}
