#!/bin/bash

NAMES=(
    'devnet-source'
    'testnet-source'
    'devnet-peggy-loadtest-lock-eth'
    'devnet-peggy-loadtest-burn-eth'
    'testnet-peggy-loadtest-lock-eth'
    'testnet-peggy-loadtest-burn-eth'
    'devnet-peggy-loadtest-lock-rowan'
    'devnet-peggy-loadtest-burn-erowan'
    'testnet-peggy-loadtest-lock-rowan'
    'testnet-peggy-testnet-burn-erowan'
    'devnet-ibctransfer-loadtest-siftocosmos'
    'devnet-ibctransfer-loadtest-cosmostosif'
    'testnet-ibctransfer-loadtest-siftocosmos'
    'testnet-ibctransfer-loadtest-cosmostosif'
    'testnet-ibctransfer-loadtest-cosmostosif-2021-08-05'
    'testnet-ibctransfer-loadtest-cosmostosif-2021-08-09'
    'testnet-ibctransfer-loadtest-siftocosmos-2021-08-09'
    'testnet-ibctransfer-loadtest-cosmostosif-2021-08-11'
    'testnet-ibctransfer-loadtest-siftocosmos-2021-08-11'
    'testnet-peggy-loadtest-lock-rowan-2021-08-16'
    'testnet-peggy-loadtest-burn-erowan-2021-08-16'
    'devnet-ibctransfer-loadtest-cosmostosif-2021-08-17'
    'devnet-ibctransfer-loadtest-siftocosmos-2021-08-17'
    'devnet-ibctransfer-loadtest-cosmostosif-2021-08-18'
)

MNEMONICS=(
    'announce world annual mushroom student eight observe panda life record radar someone amazing candy giggle time floor bottom rent expose neglect liberty violin smooth'
    'draft firm fuel eye time system cactus outside coach drama maze welcome pet okay truck garlic degree layer file solution purse what suspect level'
    'midnight gaze inquiry urban fringe gaze afford kick shuffle spirit uniform major sibling thrive away arrest crouch govern apple mansion account tone ranch you'
    'appear pair tomato toddler iron shaft heart describe mechanic ostrich same decorate seminar soccer pair cherry play oyster sing film fetch forum awake grief'
    'busy unfold puppy insane camp royal lamp drastic eager length online lemon clown example fluid anger uphold invite balance sock cake summer novel hurry'
    'pottery video entry cage silly lonely nerve wrap tribe slab party element item apology sort broken staff afraid nice express exit hire steel genre'
    'bind acoustic rent open double adjust smile struggle feature elder tip rescue basic ordinary south giraffe quiz used wild region antenna core blind depart'
    'upper turn toe thing turkey record patch return escape beyond kitchen spoil acoustic dawn rival disease little thing just divert aspect public acid auction'
    'service glory bachelor catalog ancient zebra hazard normal elbow smooth elephant soccer donate crane glimpse yellow candy high fresh summer weasel liar bacon patch'
    'bullet lyrics captain become cruel picture danger more dumb fatal blue occur course monkey finish gravity eight injury purchase nuclear choice focus project plunge'
    'quit immune nerve tunnel salon keen fiber release mix wasp hub olive picture drift plate bind rigid where armed drift window pink bless net'
    'early steak table eight battle impulse choice giant miracle polar inflict jeans echo border shell ramp pill gaze slot laundry glove risk farm second'
    'era subway lend female slush autumn october expand raw chalk since when radar skin possible family impose ceiling clarify horror crane drop antique muffin'
    'vault zebra aisle paddle census original correct indoor recall buyer athlete sweet rare park able alcohol maximum fragile holiday alter bacon enjoy slim mosquito'
    'lemon cotton cloth winter old mesh mechanic crystal clock just pair online harsh hold hour slender square crawl you reflect tiny apple midnight exit'
    'margin start safe long wage morning toilet story barely pride music base ski delay office swap capital segment steel stool warm person dose kite'
    'cactus call vehicle civil unusual effort stem anxiety awesome dumb segment acquire decline east believe friend uncle hire upon enter search document pass improve'
    'hood chicken wise scare smile turn mix muscle ritual impact barrel dish clap view town tuna entry desert blanket off buffalo upon dash shiver'
    'shoe asset ridge skill want reopen cry soap piece warm travel suit bike coyote stadium female pond gasp donate sadness expand veteran animal couple'
    'donkey add mammal pear water raise route address hip doctor city couple brief final degree reason faint bulk dial illegal gift jaguar rib earth'
    'banner right pumpkin deer cotton vehicle hero suffer outer holiday pact news skirt jar harvest jump tide rib miss phone slam thrive rude agent'
    'second call winner rescue figure during patrol junior cram soup mutual tattoo sphere drama clutch toe hello learn grunt gauge mobile attract legal hope'
    'helmet pear urban word cricket corn cliff width obtain betray wink hero search fish slush puzzle buffalo axis drum rich never glide shiver visa'
    'coin question sleep cheap fold buzz peace today february net clap kit amount umbrella beach skirt outdoor glimpse soldier outside finish welcome pave report'
)

for i in "${!NAMES[@]}"; do
    echo "recover ${NAMES[$i]}"
    sifnoded keys delete ${NAMES[$i]} --keyring-backend test -y 2> /dev/null
    printf "%s\n\n" "${MNEMONICS[$i]}" | sifnoded keys add ${NAMES[$i]} -i --recover --keyring-backend test
done