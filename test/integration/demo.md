# Peg 10 ETH

### Get token balance

yarn peggy:getTokenBalance 0xC5fdf4076b8F3A5357c5E395ab970B5B54098Fef

### Lock tokens

yarn peggy:lock \
 $(sifnodecli keys show juniper -a) \
 0x0000000000000000000000000000000000000000 \
 10000000000000000000

y advance 52

sifnodecli query account $(sifnodecli keys show juniper -a)

# Add liquidity to the pre-existing liquidity pool

sifnodecli tx clp add-liquidity \
 --from juniper \
 --symbol ceth \
 --nativeAmount 5000000000000000000 \
 --externalAmount 5000000000000000000

sifnodecli query account $(sifnodecli keys show juniper -a)

# Swap 10000 cUSDC for cETH

# Unpeg 10 cETH
