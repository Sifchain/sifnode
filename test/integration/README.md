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
pip3  install -r src/py/requirements.txt
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

### Start with exampleenv.sh

./exampleenv.sh is designed to be used in a shell like this:

```
# cd ~/workspace/sifnode/test/integration 
# source ./exampleenv.sh
# sifnoded q auth account --node tcp://44.241.55.154:26657 sif1pvnu2kh826vn8r0ttlgt82hsmfknvcnf7qmpvk
# ...
```

It should echo a set of sample commands and tests that you can run 
after it sets up your environment.

 Don't use exampleenv.sh itself, make a copy modify the variables to match your setup.

Don't mix this with vagrantenv.sh; use one or the other in a shell.

## Distributing ropsten test tokens

Set up your environment as if you were going to run tests, then run:
```
TOKEN_AMOUNT=120000000000 python3 -m pytest --color=yes -x -olog_cli=true -olog_level=DEBUG -v -olog_file=vagrant/data/pytest.log src/py/token_refresh.py
```

This creates a new sifchain account with 120000000000 tokens of each of the tokens in the whitelist.  

### To set up all test accounts with tokens:

Use the test/integration/distribute_tokens.sh script.

### Update the UI json files

When you have an environment set up as above, the build_ui_token_files.sh script will update the appropriate files in 
$BASEDIR/ui/core/src/assets.*.json.

where ROWAN_SOURCE is a sifnode address with sufficient tokens to distribute.

## Github actions

See [the github action file](../../.github/workflows/integrationtest.yml) for the description of what's executed in the integration test environment.
