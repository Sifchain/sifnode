#!/bin/bash

NAMES=(
    'devnet-ibctransfer-loadtest-siftocosmos'
    'devnet-ibctransfer-loadtest-cosmostosif'
    'testnet-ibctransfer-loadtest-siftocosmos'
    'testnet-ibctransfer-loadtest-cosmostosif'
    'testnet-ibctransfer-loadtest-cosmostosif-2021-08-05'
    'testnet-ibctransfer-loadtest-cosmostosif-2021-08-09'
    'testnet-ibctransfer-loadtest-siftocosmos-2021-08-09'
    'testnet-ibctransfer-loadtest-cosmostosif-2021-08-11'
    'testnet-ibctransfer-loadtest-siftocosmos-2021-08-11'
    'devnet-ibctransfer-loadtest-cosmostosif-2021-08-17'
    'devnet-ibctransfer-loadtest-siftocosmos-2021-08-17'
    'devnet-ibctransfer-loadtest-cosmostosif-2021-08-18'
)

MNEMONICS=(
    'defy security coast rebuild wealth common horn ignore ready fade glare stove frame decrease diagram fuel only love siren into reform artefact silly lecture'
    'loop just iron lizard hurdle humble barrel average blue cruel intact top divide axis allow achieve put swear enforce riot soda abuse rude hello'
    'leg struggle minor powder bronze belt diagram patient skull cage device smoke charge army pretty wage infant bone clap tissue water walk innocent tip'
    'wrap mystery panda bird segment kangaroo among country february scrap bounce clip rotate angle auction venue acoustic motion armor sheriff wild six accident stuff'
    'frame slice other bless safe describe live lend office cup patch fame liberty jealous security calm six cereal deal guide motor never magnet despair'
    'someone chaos admit print bounce force zoo frequent weapon cat letter alarm magnet seed burden abuse spell gauge peanut mutual wall spirit alpha document'
    'drink language provide shrug surprise window olive useless guitar industry drama much comfort dismiss fragile beach simple stay derive energy mouse enforce title chat'
    'choice swim actress april reunion head spell best below october faith beef weather noble silly bulk opinion accident empower stable scout life siren miracle'
    'off forest caution mom desk benefit wasp whale prevent dinosaur educate donor empower auto action path warfare venue cattle frame enjoy phrase host slice'
    'obey inquiry eyebrow mad canoe volcano lounge alpha found solution remind defy orchard blind thunder chalk vote august asset dentist neutral gentle reward sheriff'
    'idea horse page parrot slab butter industry glow never profit barrel push almost holiday post aim rent large typical primary town off energy ivory'
    'subway system fold gospel slice allow bus bridge control empty build tool alien holiday usage velvet paddle such hundred series vague key repair pink'
)

for i in "${!NAMES[@]}"; do
    echo "recover ${NAMES[$i]}"
    gaiad keys delete ${NAMES[$i]} --keyring-backend test -y 2> /dev/null
    printf "%s\n\n" "${MNEMONICS[$i]}" | gaiad keys add ${NAMES[$i]} -i --recover --keyring-backend test
done