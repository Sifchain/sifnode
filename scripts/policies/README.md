# monetary policy manual testing

contains a set of JS scripts to prepare test environment for monetary policy testing.

# troubleshooting

in case `getAccountNumber` does not work, it is because the account info we are trying to get from a sif address does not exist yet in the chain, a workaround is to first send a transaction to this address from another one.
