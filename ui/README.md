# Frontend repo

- `./core` - All business functionality
- `./app` - A Vue interface that uses Core
- `./chains` - Blockchain testing

## Run App and Core tests

`yarn app:serve:all` - Serve frontend app with the background blockchains
`yarn core:test:all` - Run tests on `core` module with the background blockchains
`yarn chain:start:all` - Start the background blockchains

## Having more control

`yarn chain:eth` - Start the background ethereum blockchain
`yarn chain:sif` - Start the background sifnode blockchain
`yarn chain:migrate` - Migrate the background blockchain must have the chain started
`yarn app:serve` - Serve frontend app with no background chain
`yarn core:test` - Run core tests with no background chain
`yarn core:watch` - Compile core code in watch mode
