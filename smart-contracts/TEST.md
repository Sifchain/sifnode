To run unit tests that rely on devenv, start devenv in a terminal:

    npx hardhat run scripts/devenv.ts

And then run the tests in another terminal:

    npx hardhat test test/devenv/test_lockburn.ts --network localhost

The VERBOSE environment variable can be set to:

* summary - only print summary lines
* (not set) - no verbose output
* any string other than summary - full json output
