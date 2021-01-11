. ../credentials.sh

# if we don't sleep there are issues
sleep 10

# create liquidity pool from catk:rowan
echo "create liquidity pool from catk:rowan"
sifnodecli tx clp create-pool \
 --from akasha \
 --symbol catk \
 --nativeAmount 1000000 \
 --externalAmount 1000000 \
 --yes

# if we don't sleep there are issues
sleep 5

echo "create liquidity pool from cbtk:rowan"
# create liquidity pool from cbtk:rowan
sifnodecli tx clp create-pool \
 --from akasha \
 --symbol cbtk \
 --nativeAmount 1000000 \
 --externalAmount 1000000 \
 --yes

# should now be able to swap from catk:cbtk