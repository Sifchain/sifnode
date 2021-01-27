# yarn ganache-cli -m "race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow" -p 7545 --networkId 5777

. ../credentials.sh

yarn && yarn ganache-cli \
  -m "$ETHEREUM_ROOT_MNEMONIC" \
  -p 7545 \
  --networkId 5777 \
  -g 20000000000 \
  --gasLimit 6721975