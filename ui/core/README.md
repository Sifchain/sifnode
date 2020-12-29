Principles

- Usecases define functionality
- Services can only talk to other services via `app` layer
- No frameworks
- Functions pass messages downstream
- Event handlers pass messages upstream

* `app` - Application Logic (Usecases)
  - `clp` - Continuous Liquidity Pool Logic
  - `peg` - Peggy Logic
  - `wallets` - Wallet Logic
* `services` - IO devices / detail
  - `ethereum` <-> Ethereum
  - `sifchain` <-> Sifchain
  - `view` <-> Frontend
  - `notifications` <-> Notification Dispatcher (Frontend/Device/Logging)
* `entities`
  - `formulae` - Domain logic for calculating values
  - objects - Domain Objects
