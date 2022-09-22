# Integration Tests for Margin

In order to run the integration tests for Margin, you need to set the following environment variables:

```bash
export GOTAGS="TEST_INTEGRATION"
```

Those variables enable integration tests when running tests for Margin.

Then you can use the following command to run the tests:

```bash
go test -tags $GOTAGS -v -failfast $(go list -tags $GOTAGS ./x/margin/client/testutil/... | grep -v /vendor/)
```
