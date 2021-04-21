package keeper

import (
	"github.com/Sifchain/sifnode/x/dispensation/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the clp MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}