package cmd

import (
	"fmt"
	"strconv"

	"github.com/Sifchain/sifnode/cmd/dbtool/utils"
	"github.com/spf13/cobra"
)

func NewGetStuckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-stuck",
		Short: "Get stuck transfers",
		Long: `Get stuck transfers 

ex: dbtool stuck conncetion-22`,
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

	sendPackets := getPackets("send_packet", connection)
	ackPackets := getPackets("acknowledge_packet", connection)
	timeoutPackets := getPackets("timeout_packet", connection)

	fmt.Printf("send_packet: %d\n", len(sendPackets))
	fmt.Printf("acknowledge_packet: %d\n", len(ackPackets))
	fmt.Printf("timeout_packet: %d\n", len(timeoutPackets))

	stuckPackets := []*utils.EventInfo{}

	for seq, pkts := range sendPackets {
		_, ack := ackPackets[seq]
		_, timeout := timeoutPackets[seq]
		if !(ack || timeout) {
			stuckPackets = append(stuckPackets, pkts...)
		}
	}

	fmt.Printf("stuck packets: %d\n", len(stuckPackets))
}

func getPackets(eventType string, connection string) map[uint64][]*utils.EventInfo {
	query := fmt.Sprintf("%s.packet_src_port='transfer' AND %s.packet_connection='%s'", eventType, eventType, connection)
	txs, err := utils.DoTxSearch(query, true, 1, 100)
	if err != nil {
		panic(err)
	}

	// outgoing transfers
	filter := func(string) bool {
		return eventType == "send_packet" ||
			eventType == "acknowledge_packet" ||
			eventType == "timeout_packet"

	}
	events := utils.FilterEvents(txs, filter)

	result := make(map[uint64][]*utils.EventInfo)
	for _, ev := range events {
		for _, attr := range ev.RealAttributes {
			if string(attr.Key) == "packet_sequence" {
				seq, err := strconv.ParseUint(string(attr.Value), 10, 32)
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
		}
	}

	return result
}
