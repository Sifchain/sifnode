package cmd

import (
	"fmt"
	"strconv"

	"github.com/Sifchain/sifnode/cmd/dbtool/utils"
	chtypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"
	"github.com/spf13/cobra"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func NewStuckTxsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stuck-txs",
		Short: "Get stuck IBC transfers by connection ID",
		Long: `Get stuck IBC transfers by connection ID 

This command finds all IBC transfers that were sent from this chain on a given
connection and for which there is not acknowlegement or timeout.

ATTENTION: 
It is important to run this command against a FULL ARCHIVE database, otherwise
it will output false positives. Do not run this command against a snapshot.

ex: dbtool stuck-txs connection-22`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			getStuckTxs(args[0])
		},
	}
	return cmd
}

func getStuckTxs(connection string) {
	err := utils.OpenDB(datadir)
	if err != nil {
		panic(err)
	}

	channels := getChannels(connection)

	channel := channels[0]

	fmt.Printf("channel [%s] => [%s]\n", channel.ChannelId, channel.Counterparty.ChannelId)

	sendPackets := getPackets("send_packet", channel.ChannelId, channel.Counterparty.ChannelId)
	ackPackets := getPackets("acknowledge_packet", channel.ChannelId, channel.Counterparty.ChannelId)
	timeoutPackets := getPackets("timeout_packet", channel.ChannelId, channel.Counterparty.ChannelId)

	fmt.Printf("send_packets: %d\n", len(sendPackets))
	fmt.Printf("acknowledge_packets: %d\n", len(ackPackets))
	fmt.Printf("timeout_packets: %d\n", len(timeoutPackets))

	stuckPackets := []*utils.EventInfo{}

	for seq, pkts := range sendPackets {
		_, isAck := ackPackets[seq]
		_, isTimeout := timeoutPackets[seq]
		if !(isAck || isTimeout) {
			stuckPackets = append(stuckPackets, pkts...)
		}
	}

	fmt.Printf("stuck packets: %d\n", len(stuckPackets))

	for _, info := range stuckPackets {
		fmt.Printf("tx: %s\n", info.TxHash)
		fmt.Printf("type: %s\n", info.Type)
		fmt.Printf("attributes:\n")
		for _, attr := range info.Attributes {
			fmt.Printf("	%s\n", attr)
		}
		fmt.Printf("\n")
	}
}

func getChannels(connectionID string) []chtypes.IdentifiedChannel {
	sifApp, err := utils.NewSifApp(datadir)
	if err != nil {
		panic(err)
	}

	lastBlockHeight := sifApp.LastBlockHeight()

	ctx := sifApp.NewContext(
		true,
		tmproto.Header{Height: lastBlockHeight},
	)

	allChannels := sifApp.IBCKeeper.ChannelKeeper.GetAllChannels(ctx)

	connectionChannels := []chtypes.IdentifiedChannel{}
	for _, ch := range allChannels {
		for _, hop := range ch.ConnectionHops {
			if hop == connectionID {
				connectionChannels = append(connectionChannels, ch)
			}
		}
	}

	return connectionChannels
}

func getPackets(eventType string, srcChannelID string, dstChannelID string) map[uint64][]*utils.EventInfo {
	query := fmt.Sprintf("%s.packet_src_channel='%s' AND %s.packet_dst_channel='%s'", eventType, srcChannelID, eventType, dstChannelID)
	txs, err := utils.DoTxSearch(query, true, 1, 100)
	if err != nil {
		panic(err)
	}

	// outgoing transfers
	filter := func(evType string) bool {
		return evType == "send_packet" ||
			evType == "acknowledge_packet" ||
			evType == "timeout_packet"

	}
	events := utils.FilterEvents(txs, filter)

	result := make(map[uint64][]*utils.EventInfo)
	for _, ev := range events {
		seqS := ev.GetAttribute("packet_sequence")
		seq, err := strconv.ParseUint(seqS, 10, 32)
		if err != nil {
			panic(err)
		}
		seqEvs, ok := result[seq]
		if !ok {
			seqEvs = []*utils.EventInfo{}
		}
		seqEvs = append(seqEvs, ev)
		result[seq] = seqEvs
	}

	return result
}
