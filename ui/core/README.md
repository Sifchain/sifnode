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
  - `clp` <-> go `/x/clp` module
  - `ethbridge` <-> go `/x/ethbridge` module
  - `ethereum` <-> General ethereum blockchain
  - `sifchain` <-> General Sifchain blockchain
  - `view` <-> Frontend app
  - `notifications` <-> Notification Dispatcher (Frontend/Device/Logging)
* `entities`
  - `formulae` - Domain logic for calculating values
  - objects - Domain Objects
