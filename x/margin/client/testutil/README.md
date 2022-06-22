# Integration Tests for Margin

In order to run the integration tests for Margin, you need to set the build-tag feature toggle `FEATURE_TOGGLE_SDK_045` and use the 045 mod file as follows:

```bash
FEATURE_TOGGLE_SDK_045=1 \
GOFLAGS="-modfile=go_045.mod" \
go test \
    -tags FEATURE_TOGGLE_SDK_045 \
    -v $(go list -tags FEATURE_TOGGLE_SDK_045 ./x/margin/client/testutil/... | grep -v /vendor/)
```
