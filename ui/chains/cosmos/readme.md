For reference when playing with cosmos cain

# Create liquidity pool catk:rwn

```
sifnodecli tx clp create-pool \
 --from akasha \
 --sourceChain ETH \
 --symbol ETH \
 --ticker catk \
 --nativeAmount 500 \
 --externalAmount 500
```

# Create liquidity pool cbtk:rwn

```
sifnodecli tx clp create-pool \
 --from akasha \
 --sourceChain ETH \
 --symbol ETH \
 --ticker cbtk \
 --nativeAmount 500 \
 --externalAmount 500
```

# Verify pool created

```
sifnodecli query clp pools
```

# Execute swap

```
sifnodecli tx clp swap \
 --from shadowfiend \
 --sentSourceChain ETH \
 --sentSymbol ETH \
 --sentTicker catk \
 --receivedSourceChain ETH \
 --receivedSymbol ETH \
 --receivedTicker cbtk \
 --sentAmount 20
```
