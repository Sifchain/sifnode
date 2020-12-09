For reference when playing with cosmos cain

# Create liquidity pool catk:rowan

```
sifnodecli tx clp create-pool \
 --from akasha \
 --symbol catk \
 --nativeAmount 500 \
 --externalAmount 500
```

# Create liquidity pool cbtk:rowan

```
sifnodecli tx clp create-pool \
 --from akasha \
 --symbol cbtk \
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
 --sentSymbol catk \
 --receivedSymbol cbtk \
 --sentAmount 20
```
