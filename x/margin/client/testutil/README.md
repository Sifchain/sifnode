# Integration Tests for Margin

In order to run the integration tests for Margin, you need to set the following feature toggles build-tag `FEATURE_TOGGLE_SDK_045` and `FEATURE_TOGGLE_MARGIN_CLI_ALPHA` and use the 045 go mod file as follows:

```bash
FEATURE_TOGGLE_SDK_045=1 \
FEATURE_TOGGLE_MARGIN_CLI_ALPHA=1 \
GOFLAGS="-modfile=go_045.mod" \
go test \
    -tags FEATURE_TOGGLE_SDK_045,FEATURE_TOGGLE_MARGIN_CLI_ALPHA \
    -v $(go list -tags FEATURE_TOGGLE_SDK_045,FEATURE_TOGGLE_MARGIN_CLI_ALPHA ./x/margin/client/testutil/... | grep -v /vendor/)
```
