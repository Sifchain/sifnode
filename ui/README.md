# Frontend repo

- `./core` - all business functionality
- `./app` - A Vue interface that uses Core
- `./chains` - Blockchain testing

## Run App and Core tests

`yarn app:serve` - Serve frontend app with the background blockchain
`yarn core:test` - Run tests on `core` module with the background blockchain

## Having more control

`yarn chain:start` - Start the background blockchain
`yarn chain:migrate` - Migrate the background blockchain must have the chain started
`yarn app:serve:nochain` - Serve frontend app with no background chain
`yarn core:test:nochain` - Run core tests with no background chain
