# **COSMOS 0.42 UPDATE**

The Sifnode project is built with Cosmos SDK &amp; Tendermint engine. This document describes the most noticeable changes that have been made as part of the major update of the Sifnode codebase in order to to migrate from v0.39 to v0.42 of the Cosmos SDK. It aims to give context required to be able to interact with Sifnode after the upgrade, rather than doing a complete overview of modifications.

**Reference Material**

CosmosSDK Chain Upgrade Guide to .42 - [https://docs.cosmos.network/master/migrations/chain-upgrade-guide-040.html](https://docs.cosmos.network/master/migrations/chain-upgrade-guide-040.html)

**Document Outline**

Following is a summary of important areas that have been impacted by the upgrade and likely require review. This list is not meant to be complete:

- Changes to the REST/gRPC API interface and code that implements the interface

- Changes to the CLI commands, syntax, and code that implements the CLI
- Changes to the data model and overview of data types that have been modified or converted from Amino to Protobuf
- Changes to the state that impact genesis, state export, import, and migration
- Changes to integration, load, and regression tests, test coverage updates

**Changes to the REST/gRPC API interface and code that implements the interface**

Sifnode API interface has been changed to gRPC to be able to fully migrate to v0.42 of the Cosmos SDK. REST interface support has been preserved for backwards compatibility. Existing API endpoints have the same structure to make sure that no additional changes are needed on the client side. No new endpoints were added.

Key files related to API changes:

- **app.go** - all app module routes are registered in **RegisterAPIRoutes** (both legacy and grpc-gateway).

- **x/module_name/module.go** - now has **LegacyQuerierHandler** functionthat returns **NewQuerier** with gRPC &amp; legacy support
- **x/module_name/client/rest/rest.go|tx.go|query.go** - contain REST endpoints specific to a given module
- **x/module_name/keeper/grpc_query.go** - defines module-specific gRPC queries to query data from the corresponding Keeper.
- **x/module_name/keeper/querier.go** - defines module-specific queries that combine gPRC and legacy querier codec calls
- **go.mod** - has new dependencies for **grpc** and **grpc-gateway** support

**Changes to the CLI commands, syntax, and code that implements the CLI**

CLI commands and related code has undergone various changes, and improvements, including syntax modifications, fixes, and refactoring.

**sifnoded** and **sifnoded** have been combined into **sifnoded** , so all commands previously executed through sifnoded are now executed through **sifnoded**. Aside from the merge into sidnoded, command syntax remains largely unchanged. The changes to scripts in **scripts/demo** are minimal beyond the shift from sifnoded to sifnoded.

Syntax changes:

- **send** tx is now under the **bank** route rather than a root **tx** command

  - **Old:** sifnoded tx send [from\_key\_or\_address] [to\_address] [amount] [flags]
  - **New:** sifnoded tx bank send [from\_key\_or\_address] [to\_address] [amount] [flags]

Commands are no longer added in **cmd/sifnoded/main.go** , each route is now added in the following files:

- **cmd/sifnoded/cmd/root.go**
- **cmd/sifnoded/cmd/oracle.go**
- **cmd/sifnoded/cmd/migrate.go**
- **cmd/sifnoded/cmd/gentx.go**
- **cmd/sifnoded/cmd/genaccounts.go**
- **cmd/sifnoded/cmd/clpadmin.go**

Compatibility changes for each module&#39;s cli can be found in the modules&#39; **x/module_name/client/cli/** directories.

One quirk we ran into was that the **-o** output **flag** needed to be manually added to a module&#39;s transaction such as in **x/clp/client/cli/tx.go**

**Changes to the data model and overview of data types that have been converted**

Significant share of the changes required for the upgrade are related to the data model. To be able to support v0.42 of the Cosmos SDK, most of the existing types and structs were fully redesigned. Protobuf is now used for data serialization and encoding instead of Amino protocol.

As part of the process, data types definitions were removed from **.go** files and re-defined in **.proto** files with various changes forced by Protobuf specifics and limitations. These files are required to generate **.pb.go** files containing finalized type definitions.

Additionally, multiple predefined custom types that were added or imported from third party repositories like Cosmos (e.g. Coin type).

Key files related to data model changes:

- **proto/sifnode/module_name/v1/\*.proto** - define module-specific types as Protobuf messages/services.
- **x/module_name/types/\*.pb.go** - contain final type definitions, generated from messages and services defined in **\*.proto** files
- **x/module_name/types/\*.go** - files located in **../types** directory can still contain some factory type definitions and related utility functions
- **x/module_name/keeper/\*.go** - aforementioned types are being used in multiple files from the **../keeper** directory of each module.

After switching to Protobuf, many custom types were completely removed and replaced with **strings** (where possible). For example:

- **proto/sifnode/clp/v1/types.proto** - **LiquidityProvider.liquidity_provider_address** type has been changed from **sdk.AccAddress** to **string**.
- **proto/sifnode/clp/v1/querier.proto** - **PoolRes.clp_module_address** type has been changed from **ClpModuleAddress** to **string**
- **proto/sifnode/clp/v1/querier.proto** - **LiquidityProviderRes.native_asset_balance** type has been changed from **NativeAssetBalance** to **string**
- **proto/sifnode/clp/v1/querier.proto** - **LiquidityProviderRes.external_asset_balance** type has been changed from **ExternalAssetBalance** to **string**.

Custom types, imported from third-party repos like Cosmos and Google:

- **third_party/proto/cosmos/base/coin.proto**
- **third_party/proto/cosmos/base/query/v1beta1/pagination.proto**
- **third_party/proto/google/api/http.proto**
- **third_party/proto/google/api/annotations.proto**
- **third_party/proto/google/api/httpbody.proto**
- **third_party/proto/gogoproto/gogo.proto**

**Changes to the state that impact genesis, state export, import, and migration**

Genesis and the state migration process were mainly impacted by changes required to migrate to different protocols for data encoding and serialization, as well as various modifications to functionality and architecture forced by the 0.42 upgrade.

Missing genesis export functions were added to modules as they play a crucial role in restarting the chain when upgraded to .42.

A custom migrate-data command was introduced which migrates the genesis structure of our custom modules from the 39 branch to that on the 42 branch.

Slices should not be serialised directly to storage, so a proto message wrapper should be used where slices were previously serialised.

**Changes to integration, load, and regression tests, test coverage updates**

Removed integration tests:

- **get_faucet_balance**
- **get_currency_faucet_balance**
- **transfer_new_currency**
