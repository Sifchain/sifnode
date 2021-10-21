# Scalability + localnet test framework

The current directly hosts two important testing components known as:

- Scalability test framework: run load tests against IBC and peggy.
- Localnet framework: initite and run IBC chains and relayers locally and run tests against a local IBC network of chains.

# Localnet

## Context

The current test environment Sifchain relies on requires several minutes to hours to run tests and are not conveninant for fast iteration development cycles as the developers have to wait a long time before getting any meaningful test result. It also discourages anyone involved in test engineering to write further tests in such environment as the short-term benefits of testing manually but poorly outpaces writing test code that takes long to provide results but do not improve testing and QA processes overall.

The localnet test environment solves those issues by running a local stack of services that combines IBC chains and IBC relayers all hosted within the same local machine and network layer. Therefore removing any costly latency due to containerization or network reliability.

## Getting started

In order to use the localnet test environment, use the following command lines from the current directory.

### Install dependencies

First install the project dependencies:

```
yarn
```

### Initiate all IBC chains

Then initiate all the IBC chains that are defined in the [config/chains.json](./config/chains.json) file. Feel free to add more IBC chains to it along the required information.

```
yarn initAllChains
```

All the IBC chains have been created along with their genesis accounts and validators nodes.

### Start all the IBC chains

Start now all the IBC chains that have been previously created using the following command line:

```
yarn startAllChains
```

To stop all the running processes, use `CTRL+C` or any other combinations to stop a running command.

### Initiate all the IBC relayers

Let's create now all the IBC relayers connection point between the chains using the following command line:

```
yarn initAllRelayers
```

All the IBC relayers files have been now created along with their connections and channels information.

### Start all the IBC relayers

Now start the relayers with this command line:

```
yarn startAllRelayers
```

To stop all the running processes, use `CTRL+C` or any other combinations to stop a running command.

### Testing

Run the following test command that transfers a certain amonut from IBC chain to another:

```
yarn zx test/test-ibc-transfer-between-two-chains.mjs
```

### Snapshot

TBD
