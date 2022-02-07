package events

import (
	"database/sql"

	"github.com/Sifchain/sifnode/tools/siflisten/chain"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
)

func Sync(clientCtx sdkclient.Context, _ *sql.DB) {
	chain.QueryMarginParams(clientCtx)
	chain.QueryPools(clientCtx)
	chain.BlockEvents(clientCtx)

}
