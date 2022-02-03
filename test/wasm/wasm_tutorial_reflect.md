# Sifchain - Wasm (Custom Messages)

Testng custom messages to talk to our own modules

## Setup

1. Initialize the local chain: `make init`

2. Start the chain: `make run`

## Store and Initialize

Store the contract:

```
sifnoded tx wasm store ./sc/reflect.wasm \
--from sif --keyring-backend test \
--gas 1000000000000000000 \
--broadcast-mode block \
--chain-id localnet  \
-y
```

Instantiate

```
sifnoded tx wasm instantiate 1 '{}' \
--amount 50000rowan \
--label "reflect" \
--from sif --keyring-backend test \
--gas 1000000000000000000 \
--broadcast-mode block \
--chain-id localnet \
-y
```

```
sifnoded query wasm contract sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6
```

```
sifnoded query wasm contract-state smart sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 '{"owner":{}}'
```

```
sifnoded tx wasm execute sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 \
'{"reflect_msg":{"msgs":[{"bank": {"send": {"to_address": "sif1l7hypmqk2yc334vc6vmdwzp5sdefygj2ad93p5", "amount": [{"denom": "rowan", "amount": "10"}]}}}]}}' \
--from sif --keyring-backend test \
--chain-id localnet \
--broadcast-mode block \
-y 
```

```
sifnoded query wasm contract-state smart sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 '{"capitalized":"hello world!"}'
```

```
sifnoded tx wasm execute sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 \
'{"reflect_msg":{"msgs":[{"bank":{"send":{"to_address":"sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd","amount":[{"denom":"rowan","amount":"15000"}]}}}]}}' \
  --from sif \
  --keyring-backend test \
  --chain-id localnet \
  --broadcast-mode block \
  -y
```

```
sifnoded tx wasm execute sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 \
'{"reflect_msg":{"msgs":[{"custom":{"debug":"this is the input message"}}]}}' \
  --from sif \
  --keyring-backend test \
  --chain-id localnet \
  --broadcast-mode block \
  -y
```

```
sifnoded tx wasm execute sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 \
'{"reflect_msg":{"msgs":[{"custom":{"raw": "12345"}}]}}' \
  --from sif \
  --keyring-backend test \
  --chain-id localnet \
  --broadcast-mode block \
  -y
```

```
sifnoded tx wasm execute sif14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s62cvu6 \
'{"reflect_msg":{"msgs":[{"custom":{"raw":"eyJAdHlwZSI6Ii9zaWZub2RlLmNscC52MS5Nc2dTd2FwIiwic2lnbmVyIjoiIiwic2VudF9hc3NldCI6eyJzeW1ib2wiOiJyb3dhbiJ9LCJyZWNlaXZlZF9hc3NldCI6eyJzeW1ib2wiOiJjZXRoIn0sInNlbnRfYW1vdW50IjoiMSIsIm1pbl9yZWNlaXZpbmdfYW1vdW50IjoiMTAifQ=="}}]}}' \
  --from sif \
  --keyring-backend test \
  --chain-id localnet \
  --broadcast-mode block \
  -y
```