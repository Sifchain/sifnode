import pytest

from integration_framework import main, common, eth, test_utils, inflate_tokens
from inflate_tokens import InflateTokens
from common import *


# Sifchain wallets to which we want to distribute
test_wallets = [
    "sif1fpq67nw66thzmf2a5ng64cd8p8nxa5vl9d3cm4",
    "sif1syavy2npfyt9tcncdtsdzf7kny9lh777yqc2nd",
    "sif1hjkgsq0wcmwdh8pr3snhswx5xyy4zpgs833akh",
    "sif1ypc5qcq5ha562xlak4xw3g6v352k39t6868jhx",
    "sif1u7cp5e5kty8xwuu7k234ah4jsknvkzazqagvl6",
    "sif1lj3rsayj4xtrhp2e3elv4nf7lazxty272zqegr",
    "sif1cffgyxgvw80rr6n9pcwpzrm6v8cd6dax8x32f5",
    "sif1dlse3w2pxlmuvsj5eda344zp99fegual958qyr",
    "sif1m7257566ehx7ya4ypeq7lj4y2h075z6u2xu79v",
    "sif1qrxylp97p25wcqn4cs9nd02v672073ynpkt4yr",
    "sif13rysrrdlhtmuc2pzve7jk0t4pwytwyxhaqcqcn",
    "sif1shywxv2g8gvjcqknvkxu4p6lkqhfclwwj2qk6h",
    "sif1gqm44p5ax4kgk6hksxgv4vuh2adue2acxvg542",
    "sif1zwgc9frcfpt3hhkqfu9u7up94ag5rp30kwrwrj",
]

assets = [
    {
        "decimals": 6,
        "imageUrl": "[https://assets.coingecko.com/coins/images/325/thumb/Tether-logo.png?1598003707](https://assets.coingecko.com/coins/images/325/thumb/Tether-logo.png?1598003707)",
        "name": "Tether USDT",
        "network": "sifchain",
        "symbol": "cusdt"
    }, {
        "decimals": 18,
        "imageUrl": "[https://assets.coingecko.com/coins/images/677/thumb/basic-attention-token.png?1547034427](https://assets.coingecko.com/coins/images/677/thumb/basic-attention-token.png?1547034427)",
        "name": "Basic Attention Token",
        "network": "sifchain",
        "symbol": "cbat"
    }, {
        "decimals": 18,
        "imageUrl": "[https://assets.coingecko.com/coins/images/9545/thumb/band-protocol.png?1568730326](https://assets.coingecko.com/coins/images/9545/thumb/band-protocol.png?1568730326)",
        "name": "Band Protocol",
        "network": "sifchain",
        "symbol": "cband"
    }, {
        "decimals": 18,
        "imageUrl": "[https://assets.coingecko.com/coins/images/9956/thumb/dai-multi-collateral-mcd.png?1574218774](https://assets.coingecko.com/coins/images/9956/thumb/dai-multi-collateral-mcd.png?1574218774)",
        "name": "Dai Stablecoin",
        "network": "sifchain",
        "symbol": "cdai"
    }, {
        "decimals": 18,
        "imageUrl": "[https://assets.coingecko.com/coins/images/10775/thumb/COMP.png?1592625425](https://assets.coingecko.com/coins/images/10775/thumb/COMP.png?1592625425)",
        "name": "Compound",
        "network": "sifchain",
        "symbol": "ccomp"
    }, {
        "decimals": 18,
        "imageUrl": "[https://assets.coingecko.com/coins/images/10951/thumb/UMA.png?1586307916](https://assets.coingecko.com/coins/images/10951/thumb/UMA.png?1586307916)",
        "name": "UMA",
        "network": "sifchain",
        "symbol": "cuma"
    }, {
        "decimals": 18,
        "imageUrl": "[https://assets.coingecko.com/coins/images/11683/thumb/Balancer.png?1592792958](https://assets.coingecko.com/coins/images/11683/thumb/Balancer.png?1592792958)",
        "name": "Balancer",
        "network": "sifchain",
        "symbol": "cbal"
    }, {
        "decimals": 18,
        "imageUrl": "[https://assets.coingecko.com/coins/images/11849/thumb/yfi-192x192.png?1598325330](https://assets.coingecko.com/coins/images/11849/thumb/yfi-192x192.png?1598325330)",
        "name": "yearn finance",
        "network": "sifchain",
        "symbol": "cyfi"
    }, {
        "decimals": 6,
        "imageUrl": "[https://assets.coingecko.com/coins/images/11970/thumb/serum-logo.png?1597121577](https://assets.coingecko.com/coins/images/11970/thumb/serum-logo.png?1597121577)",
        "name": "Serum",
        "network": "sifchain",
        "symbol": "csrm"
    }, {
        "decimals": 18,
        "imageUrl": "[https://assets.coingecko.com/coins/images/11976/thumb/Cream.png?1596593418](https://assets.coingecko.com/coins/images/11976/thumb/Cream.png?1596593418)",
        "name": "Cream",
        "network": "sifchain",
        "symbol": "ccream"
    }, {
        "decimals": 18,
        "imageUrl": "[https://assets.coingecko.com/coins/images/12129/thumb/sandbox_logo.jpg?1597397942](https://assets.coingecko.com/coins/images/12129/thumb/sandbox_logo.jpg?1597397942)",
        "name": "SAND",
        "network": "sifchain",
        "symbol": "csand"
    }, {
        "decimals": 18,
        "imageUrl": "[https://assets.coingecko.com/coins/images/12271/thumb/512x512_Logo_no_chop.png?1606986688](https://assets.coingecko.com/coins/images/12271/thumb/512x512_Logo_no_chop.png?1606986688)",
        "name": "Sushi",
        "network": "sifchain",
        "symbol": "csushi"
    }, {
        "decimals": 18,
        "imageUrl": "[https://assets.coingecko.com/coins/images/12440/thumb/esd_logo_circle.png?1603676421](https://assets.coingecko.com/coins/images/12440/thumb/esd_logo_circle.png?1603676421)",
        "name": "Empty Set Dollar",
        "network": "sifchain",
        "symbol": "cesd"
    }, {
        "decimals": 18,
        "imageUrl": "[https://assets.coingecko.com/coins/images/12504/thumb/uniswap-uni.png?1600306604](https://assets.coingecko.com/coins/images/12504/thumb/uniswap-uni.png?1600306604)",
        "name": "Uniswap",
        "network": "sifchain",
        "symbol": "cuni"
    }, {
        "decimals": 18,
        "imageUrl": "[https://assets.coingecko.com/coins/images/12645/thumb/AAVE.png?1601374110](https://assets.coingecko.com/coins/images/12645/thumb/AAVE.png?1601374110)",
        "name": "Aave",
        "network": "sifchain",
        "symbol": "caave"
    }, {
        "decimals": 18,
        "imageUrl": "[https://assets.coingecko.com/coins/images/14460/small/Tidal-mono.png?1616233894](https://assets.coingecko.com/coins/images/14460/small/Tidal-mono.png?1616233894)",
        "name": "Tidal",
        "network": "sifchain",
        "symbol": "ctidal"
    }, {
        "decimals": 18,
        "imageUrl": "[https://etherscan.io/token/images/dogekiller_32.png](https://etherscan.io/token/images/dogekiller_32.png)",
        "name": "DOGE KILLER",
        "network": "sifchain",
        "symbol": "cleash"
    }
]


@pytest.mark.skipif("on_peggy2_branch")
def test_inflate_tokens_short(ctx):
    amount =  12 * 10**10
    wallets = test_wallets[:2]

    # TODO Read tokens from file
    requested_tokens = [{
        "symbol": ctx.eth_symbol_to_sif_symbol(t.symbol),
        "name": t.name,
        "decimals": t.decimals,
        # Those are ignored
        # "imageUrl": None,
        # "network": None,
    } for t in [ctx.generate_random_erc20_token_data() for _ in range(3)]]

    script = InflateTokens(ctx)
    script.transfer(requested_tokens, amount, wallets)


@pytest.mark.skipif("on_peggy2_branch")
def disabled_test_inflate_tokens_full(ctx):
    amount =  12 * 10**10
    script = InflateTokens(ctx)
    script.transfer(assets, amount, test_wallets)


# TODO
# Advanced tests, such as create a new token such as mixture of existing/new tokens
# # Potato token is for Testnet and Tomato token is for Devnet.
# # They are all on Ropsten, but they are different for the sake of convenience because their smart contracts are
# # deployed by different OWNER_ADDRESS for Testnet/Devnet, and the owner is the only one with the minter role.
# potato_token = ["Potato Kilogram token", "potato", 18, "0xB51Ee40233758e9BD2bA24c4c1e9D46E272f169a"]  # Testnet owner
# tomato_token = ["Tomato Kilogram token", "tomato", 18, "0x8D753c4054046e7F77416726eEe5f3A981536B94"]  # Devnet owner, wrongly registered in token registry
# broccoli_token = ["Broccoli Kilogram token", "broc", 18, "0x77A5941E0111821ec8954555aEDfe0220bFbe798"]  # Devnet owner, registered in token registry
# carrot_token = ["Carrot Kilogram token", "carot", 18, "0x3D29b0a99d45D6ee8A669341e82eAd44e5336bDB"]  # Devnet owner, not registered in token registry
# beet_token = ["Beet Kilogram token", "beet", 18, "0x130aFa78eD832c5d871F5db02DB14Ee21AA5d803"]  # Devnet owner, not registered in token registry
