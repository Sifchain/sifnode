# Frontend repo

🚧 This is currently under construction and may not work. 🚧

## Installation

### Prerequisites

- [Go](https://golang.org/doc/install) (to build the sifchain to test against)
- [Node 14](https://nodejs.org/en/)
- [Yarn Classic](https://classic.yarnpkg.com/en/docs/install#mac-stable)
- A linux like environment

### Setup

1. Install the base sifnode repo in your go directory: `~/go/src/github.com/Sifchain/sifnode`
1. `cd ./ui` - To work on the frontend UI
1. `yarn` - Install dependencies
1. `yarn build` - Build the blockchains to test against

NOTE: If you are using VSCode you should use the code-workspace at `./ui/SifnodeUI.code-workspace` to ensure that Vetur works correctly.

### Launching locally

There are a few ways you can launch the project stack locally. Most of the time working on frontend you will probably just want to use:

```
yarn stack
```

NOTE: This command requires [tmux](https://github.com/tmux/tmux/wiki/Installing)

### How do I change something about the backing stack?

Let's say you want to do one of the following type of things:

- Add a new account
- Provide some genesis tokens
- Add a new token for localnet
- Change anything about our blockchain setup
- Respond to an environment request from one of the blockchain teams

For any of these things you will want ot create a new snapshot. Backend state is saved to a snapshot that is shared with the team or quick start development and also affects our e2e tests. You can create a new snapshot if you need by using the scripts below. Some of these commands are a little confusing so this table shows you when to use which.

| command                            | I want to                                                                     | quick start | ongoing | sif | eth | ebrelayer | FE  | setup scripts | save snapshot |
| ---------------------------------- | ----------------------------------------------------------------------------- | ----------- | ------- | --- | --- | --------- | --- | ------------- | ------------- |
| `yarn stack`                       | Work on frontend (requires tmux)                                              | ✅          | ✅      | ✅  | ✅  | ✅        | ✅  | 🚫            | 🚫            |
| `yarn stack:backend`               | Run only backing services from a snapshot say during CI.                      | ✅          | ✅      | ✅  | ✅  | ✅        | 🚫  | 🚫            | 🚫            |
| `yarn stack:backend-from-scripts`  | Run backing with setup scripts to manually change state and create a snapshot | 🚫          | ✅      | ✅  | ✅  | ✅        | 🚫  | ✅            | 🚫            |
| `yarn stack:save-default-snapshot` | Save new setup scripts to a snapshot                                          | 🚫          | 🚫      | ✅  | ✅  | ✅        | 🚫  | ✅            | ✅            |
| `yarn stack:save-snapshot`         | Save a snapshot from whatever is running                                      | 🚫          | 🚫      | 🚫  | 🚫  | 🚫        | 🚫  | 🚫            | ✅            |

You can either

- Alter the setup scripts that configure the blockchain and save the result as a snapshot (`yarn stack:save-snapshot-from-scripts`)
- Run the setup scripts do some kind of account action manually and then create a snapshot. (`yarn stack:backend-from-scripts` -> make a transaction -> `yarn stack:save-snapshot`)

It is preferrable that you include any changes you make in the setup scripts (say in `chains/post_migrate.sh`) as it means that it is possible your setup might get overwritten at a later date.

### Run tests in core

`yarn test`

### Run App and Core tests

| Command                | Description                                                |
| ---------------------- | ---------------------------------------------------------- |
| `yarn test`            | Alias for `core:test:all`                                  |
| `yarn build`           | Build core, all chains and the frontend app                |
| `yarn app:serve:all`   | Serve frontend app with the background blockchains         |
| `yarn core:test:all`   | Run tests on `core` module with the background blockchains |
| `yarn chain:start:all` | Start the background blockchains                           |

## Having more control

| Command              | Description                                                   |
| -------------------- | ------------------------------------------------------------- |
| `yarn chain:eth`     | Start the background ethereum blockchain                      |
| `yarn chain:sif`     | Start the background sifnode blockchain                       |
| `yarn chain:migrate` | Migrate the background blockchain must have the chain started |
| `yarn app:serve`     | Serve frontend app with no background chain                   |
| `yarn core:test`     | Run core tests with no background chain                       |
| `yarn core:watch`    | Compile core code in watch mode                               |

## Folder structure

| Path               | Description                      |
| ------------------ | -------------------------------- |
| `./app`            | A Vue interface that uses core.  |
| `./chains`         | Blockchain projects for testing. |
| `./core`           | All business functionality.      |
| `./docs`           | Documentation.                   |
| `./docs/decisions` | Architectural decisions.         |

## Architecture

We are following architecture influenced by clean architecture.

https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

<img src="./docs/FEArchitecture.png" />

_Example action/service dependencies_

The main premise here is that we have a domain consisting of **actions** and **entities** which communicate with the outside world over `api` and `store` channels.

| Section         | Description                                                                                                   |
| --------------- | ------------------------------------------------------------------------------------------------------------- |
| `core/actions`  | Actions, (aka usecases, interactors, commands) These hold the business logic and policies for the application |
| `core/entities` | Application data types.                                                                                       |
| `core/api`      | Input/output services. This is where you write and read data to wallets remote endpoints, rpc etc.            |
| `core/store`    | Shared reactive state between the `actions` and the view                                                      |
| `app`           | View application that renders UI                                                                              |

Every part of this system is designed to facilitate easy testing.

## Testing

### Testing Actions

Actions can be grouped arbitrarily by domain aggregate and may have their dependencies injected using the supplied creator. You ask for your api and store keys by using the given TS types.

```ts
// Generic params specify what API the service expects
type ActionContext<ServiceKeys, StoreKeys>
```

```ts
export default function createAction({
  api,
  store,
}: ActionContext<"WalletService" | "SifService", "WalletStore">) {
  return {
    async disconnectWallet() {
      await api.WalletService.disconnect();
      store.WalletStore.isConnected = false;
      store.WalletStore.balances = [];
      await api.SifService.disconnect();
    },
  };
}
```

The reason we do it this way is that in testing we only need to give the action creator exactly what it needs.

```ts
const actions = createAction({ api: { WalletService: fakeWalletService } });

// Then under test the wallet service runs with it's dependencies
actions.disconnectWallet();
```

### Testing Blockchain Driven Api

In the same way that Actions have their dependencies injected we can inject dependencies to our services layer.

```ts
export default function createFooService(context: FooServiceContext) {
  return {
    async doStuff() {
      const provider = await context.getWeb3Provider();
      // ...
    },
  };
}
```

### Etherium based blockchain development.

To test our blockchain backed apps we use ganache-cli and truffle to create a local etherium chain that is migrated to contain a couple of fake tokens.

You can find the token contracts [here](../chains/eth/contracts).

Our API setup asks for getters to supply environment information. It may make sense to convert this to a function that returns a config object we inject everywhere.

To test manually run the app using serve which includes ganache running in the background

```bash
./ui> yarn app:serve:all
```

> Alternatively you can run the following processes in separate terminals:
>
> 1. `yarn chain:eth` - Will run ganache
> 1. `yarn chain:sif` - Will run a built sifnode
> 1. `yarn chain:migrate` - Will run migrations against the running blockchains then exit
> 1. `yarn app:serve` - Will run the Vue app

From the terminal window running ganache make note of the first private key that gets generated:

<img src="./docs/ganache-keys.png" />

Then fire up the app on http://localhost:8080/.

Go to metamask. Click on the right corner menu and select "Import Account"

<img src="./docs/metamask1.png" width="300" />

Paste your private key there and you will load up your account on metamask.

<img src="./docs/metamask2.png" width="300" />

Hit import and select this account. Be sure to have this account selected. Reload the page click the connect wallet button and run through the procedure to connect your wallet in metamask.

You should see the balances of your wallet in the application.

### Component Library

To view and develop components in isolation you can use Storybook.

```
yarn storybook
```

### Testing stores

Stores are created using factory functions so that their state can be set upon creation. The store is the state our view responds to. It makes sense to test the actions and resultant store effects together as a usecase as we require no further dependencies. We can supply stores to actions in a similar way to the way we supply apis.

### Testing Views

Testing views is not as important as testing core code so we can defer to e2e testing for that. That is why it is ok to share configured action and store instances with Vue components. However complex render functionality or computed properties should be contained within stores or render actions.
