# Feature Toggles

This document keeps track of all existing feature toggles (also known as feature flags) that can be enabled using the relevant environment variables. Setting any of the following flag will enable the corresponding feature.

- `FEATURE_TOGGLE_SDK_045`: Switch to Cosmos SDK v0.45.1

## Example

To compile using Cosmos SDK v0.45 you can use the following command:

```bash
FEATURE_TOGGLE_SDK_045=1 make build-sifd
```

or

```bash
FEATURE_TOGGLE_SDK_045=1 make install
```
