import pytest

import siftool_path
from siftool import eth, sifchain
from siftool.inflate_tokens import InflateTokens
from siftool.common import *


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
        "name": "Tether USDT",
        "symbol": "usdt"
    }, {
        "decimals": 18,
        "name": "Basic Attention Token",
        "symbol": "bat"
    }, {
        "decimals": 18,
        "name": "Band Protocol",
        "symbol": "band"
    }, {
        "decimals": 18,
        "name": "Balancer",
        "symbol": "bal"
    }, {
        "decimals": 18,
        "name": "yearn finance",
        "symbol": "yfi"
    }, {
        "decimals": 18,
        "name": "Cream",
        "symbol": "cream"
    }, {
        "decimals": 18,
        "name": "Sushi",
        "symbol": "sushi"
    }, {
        "decimals": 18,
        "name": "Uniswap",
        "symbol": "uni"
    }, {
        "decimals": 18,
        "name": "Aave",
        "symbol": "aave"
    }, {
        "decimals": 18,
        "name": "Tidal",
        "symbol": "tidal"
    }, {
        "decimals": 18,
        "name": "DOGE KILLER",
        "symbol": "leash"
    }
]


@pytest.mark.skipif("on_peggy2_branch")
def test_inflate_tokens_short(ctx):
    _test_inflate_tokens_parametrized(ctx, 3)


# This test takes >1h, times out in GitHub CI
@pytest.mark.skipif("on_peggy2_branch")
@pytest.mark.skipif("in_github_ci")
def test_inflate_tokens_long(ctx):
    _test_inflate_tokens_parametrized(ctx, 300)


def _test_inflate_tokens_parametrized(ctx, number_of_tokens):
    amount_in_tokens =  123
    amount_gwei = 456
    wallets = test_wallets[:2]

    # TODO Read tokens from file
    requested_tokens = [{
        "symbol": t.symbol,
        "name": t.name,
        "decimals": t.decimals,
    } for t in [ctx.generate_random_erc20_token_data() for _ in range(number_of_tokens)]]

    script = InflateTokens(ctx)

    balances_before = [ctx.get_sifchain_balance(w) for w in wallets]
    script.transfer(requested_tokens, amount_in_tokens, wallets, amount_gwei)
    balances_delta = [sifchain.balance_delta(balances_before[i], ctx.get_sifchain_balance(w)) for i, w in enumerate(wallets)]

    for balances_delta in balances_delta:
        for t in requested_tokens:
            assert balances_delta[ctx.eth_symbol_to_sif_symbol(t["symbol"])] == amount_in_tokens * 10**t["decimals"]
        assert balances_delta.get(ctx.ceth_symbol, 0) == amount_gwei * eth.GWEI


@pytest.mark.skipif("on_peggy2_branch")
def disabled_test_inflate_tokens_full(ctx):
    amount =  12 * 10**10
    script = InflateTokens(ctx)
    script.transfer(assets, amount, test_wallets, 0)


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
