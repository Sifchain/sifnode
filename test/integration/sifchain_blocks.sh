height=$(sifnoded --home $CHAINDIR/.sifnoded q block | jq -r .block.header.height)
seq $height | parallel -k sifnoded --home $CHAINDIR/.sifnoded q block {}
