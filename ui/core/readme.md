```
#!/bin/bash
# to compile into ./lib then watch for changes
# this enables parallel repo work,
# with `../app` using this as a dependency

tsc -w
```

# Prerequisites

1. Install ganache [ganache-cli](https://github.com/trufflesuite/ganache-cli) (Globally currently)

2. Install [truffle](https://www.trufflesuite.com/docs/truffle/getting-started/installation)

3. Install dependencies to the truffle project under [fixtures/ethereum](fixtures/ethereum) by running `yarn` in that folder.
