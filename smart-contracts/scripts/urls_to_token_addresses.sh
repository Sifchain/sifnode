# given input of lines that have a URL in quotes,
# get the contents of the url
# extract the etherscan url
# extract the token address
#
# Input probably looks like this:

# - [ ]  [https://www.coingecko.com/en/coins/don-key](https://www.coingecko.com/en/coins/don-key)

# but will work with anything as long as there's a single url in parens

sed -e 's/.*(//' | sed -e 's/).*//' | \
  parallel -j1 wget -q -O - {} \| grep https://etherscan.io/token/ \| sed -e 's,.*etherscan.io/token/0x,0x,' \| sed -e 's,\?.*,,' \| sed -e 's,\".*,,' \| head -1
