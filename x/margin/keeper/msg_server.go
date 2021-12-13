package keeper

import "github.com/Sifchain/sifnode/x/margin/types"

type msgServer struct {
	KeeperI
}

var _ types.MsgServer = msgServer{}

func NewMsgServerImpl(k KeeperI) types.MsgServer {
	return msgServer{
		k,
	}
}
