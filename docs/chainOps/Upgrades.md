# Sifnode Upgrades

## Packaging and proposing an upgrade

Upgrades on the Sifnodes themselves are automated, however there are several actions required to be performed prior.

1. Generate a new release [here](https://github.com/Sifchain/sifnode/releases).

2. Github actions will then build a new sifnoded binary and add it to the assets of the release. A sha256 will also be published (you will need this also).

3. Submit an upgrade proposal to the network:

```
sifnodecli tx gov submit-proposal software-upgrade <upgrade_name> \
    --from <from> \
    --deposit <deposit> \
    --upgrade-height <height> \
    --info '{"binaries":{"linux/amd64":"<url>"}}' \
    --title <title> \
    --description <description>
```

Where:

| Parameter | Description |
|-----------|-------------|
| `<upgrade_name>` | Name of the upgrade. This must match the upgrade name being used by the upgrade handler in the new binary. |
| `<from>` | The moniker of the validator proposing the upgrade. |
| `<deposit>` | The deposit/fee for proposing the upgrade (this is configurable in genesis). |
| `<height>` | The block height at which the upgrade should take place (must be greater than the voting period). |
| `<url>` | The URL to the new binary, including the SHA256 checksum as a query parameter. |
| `<title>` | The title of the upgrade. |
| `<description>` | A short description of the upgrade. |

e.g.:

```
sifnodecli tx gov submit-proposal software-upgrade sifnoded \
    --from my-node-moniker \
    --deposit 10000000000rowan \
    --upgrade-height 123456789 \
    --info '{"binaries":{"linux/amd64":"https://example.com/sifnode.zip?checksum=sha256:8630d1e36017ca680d572926d6a4fc7fe9a24901c52f48c70523b7d44ad0cfb2"}}' \
    --title 'Brave new world' \
    --description 'A special new upgrade'
```

## Voting on an upgrade proposal

To vote on a proposal, simply run:

```
sifnodecli tx gov vote <proposal_id> yes \
    --from <from> \
    --keyring-backend file \
    --chain-id <chain_id>
```

| Parameter | Description |
|-----------|-------------|
| `<proposal_id>` | The proposal ID. |
| `<from>` | The moniker of the validator voting on the upgrade. |
| `<chain_id>` | The chain ID of the network. |

e.g.:
 
```
sifnodecli tx gov vote 1 yes \
    --from my-node-moniker \
    --keyring-backend file \
    --chain-id sifchain
```
