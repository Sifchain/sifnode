# if we don't sleep there are issues
sleep 10

# create liquidity pool from catk:rwn
echo "create liquidity pool from catk:rwn"
sifnodecli tx clp create-pool \
 --from akasha \
 --sourceChain ETH \
 --symbol ETH \
 --ticker catk \
 --nativeAmount 1000000 \
 --externalAmount 1000000 \
 --yes

# if we don't sleep there are issues
sleep 5

echo "create liquidity pool from cbtk:rwn"
# create liquidity pool from cbtk:rwn
sifnodecli tx clp create-pool \
 --from akasha \
 --sourceChain ETH \
 --symbol ETH \
 --ticker cbtk \
 --nativeAmount 1000000 \
 --externalAmount 1000000 \
 --yes

# should now be able to swap from catk:cbtk