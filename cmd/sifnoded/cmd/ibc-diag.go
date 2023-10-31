package cmd

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/query"
	chtypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	"github.com/spf13/cobra"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

var (
	srcNode      string
	dstNode      string
	srcChannel   string
	dstChannel   string
	transferPort = "transfer"
)

func NewIBCDiagCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ibc-diag",
		Short: "IBC diagnostics commands",
	}
	cmd.AddCommand(
		NewGetStuckTransfersCmd(),
	)
	return cmd
}

func NewGetStuckTransfersCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "stuck-txs",
		Short: "get stuck IBC transfers",
		Long: `get stuck IBC transfers

Example: Getting stuck transfers between sifchain and terra

sifnoded ibc-diag stuck-txs \
  --src-node http://rpc.sifchain.finance:80 \
  --dst-node http://public-node.terra.dev:26657 \
  --src-channel channel-18 \ 
  --dst-channel channel-7

Use the regular IBC commands to find the src and dst channels of a connection
`,
		Run: func(cmd *cobra.Command, args []string) {
			getStuckTransfers(cmd)
		},
	}
	command.Flags().StringVar(&srcNode, "src-node", srcNode, "rpc endpoint of source node")
	command.Flags().StringVar(&dstNode, "dst-node", dstNode, "rpc endpoint of destination node")
	command.Flags().StringVar(&srcChannel, "src-channel", dstChannel, "source channel id")
	command.Flags().StringVar(&dstChannel, "dst-channel", dstChannel, "destination channel id")
	_ = command.MarkFlagRequired("src-node")
	_ = command.MarkFlagRequired("dst-node")
	_ = command.MarkFlagRequired("src-channel")
	_ = command.MarkFlagRequired("dst-channel")
	return command
}

func getStuckTransfers(cmd *cobra.Command) {
	commitments, err := getCommittedPackets(cmd, srcNode, transferPort, srcChannel)
	if err != nil {
		panic(err)
	}

	unreceived, err := getUnreceivedPackets(
		cmd,
		dstNode,
		commitments,
		transferPort,
		dstChannel)
	if err != nil {
		panic(err)
	}

	transfers, err := getTransfers(srcNode, srcChannel, unreceived)
	if err != nil {
		panic(err)
	}

	// print results in csv format
	fmt.Println("packet_sequence, tx_hash, amount, denom, receiver, sender")
	for _, t := range transfers {
		fmt.Printf("%s, %s, %s, %s, %s, %s\n",
			t.PacketSequence,
			t.TxHash,
			t.PacketData.Amount,
			t.PacketData.Denom,
			t.PacketData.Receiver,
			t.PacketData.Sender,
		)
	}
}

// getCommmittedPackets returns the list of packets that were sent on a given
// channel/port but for which there is still a PacketCommitment in the
// underlying DB . A packet that still has a PacketCommitment is a packet whose
// receipt was never acknowledged and which hasn't yet timed out.
func getCommittedPackets(
	cmd *cobra.Command,
	nodeURI string,
	portID string,
	channelID string) ([]uint64, error) {

	clientCtx, err := getClientContext(cmd, nodeURI)
	if err != nil {
		return nil, err
	}

	queryClient := chtypes.NewQueryClient(clientCtx)

	packets := []uint64{}

	page := uint64(1)
	limit := uint64(100)

	for {
		pageReq := &query.PageRequest{
			Offset: (page - 1) * limit,
			Limit:  limit,
		}

		req := &chtypes.QueryPacketCommitmentsRequest{
			PortId:     portID,
			ChannelId:  channelID,
			Pagination: pageReq,
		}

		res, err := queryClient.PacketCommitments(context.Background(), req)
		if err != nil {
			return nil, err
		}

		packetSequences := make([]uint64, len(res.Commitments))
		for i, p := range res.Commitments {
			packetSequences[i] = p.Sequence
		}

		packets = append(packets, packetSequences...)

		if len(res.Commitments) < int(limit) {
			break
		} else {
			page++
		}
	}

	return packets, nil
}

// getUnreceivedPackets takes a list of packets and returns the subset that
// hasn't been received on the destination channel/port.
func getUnreceivedPackets(
	cmd *cobra.Command,
	nodeURI string,
	committedPackets []uint64,
	portID string,
	channelID string) ([]uint64, error) {

	clientCtx, err := getClientContext(cmd, nodeURI)
	if err != nil {
		panic(err)
	}

	queryClient := chtypes.NewQueryClient(clientCtx)

	req := &chtypes.QueryUnreceivedPacketsRequest{
		PortId:                    portID,
		ChannelId:                 channelID,
		PacketCommitmentSequences: committedPackets,
	}

	res, err := queryClient.UnreceivedPackets(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return res.Sequences, nil
}

// getTransfers takes a list of packet sequences and returns the corresponding
// list of Transfers which contain the actual transfer data.
func getTransfers(nodeURI string, channelID string, packets []uint64) ([]*Transfer, error) {
	c, err := rpchttp.New(nodeURI, "/websocket")
	if err != nil {
		return nil, err
	}

	transfers := []*Transfer{}
	for _, seq := range packets {
		ev, err := getTransfer(c, channelID, seq)
		if err != nil {
			return nil, err
		}
		transfers = append(transfers, ev)
	}

	return transfers, nil
}

// getTransfer fetches transfer data corresponding to a given packet.
func getTransfer(client *rpchttp.HTTP, channelID string, packetSequence uint64) (*Transfer, error) {
	query := fmt.Sprintf("send_packet.packet_sequence=%d AND send_packet.packet_src_channel='%s'", packetSequence, channelID)

	page := 1
	perPage := 100
	res, err := client.TxSearch(
		context.Background(),
		query,
		false,
		&page,
		&perPage,
		"asc",
	)
	if err != nil {
		return nil, err
	}

	filter := func(eventType string) bool {
		return eventType == "send_packet"
	}
	filteredEvents := FilterEvents(res.Txs, filter)

	if len(filteredEvents) == 0 {
		return nil, fmt.Errorf("pruned send_packet (sequence %d)", packetSequence)
	}
	if len(filteredEvents) > 1 {
		return nil, fmt.Errorf("multiple events (%d) for %s", len(filteredEvents), query)
	}

	ev := filteredEvents[0]
	var fungibleTokenPacket FungibleTokenPacketData
	err = json.Unmarshal(
		[]byte(ev.GetAttribute("packet_data")),
		&fungibleTokenPacket,
	)
	if err != nil {
		return nil, err
	}

	t := &Transfer{
		PacketSequence: ev.GetAttribute("packet_sequence"),
		TxHash:         ev.TxHash,
		PacketData:     fungibleTokenPacket,
	}

	return t, nil
}

func getClientContext(cmd *cobra.Command, nodeURI string) (*client.Context, error) {
	ctx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		return nil, err
	}

	ctx = ctx.WithNodeURI(nodeURI)

	srcClient, err := client.NewClientFromNode(nodeURI)
	if err != nil {
		return nil, err
	}

	ctx = ctx.WithClient(srcClient)

	return &ctx, nil
}

type Transfer struct {
	TxHash         string
	PacketSequence string
	PacketData     FungibleTokenPacketData
}

type FungibleTokenPacketData struct {
	Denom    string
	Amount   string
	Sender   string
	Receiver string
}

func FilterEvents(
	txs []*coretypes.ResultTx,
	typeFilter func(string) bool,
) []*EventInfo {
	infos := []*EventInfo{}
	for _, tx := range txs {
		txHash := hex.EncodeToString(tx.Tx.Hash())
		for _, ev := range tx.TxResult.Events {
			if typeFilter(ev.Type) {
				attributes := []string{}
				for _, attr := range ev.Attributes {
					attributes = append(attributes, attr.String())
				}
				info := &EventInfo{
					Type:           ev.Type,
					TxHash:         txHash,
					Attributes:     attributes,
					RealAttributes: ev.Attributes,
				}
				infos = append(infos, info)
			}
		}
	}
	return infos
}

type EventInfo struct {
	Type           string
	TxHash         string
	Attributes     []string
	RealAttributes []abcitypes.EventAttribute
}

func (ev *EventInfo) GetAttribute(key string) string {
	for _, attr := range ev.RealAttributes {
		if string(attr.Key) == key {
			return string(attr.Value)
		}
	}
	return ""
}
