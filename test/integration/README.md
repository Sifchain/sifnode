# Running the Integration Test Suite

## Tooling

The [setup-linux-environment.sh](./setup-linux-environment.sh) script will install all the tools you need to run in a fresh Linux environment (go, make, etc).  This is the script that github actions use to set up that environment.

## Execute in a local environment (not github actions)

Run `make` in test/integration/vagrant.  That uses [vagrant](https://www.vagrantup.com/docs/installation) to set up a fresh Linux environment with all the tools necessary for building and running the tests.  It will:

*  Mount local files to the virtualbox instance
*  Create a new Linux machine (using virtualbox).
*  Install the tools (using setup-linux-environment.sh)
*  Run the tests.
*  Leave a virtual machine running with the full test environment available for use.
*  Copies logs to `data/*` and tars them up into `datafiles.12-11-16-15-53.tar`.

Running `make` again will run the tests again in the existing environment.

## Docker

To build the docker container, run:

```
make sifdocker  # builds the docker image
make sifdocker-start  # starts the docker container and leaves it running
make sifdocker-sh     # gives you a shell in the running container
```
## Execute

[start-integration-env.sh](./start-integration-env.sh) starts 
sifnoded and ganache.

Run the tests in a container with:

```
cd /sifnode/test/integration
./start-integration-env.sh && . vagrantenv.sh
python3 -m pytest -v src/py/test_*
```

You can control the log level and which tests are run
with standard pytest options:

```
python3 -m pytest -olog_cli=true -olog_level=DEBUG -olog_file=/tmp/log.txt -v src/py/test_rowan_transfers.py::test_transfer_rowan_to_erowan
```
If you have a clean Ubuntu environment, these two commands will set up everything you need:

```
test/integration/setup-linux-environment.sh
test/integration/start-integration-env.sh
```

## Running tests against ropsten

###  Set the appropriate environment variables:

```
# only once, add a sifchain account that has rowan
# sifnodecli keys import rowansource key.txt --keyring-backend test

# Set an ethereum private key for an acount that has eth on the testnet
export SMART_CONTRACTS_DIR=...
export ETHEREUM_ADDRESS=...
export ROWAN_SOURCE=sif1pvnu2kh826vn8r0ttlgt82hsmfknvcnf7qmpvk
export ROWAN_SOURCE_KEY=thekey
export ETHEREUM_NETWORK=ropsten
export SIFNODE=http://xx.xx.xx.xx:26657
export INFURA_PROJECT_ID=yourId
cd $yourtestintegrationdirectory
source <(smart_contract_env.sh ~/workspace/sifnode/smart-contracts/deployments/sandpit/)
export ETHEREUM_PRIVATE_KEY=...

python3 -m pytest -x -olog_cli=true -olog_level=DEBUG -v -olog_file=/tmp/log.txt -v src/py/test_eth_transfers
```


1.  
## Github actions

See [the github action file](../../.github/workflows/integrationtest.yml) for the description of what's executed in the integration test environment.
