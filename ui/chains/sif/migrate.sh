#!/bin/bash
. ../credentials.sh

# if we don't sleep there are issues
sleep 10

# create liquidity pool from catk:rowan
echo "create liquidity pool from catk:rowan"


# nativeAmount 10000000 catk
# externalAmount 10000000 rowan
sifnodecli tx clp create-pool \
 --from akasha \
 --symbol catk \
 --nativeAmount   10000000000000000000000000 \
 --externalAmount 10000000000000000000000000 \
 --yes

# if we don't sleep there are issues
sleep 5

echo "create liquidity pool from cbtk:rowan"
# create liquidity pool from cbtk:rowan
# nativeAmount 10000000 cbtk
# externalAmount 10000000 rowan
sifnodecli tx clp create-pool \
 --from akasha \
 --symbol cbtk \
 --nativeAmount   10000000000000000000000000 \
 --externalAmount 10000000000000000000000000 \
 --yes

# should now be able to swap from catk:cbtk

sleep 5

echo "create liquidity pool from ceth:rowan"
# nativeAmount 8300 ceth
# externalAmount 10000000 rowan
sifnodecli tx clp create-pool \
 --from akasha \
 --symbol ceth \
 --nativeAmount   10000000000000000000000000 \
 --externalAmount 8300000000000000000000 \
 --yes

 # should now be able to swap from x:ceth

sleep 5

echo "create liquidity pool from cusdc:rowan"
sifnodecli tx clp create-pool \
 --from akasha \
 --symbol cusdc \
 --nativeAmount   10000000000000000000000000 \
 --externalAmount 10000000000000000000000000 \
 --yes

sleep 5

echo "create liquidity pool from clink:rowan"
sifnodecli tx clp create-pool \
 --from akasha \
 --symbol clink \
 --nativeAmount   10000000000000000000000000 \
 --externalAmount 588235000000000000000000 \
 --yes

