height=$(sifnodecli --home $CHAINDIR/.sifnodecli q block | jq -r .block.header.height)
seq $height | parallel -k sifnodecli --home $CHAINDIR/.sifnodecli q block {}
