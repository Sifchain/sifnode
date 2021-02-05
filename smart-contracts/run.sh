sifaddress=( 'sif1fpq67nw66thzmf2a5ng64cd8p8nxa5vl9d3cm4' )
contracts=( '0xfA8fC9C22C33FE62BabD5D92DD38Aa27B730d562' '0x64C4f2e91288FFA7661514eb9e1eb55D3b4013dC' '0x7cd621B1B092628F67591005AbDa07E27AC036d4' '0x9e8bD20374898F61B4e5BD32b880b7fE198E44a1' )

for address in "${sifaddress[@]}"
do
    # for contract in "${contracts[@]}"
    # do
    # for i in {1..20}
        echo "sifaddress: ${address} token address: ${contract}"
        BRIDGEBANK_ADDRESS=0x90DefDBd69a5cE54c3ECF1a9069856ccfdd739a6 npx truffle exec scripts/sendLockTx.js --network ropsten $address 0xfA8fC9C22C33FE62BabD5D92DD38Aa27B730d562 100000000000000000000000000000 >> run.txt &
        # sleep 1
    # done
done
