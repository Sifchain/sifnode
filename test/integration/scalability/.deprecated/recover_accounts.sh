#!/bin/bash

NAMES=(
    'sifchain-devnet-source'
    'sifchain-testnet-source'
    'devnet-ibctransfer-loadtest-siftochain-2021-08-18'
    'devnet-ibctransfer-loadtest-chaintosif-2021-08-18'
    'devnet-ibctransfer-loadtest-siftochain-2021-08-19'
    'devnet-ibctransfer-loadtest-chaintosif-2021-08-19'
    'devnet-ibctransfer-loadtest-siftochain-2021-08-22'
    'devnet-ibctransfer-loadtest-chaintosif-2021-08-22'
    'testnet-ibctransfer-loadtest-chaintosif-2021-08-22'
    'devnet-ibctransfer-loadtest-siftochain-2021-08-22-2'
    'devnet-ibctransfer-loadtest-chaintosif-2021-08-22-2'
    'testnet-ibctransfer-loadtest-siftochain-2021-08-22-2'
    'testnet-ibctransfer-loadtest-chaintosif-2021-08-22-2'
    'devnet-ibctransfer-loadtest-chaintosif-2021-08-22-3'
    'devnet-ibctransfer-loadtest-chaintosif-2021-08-22-4'
    'devnet-create-pool-2021-08-22'
    'devnet-ibc-token-to-eth'
    'sifchain-devnet-genesis-account'
    'devnet-ibctransfer-loadtest-chaintosif-2021-08-25'
    'testnet-ibctransfer-loadtest-chaintosif-2021-08-30'
)

MNEMONICS=(
    'announce world annual mushroom student eight observe panda life record radar someone amazing candy giggle time floor bottom rent expose neglect liberty violin smooth'
    'draft firm fuel eye time system cactus outside coach drama maze welcome pet okay truck garlic degree layer file solution purse what suspect level'
    'zero render knock nest sound room evoke dolphin deer beauty job torch wild west teach test holiday opera world tuition tiny pyramid bleak palm'
    'squeeze almost order behave twelve chest library upgrade citizen chaos spirit catch build dish sentence horror outer rescue voyage blanket prosper hurry issue click'
    'club stairs employ gauge chapter text elbow task mammal glue page improve section warrior disease toddler forum note wear balcony scan describe resist lift'
    'desert portion addict expose mandate again one young canoe total leopard cup clip sentence smart company match sniff town alert sick cat mobile false'
    'sing pole become zero bench oyster key decline outer tissue crunch inmate people own craft sentence attitude spend gentle scrub obvious claw solve audit'
    'give digital vintage brief normal fire forest glimpse october nuclear print question shop penalty swing inhale acid clown cup vote upset trade split planet'
    'absent catalog minute raccoon picnic gospel clog frozen donor dress ability outer topic thunder weird capable door vintage holiday garment begin museum jelly skill'
    'mouse recall romance casual walk renew minor include butter kite base couch holiday lab caught craft isolate clerk crisp fetch decide foam crazy exile'
    'antenna digital element ability donor sunset defense place turtle web ride hockey inhale price guess trap job prize way distance announce until bacon endless'
    'day film uncover certain foil avocado attack hotel near choose emotion young solid mammal guitar example crack thumb affair enable struggle ghost twelve label'
    'swing under sound uphold always fragile toss peanut fancy prison cash betray myself plate wave apart grocery whisper midnight attend bounce assume powder delay'
    'edge adjust guilt season wisdom oil dignity taxi crazy cram gain degree deer curve regular mom capable plastic taste rocket man media toy elevator'
    'shiver table reveal drill degree capable conduct labor power undo muffin mammal helmet oblige salon fog device injury sail please mother negative remind indicate'
    'disorder only cloth style bullet hybrid fan syrup drink wing fresh flat absorb install biology occur warfare fix object fuel skull misery horse rate'
    'winner flag march vivid lesson muscle ecology addict segment game impulse hungry humor damage biology spirit reunion such puzzle switch rate apology original fruit'
    'connect rocket hat athlete kind fall auction measure wage father bridge tackle midnight athlete benefit faculty shove okay win entire reveal kit era truly'
    'salute physical sand predict merry source axis film nuclear soldier reform derive miracle work student noise night cram solid desert picnic firm ranch husband'
    'imitate taxi catch add salad spike buyer slide grow element kick crucial flush husband ghost profit jump walk upgrade useful divert despair chronic any'
)

CHAINS_BINARY="sifnoded,${CHAINS_BINARY}"

for i in "${!NAMES[@]}"; do
    echo "import ${NAMES[$i]}"
    for chain_binary in ${CHAINS_BINARY//,/ }; do
        echo "-> to ${chain_binary}"
        ${chain_binary} keys delete ${NAMES[$i]} --keyring-backend test -y 2> /dev/null
        printf "%s\n\n" "${MNEMONICS[$i]}" | ${chain_binary} keys add ${NAMES[$i]} -i --recover --keyring-backend test
    done
done