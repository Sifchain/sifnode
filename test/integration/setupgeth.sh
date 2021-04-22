initial_addresses="$*"

# httpport is the same port ganache uses
httpport=7545
wsport=8646

pkill geth || true
sleep 1

apis=personal,eth,net,web3,debug

nohup geth --networkid 5777 --datadir /tmp/gethdata \
  --dev \
  --ws --ws.addr 0.0.0.0 --ws.port $wsport --ws.api $apis \
  --http --http.addr 0.0.0.0 --http.port $httpport --http.api $apis \
  --dev.period 1 \
  --mine --miner.threads=1 > /tmp/gethlog.txt 2>&1 &

while ! nc -z localhost $wsport; do
  sleep 1
done

one_hundred_eth=100000000000000000000
for i in $initial_addresses
do
  geth attach /tmp/gethdata/geth.ipc --exec "eth.sendTransaction({from:eth.coinbase, to:\"$i\", value:$one_hundred_eth})"
  geth attach /tmp/gethdata/geth.ipc --exec "eth.getBalance(\"$i\")"
done

# tail -F /dev/null
