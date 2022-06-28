# Integration Tests for Margin

In order to run the integration tests for Margin, you need to set the following environment variables:

```bash
export FEATURE_TOGGLE_SDK_045=1
export FEATURE_TOGGLE_MARGIN_CLI_ALPHA=1
export GOFLAGS="-modfile=go_045.mod"
export GOTAGS="FEATURE_TOGGLE_SDK_045,FEATURE_TOGGLE_MARGIN_CLI_ALPHA"
```

Those variables enable the relevant feature toggles required to run the integration tests for Margin.

Then you can use the following command to run the tests:

```bash
go test -tags $GOTAGS -v -failfast $(go list -tags $GOTAGS ./x/margin/client/testutil/... | grep -v /vendor/)
```
