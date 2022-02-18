#!/bin/bash

go run cmd/dbtool/main.go --data ~/.sifnoded_prod/data/  \
--query "send_packet.packet_dst_port='transfer'" \
--out ~/send_packets.data \
--pages 2 \
--per-page 100 

go run cmd/dbtool/main.go --data ~/.sifnoded_prod/data/  \
--query "acknowledge_packet.packet_dst_port='transfer'" \
--out ~/acknowledge_packets.data \
--pages 2 \
--per-page 100 

go run cmd/dbtool/main.go --data ~/.sifnoded_prod/data/  \
--query "timeout_packet.packet_dst_port='transfer'" \
--out ~/timeout_packets.data \
--pages 2 \
--per-page 100 