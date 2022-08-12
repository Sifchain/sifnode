# pool.proto

The Feature Toggle `FEATURE_TOGGLE_MARGIN_CLI_ALPHA` as described in `docs/general/FeatureToggles.md`, introduces additional fields to the `Pool` message that should not be exposed to the mainline.

To avoid having those fields exposed in the mainline code, we generated two versions of the `pool.pb.go` file:

- `pool_FEATURE_TOGGLE_MARGIN_CLI_ALPHA.pb.go`: generated protobuf code that contains the `FEATURE_TOGGLE_MARGIN_CLI_ALPHA` additional fields of the Pool message
- `pool_NO_FEATURE_TOGGLE_MARGIN_CLI_ALPHA.pb.go`: generated protobuf code that contains the original Pool message without the additional fields from `FEATURE_TOGGLE_MARGIN_CLI_ALPHA`

Whenever the protobuf files need to be re-generated, we should make sure to include the relevant changes of the Pool message generated codes to one of those files.

Also please make sure that the `pool.pb.go` file is not added and commited to the repository.
