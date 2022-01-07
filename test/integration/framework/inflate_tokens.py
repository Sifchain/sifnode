# This is a replacement for test/integration/inflate_tokens.sh.
# The original script had a lot of problems as described in https://app.zenhub.com/workspaces/current-sprint---engineering-615a2e9fe2abd5001befc7f9/issues/sifchain/issues/719.
# See https://www.notion.so/sifchain/TEST-TOKEN-DISTRIBUTION-PROCESS-41ad0861560c4be58918838dbd292497

import json
import logging

import test_utils
from common import *

log = logging.getLogger(__name__)


class InflateTokens:
    def __init__(self, ctx):
        self.ctx = ctx
        self.wait_for_account_change_timeout = 1800  # For Ropsten we need to wait for 50 blocks i.e. ~20 mins

    def get_whitelisted_tokens(self):
        whitelist = self.ctx.get_whitelisted_tokens_from_bridge_bank_past_events()
        result = []
        for token_addr, value in whitelist.items():
            token_data = self.ctx.get_generic_erc20_token_data(token_addr)
            token_symbol = token_data["symbol"]
            token = {
                "address": token_addr,
                "symbol": token_symbol,
                "name": token_data["name"],
                "decimals": token_data["decimals"],
                "is_whitelisted": value,
                "sif_denom": self.ctx.eth_symbol_to_sif_symbol(token_symbol),
            }
            log.debug("Whitelisted entry: {}".format(repr(token_data)))
            assert token_symbol not in result, f"Symbol {token_symbol} is being used by more than one whitelisted token"
            result.append(token)
        erowan_token = [t for t in result if t["symbol"] == "erowan"]
        assert len(erowan_token) == 1, "erowan is not whitelisted"
        assert erowan_token[0]["is_whitelisted"], "erowan is un-whitelisted"
        return result

    def wait_for_all(self, pending_txs):
        result = []
        for txhash in pending_txs:
            txrcpt = self.ctx.eth.wait_for_transaction_receipt(txhash)
            result.append(txrcpt)
        return result

    def build_list_of_tokens_to_create(self, existing_tokens, requested_tokens):
        """
        This part deploys SifchainTestoken for every requested token that has not yet been deployed.
        The list of requested tokens is (historically) read from assets.json, but in practice it can be
        a subset of tokens that are whitelisted in production.
        The list of existing tokens is reconstructed from past LogWhiteListUpdate events of the BridgeBank
        smart contract (since there is no way to "dump" the contents of a mapping in Solidity).
        Deployed tokens are whitelisted with BridgeBank, minted to owner's account and approved to BridgeBank.
        This part only touches EVM chain through web3.
        """

        # Strictly speaking we could also skip tokens that were un-whitelisted (value == False) since the fact that
        # their addresses appear in BridgeBank's past events implies that the corresponding ERC20 smart contracts have
        # been deployed, hence there is no need to deploy them.

        # This assumes that requested token symbols are in sifchain format (c-prefixed, i.e. "cusdt", "csushi" etc.).
        # There can also be "ceth" and "rowan" in this list, which we ignore as they represent special cases.
        # To compare it to entries on existing_whitelist, we need to prefix entries on existing_whitelist with "c".
        # TODO It would be better if the requested tokens didn't have "c" prefixes. For now we keep it for
        #      compatibility. Ask people who use this script.
        token_symbols_to_skip = set()
        token_symbols_to_skip.add(test_utils.CETH)  # ceth is special since we can't just mint it or create an ERC20 contract for it
        token_symbols_to_skip.add(test_utils.ROWAN)
        tokens_to_create = []  # = requested - existing - {rowan, ceth}
        for token in requested_tokens:
            token_symbol = token["symbol"]
            if (token_symbol == test_utils.CETH) or (token_symbol == test_utils.ROWAN):
                assert False, f"Token {token_symbol} cannot be used by this procedure, please remove it from list of requested assets"
            if not token_symbol.startswith("c"):
                assert False, f"Token {token_symbol} is invalid - should start with 'c'"
            eth_token_symbol = token_symbol[1:]  # Strip "c", e.g. "cusdt" -> "usdt"

            existing_token = zero_or_one(find_by_value(existing_tokens, "symbol", eth_token_symbol))
            if existing_token is None:
                tokens_to_create.append({
                    "name": token["name"],
                    "symbol": eth_token_symbol,
                    "decimals": token["decimals"],
                })
            else:
                if not all([existing_token[f] == token[f] for f in ["name", "decimals"]]):
                    assert False, "Existing token's name/decimals does not match requested for token: " \
                        "requested={}, existing={}".format(repr(token), repr(existing_token))
                if existing_token["is_whitelisted"]:
                    log.info(f"Skipping deployment of smmart contract for token {token_symbol} as it should already exist")
                else:
                    log.warning(f"Skipping token {token_symbol} as it is currently un-whitelisted")
        return tokens_to_create

    def create_new_tokens(self, tokens_to_create):
        amount_in_token_units = 0
        pending_txs = []
        for token in tokens_to_create:
            token_name = token["name"]
            token_symbol = token["symbol"]
            token_decimals = token["decimals"]
            log.info(f"Creating token {token_symbol}...")
            amount = amount_in_token_units * (10**token_decimals)
            # Deploy a SifchainTestToken
            # call BridgeBank.updateEthWhiteList with its address
            # Mint amount_in_token_units to operator_address
            # Approve entire minted amount to BridgeBank
            # TODO We don't really need create_new_currency here, we only need to deploy the smart contract
            #      since we do the minting and approval in next step (token_refresh).

            # token_addr = self.ctx.create_new_currency(token_symbol, token_name, token_decimals, amount, minted_tokens_recipient)

            txhash = self.ctx.tx_deploy_new_generic_erc20_token(self.ctx.operator, token_name, token_symbol, token_decimals)
            pending_txs.append(txhash)

        token_contracts = [self.ctx.tx_get_generic_erc20_token_at(txrcpt.contractAddress) for txrcpt in self.wait_for_all(pending_txs)]

        new_tokens = []
        pending_txs = []
        for token_to_create, token_sc in [[tokens_to_create[i], c] for i, c in enumerate(token_contracts)]:
            token_symbol = token_to_create["symbol"]
            token_name = token_to_create["name"]
            token_decimals = token_to_create["decimals"]
            assert token_sc.functions.totalSupply().call() == 0
            assert token_sc.functions.name().call() == token_name
            assert token_sc.functions.symbol().call() == token_symbol
            assert token_sc.functions.decimals().call() == token_decimals
            new_tokens.append({
                "address": token_sc.address,
                "symbol": token_symbol,
                "name": token_name,
                "decimals": token_decimals,
                "is_whitelisted": True,
                "sif_denom": self.ctx.eth_symbol_to_sif_symbol(token_symbol),
            })
            txhash = self.ctx.tx_update_bridge_bank_whitelist(token_sc.address, True)
            pending_txs.append(txhash)

        self.wait_for_all(pending_txs)
        return new_tokens

    def mint(self, list_of_tokens_addrs, amount_in_tokens, mint_recipient):
        pending_txs = []
        for token_addr in list_of_tokens_addrs:
            token_sc = self.ctx.tx_get_generic_erc20_token_at(token_addr)
            decimals = token_sc.functions.decimals().call()
            amount = amount_in_tokens * 10**decimals
            txhash = self.ctx.tx_testing_token_mint(token_sc, self.ctx.operator, amount, mint_recipient)
            pending_txs.append(txhash)
        self.wait_for_all(pending_txs)

    def approve_and_lock(self, token_addr_list, eth_addr, to_sif_addr, amount):
        pending_txs = []
        for token_addr in token_addr_list:
            token_sc = self.ctx.tx_get_generic_erc20_token_at(token_addr)
            pending_txs.extend(self.ctx.tx_approve_and_lock(token_sc, eth_addr, to_sif_addr, amount))
        return self.wait_for_all(pending_txs)

    def transfer_from_eth_to_sifnode(self, eth_addr, sif_addr, tokens_to_transfer, amount):
        sif_balances_before = self.ctx.get_sifchain_balance(sif_addr)
        self.approve_and_lock([t["address"] for t in tokens_to_transfer], eth_addr, sif_addr, amount)

        # Wait for intermediate_sif_account to receive all funds across the bridge
        self.ctx.advance_blocks()
        send_amounts = [[amount, t["sif_denom"]] for t in tokens_to_transfer]
        self.ctx.wait_for_sif_balance_change(sif_addr, sif_balances_before,
            min_changes=send_amounts, timeout=self.wait_for_account_change_timeout)

    def distribute_tokens_to_wallets(self, from_sif_account, tokens_to_transfer, amount, target_sif_accounts):
        # Distribute from intermediate_sif_account to each individual account
        # Note: firing transactions with "sifnoded tx bank send" in rapid succession does not work. This is currently a
        # known limitation of Cosmos SDK, see https://github.com/cosmos/cosmos-sdk/issues/4186
        # Instead, we take advantage of batching multiple denoms to single account with single send command (amounts
        # separated by by comma: "sifnoded tx bank send ... 100denoma,100denomb,100denomc") and wait for destination
        # account to show changes for all denoms after each send.
        send_amounts = [[amount, t["sif_denom"]] for t in tokens_to_transfer]
        target_sif_balances_before = []
        for sif_acct in target_sif_accounts:
            target_sif_balances_before.append(self.ctx.get_sifchain_balance(sif_acct))
            sif_balance_before = self.ctx.get_sifchain_balance(sif_acct)
            self.ctx.send_from_sifchain_to_sifchain(from_sif_account, sif_acct, send_amounts)
            self.ctx.wait_for_sif_balance_change(sif_acct, sif_balance_before, min_changes=send_amounts)

    def run(self, requested_tokens, amount, target_sif_accounts):
        """
        It goes like this:
        1. Starting with assets.json of your choice, It will first compare the list of tokens to existing whitelist and deploy any new tokens (ones that have not yet been whitelisted)
        2. For each token in assets.json It will mint the given amount of all listed tokens to OPERATOR account
        3. It will do a single transaction across the bridge to move all tokens from OPERATOR to sif_broker_account
        4. It will distribute tokens from sif_broker_account to each of given target accounts
        The sif_broker_account and OPERATOR can be any Sifchain and Ethereum accounts, we might want to use something
        familiar so that any tokens that would get stuck in the case of interrupting the script can be recovered.
        """

        assert not on_peggy2_branch, "Not supported yet on peggy2.0 branch"

        self.ctx.sanity_check()

        amount_per_token = amount * len(target_sif_accounts)
        # sif_broker_account = self.ctx.rowan_source
        fund_rowan = [5 * test_utils.sifnode_funds_for_transfer_peggy1, "rowan"]
        sif_broker_account = self.ctx.create_sifchain_addr(fund_amounts=[fund_rowan])
        eth_broker_account = self.ctx.operator

        log.info("Using eth_broker_account {}".format(eth_broker_account))
        log.info("Using sif_broker_account {}".format(sif_broker_account))

        # Check first that we have the key for ROWAN_SOURCE since the script uses it as an intermediate address
        keys = self.ctx.cmd.sifnoded_keys_list(keyring_backend="test", sifnoded_home=self.ctx.sifnoded_home)
        rowan_source_key = zero_or_one([k for k in keys if k["address"] == sif_broker_account])
        assert rowan_source_key is not None, "Need private key of broker account {} in sifnoded test keyring".format(sif_broker_account)

        existing_tokens = self.get_whitelisted_tokens()

        tokens_to_create = self.build_list_of_tokens_to_create(existing_tokens, requested_tokens)
        new_tokens = self.create_new_tokens(tokens_to_create)
        existing_tokens.extend(new_tokens)

        tokens_to_transfer = [exactly_one(find_by_value(existing_tokens, "sif_denom", t["symbol"]))
            for t in requested_tokens]

        self.mint([t["address"] for t in tokens_to_transfer], amount_per_token, eth_broker_account)
        self.transfer_from_eth_to_sifnode(eth_broker_account, sif_broker_account, tokens_to_transfer, amount_per_token)
        self.distribute_tokens_to_wallets(sif_broker_account, tokens_to_transfer, amount, target_sif_accounts)


def run(assets_file, amount, target_accounts_file):
    ctx = test_utils.get_env_ctx()
    requested_tokens = json.loads(ctx.cmd.read_text_file(assets_file))
    target_accounts = json.loads(ctx.cmd.read_text_file(target_accounts_file))
    amount = int(amount)
    script = InflateTokens(ctx)
    script.run(requested_tokens, amount, target_accounts)

if __name__ == "__main__":
    import sys
    basic_logging_setup()
    run(*sys.argv[1:])
