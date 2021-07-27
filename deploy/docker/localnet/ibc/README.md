# IBC LocalNet

This will launch the following:

* Two sifnode instances running independent chains.
* An [IBC relayer](https://github.com/confio/ts-relayer) to connect both chains.

## Setup/Build

1. Switch to the `deploy/docker/localnet/ibc` directory.

2. Build a new `sifnode` image:

```bash
SERVICE=sifnode make build-image
```

3. Build a new `ts-relayer` image:

```bash
SERVICE=ts-relayer make build-image
```

## Run

1. Switch to the `deploy/docker/localnet/ibc` directory.

2. Launch `docker-compose` as follows:

```bash
CHAINNET0=${CHAINNET0} \
CHAINNET1=${CHAINNET1} \
IPADDR0=${IPADDR0} \
IPADDR1=${IPADDR1} \
IPADDR2=${IPADDR2} \
SUBNET=${SUBNET} \
MNEMONIC='${MNEMONIC}' docker-compose up
```

Where:

|Var|Description|
|---|-----------|
|`${CHAINNET0}`|The Chain ID of the first sifnode (e.g.: `sifchain-ibc-0`)|
|`${CHAINNET1}`|The Chain ID of the second sifnode (e.g.: `sifchain-ibc-1`)|
|`${IPADDR0}`|The IP Address of the first sifnode (e.g.: `192.168.65.2`)|
|`${IPADDR1}`|The IP Address of the second sifnode (e.g.: `192.168.65.3`)|
|`${IPADDR2}`|The IP Address of the relayer (e.g.: `192.168.65.3`)|
|`${SUBNET}`|The subnet of the bridged network that Docker needs to create (e.g.: `192.168.65.1/24`)| 
|`${MNEMONIC}`|The mnemonic both sifnode's will use for their genesis accounts.|

e.g.:

```bash
CHAINNET0=sifchain-ibc-0 \
CHAINNET1=sifchain-ibc-1 \
IPADDR0=192.168.65.2 \
IPADDR1=192.168.65.3 \
IPADDR2=192.168.65.4 \
SUBNET=192.168.65.1/24 \
MNEMONIC='toddler spike waste purpose neutral beach science dawn joke stock help beyond' docker-compose up
```

## Notes

Currently, the relayer generates its own mnemonic, and the resulting address needs to be funded on both chains. The script `ts-relayer.sh` will perform this automatically when the container boots. This TypeScript implementation was used simply as a proof of concept, given the [issues experienced](https://discord.com/channels/669268347736686612/773388941947568148/839049449551691797) when attempting to use Cosmos' own IBC relayer.
