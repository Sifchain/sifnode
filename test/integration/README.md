# Running the Integration Test Suite

## Tooling

The [setup-linux-environment.sh](./setup-linux-environment.sh) script will install all the tools you need to run in a fresh Linux environment (go, make, etc).  This is the script that github actions use to set up that environment.

## Execute in test environment

Run `make` in test/integration/vagrant.  That uses [vagrant](https://www.vagrantup.com/docs/installation) to set up a fresh Linux environment with all the tools necessary for building and running the tests.  It will:

1.  Create a new Linux machine (using virtualbox).
2.  Install the tools (using setup-linux-environment.sh)
3.  Run the tests.
4.  Leave a virtual machine running with the full test environment available for use.

Running `make` again will run the tests again in the existing environment.

## Execute

[start-integration-env.sh](./start-integration-env.sh) starts docker instances and runs the integration tests.  The tests are scripts that exit with a non-zero status if they fail.

If you have a clean Ubuntu environment, these two commands will install everything and run the tests:

```
test/integration/setup-linux-environment.sh
test/integration/start-integration-env.sh
```

## Github actions

See [the github action file](../../.github/workflows/integrationtest.yml) for the description of what's executed in the integration test environment.