# This is a replacement for test/integration/inflate_tokens.sh.
# The original script had a lot of problems as described in https://app.zenhub.com/workspaces/current-sprint---engineering-615a2e9fe2abd5001befc7f9/issues/sifchain/issues/719.
# See https://www.notion.so/sifchain/TEST-TOKEN-DISTRIBUTION-PROCESS-41ad0861560c4be58918838dbd292497

import json
import re
from typing import Any, Mapping, Iterable, Sequence

from siftool import eth, test_utils, cosmos
from siftool.common import *

log = siftool_logger(__name__)

TokenDict = Mapping[str, Any]


class InflateTokens:
    def __init__(self, ctx: test_utils.EnvCtx):
        self.ctx = ctx
        self.wait_for_account_change_timeout = 120
        self.excluded_token_symbols = ["erowan"]  # TODO peggy1 only

        # Only transfer this tokens in a batch for Peggy1. See #2397. You would need to adjust this if
        # test_inflate_tokens_short is passing, but test_inflate_tokens_long is timing out. It only applies to Peggy 1.
        # The value of 3 is experimental; if tokens are still not getting across the bridge reliably, reduce the value
        # down to 1 (minimum). The lower the value the more time the transfers will take as there will be more
        # sequential transfers instead of parallel.
        self.max_ethereum_batch_size = 0

    def get_whitelisted_tokens(self) -> List[TokenDict]:
        whitelist = self.ctx.get_whitelisted_tokens_from_bridge_bank_past_events()
        ibc_pattern = re.compile("^ibc/([0-9a-fA-F]{64})$")
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
            m = ibc_pattern.match(token_symbol)
            if m:
                token["ibc"] = m[1].lower()
            log.debug("Found whitelisted entry: {}".format(repr(token_data)))
            assert token_symbol not in result, f"Symbol {token_symbol} is being used by more than one whitelisted token"
            result.append(token)
        erowan_token = [t for t in result if t["symbol"] == "erowan"]
        # These assertions are broken in Tempnet, possibly indicating missing/incomplete chain init, see README.md
        # for comparision of steps
        assert len(erowan_token) == 1, "erowan is not whitelisted, probably bad/incomplete deployment"
        assert erowan_token[0]["is_whitelisted"], "erowan is un-whitelisted"
        return result

    def wait_for_all(self, pending_txs):
        result = []
        for txhash in pending_txs:
            txrcpt = self.ctx.eth.wait_for_transaction_receipt(txhash)
            result.append(txrcpt)
        return result

    def build_list_of_tokens_to_create(self, existing_tokens: Iterable[TokenDict], requested_tokens: Iterable[TokenDict]
    ) -> Sequence[Mapping[str, Any]]:
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

        tokens_to_create = []
        for token in requested_tokens:
            token_symbol = token["symbol"]
            if token_symbol in self.excluded_token_symbols:
                assert False, f"Token {token_symbol} cannot be used by this procedure, please remove it from list of requested assets"

            existing_token = zero_or_one(find_by_value(existing_tokens, "symbol", token_symbol))
            if existing_token is None:
                tokens_to_create.append(token)
            else:
                if not all(existing_token[f] == token[f] for f in ["name", "decimals"]):
                    assert False, "Existing token's name/decimals does not match requested for token: " \
                        "requested={}, existing={}".format(repr(token), repr(existing_token))
                if existing_token["is_whitelisted"]:
                    log.info(f"Skipping deployment of smmart contract for token {token_symbol} as it should already exist")
                else:
                    log.warning(f"Skipping token {token_symbol} as it is currently un-whitelisted")
        return tokens_to_create

    def create_new_tokens(self, tokens_to_create: Iterable[TokenDict]) -> Sequence[TokenDict]:
        pending_txs = []
        for token in tokens_to_create:
            token_name = token["name"]
            token_symbol = token["symbol"]
            token_decimals = token["decimals"]
            log.info(f"Deploying generic ERC20 smart contract for token {token_symbol}...")
            txhash = self.ctx.tx_deploy_new_generic_erc20_token(self.ctx.operator, token_name, token_symbol, token_decimals)
            pending_txs.append(txhash)

        token_contracts = [self.ctx.get_generic_erc20_sc(txrcpt.contractAddress) for txrcpt in self.wait_for_all(pending_txs)]

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
            if not on_peggy2_branch:
                txhash = self.ctx.tx_update_bridge_bank_whitelist(token_sc.address, True)
                pending_txs.append(txhash)

        self.wait_for_all(pending_txs)
        return new_tokens

    def mint(self, list_of_tokens_addrs, amount_in_tokens, mint_recipient):
        pending_txs = []
        for token_addr in list_of_tokens_addrs:
            token_sc = self.ctx.get_generic_erc20_sc(token_addr)
            decimals = token_sc.functions.decimals().call()
            amount = amount_in_tokens * 10**decimals
            txhash = self.ctx.tx_testing_token_mint(token_sc, self.ctx.operator, amount, mint_recipient)
            pending_txs.append(txhash)
        self.wait_for_all(pending_txs)

    def transfer_from_eth_to_sifnode(self, from_eth_addr, to_sif_addr, tokens_to_transfer, amount_in_tokens, amount_eth_gwei):
        sif_balances_before = self.ctx.get_sifchain_balance(to_sif_addr)
        sent_amounts = []
        pending_txs = []
        for token in tokens_to_transfer:
            token_addr = token["address"]
            decimals = token["decimals"]
            token_sc = self.ctx.get_generic_erc20_sc(token_addr)
            amount = amount_in_tokens * 10**decimals
            pending_txs.extend(self.ctx.tx_approve_and_lock(token_sc, from_eth_addr, to_sif_addr, amount))
            sent_amounts.append([amount, token["sif_denom"]])
        if amount_eth_gwei > 0:
            amount = amount_eth_gwei * eth.GWEI
            pending_txs.append(self.ctx.tx_bridge_bank_lock_eth(from_eth_addr, to_sif_addr, amount))
            sent_amounts.append([amount, self.ctx.ceth_symbol])
        self.wait_for_all(pending_txs)
        log.info("{} Ethereum transactions commited: {}".format(len(pending_txs), repr(sent_amounts)))

        # Wait for intermediate_sif_account to receive all funds across the bridge
        previous_block = self.ctx.eth.w3_conn.eth.block_number
        self.ctx.advance_blocks()
        log.info("Ethereum blocks advanced by {}".format(self.ctx.eth.w3_conn.eth.block_number - previous_block))
        self.ctx.sifnode.wait_for_balance_change(to_sif_addr, sif_balances_before, min_changes=sent_amounts,
            polling_time=5, timeout=0, change_timeout=self.wait_for_account_change_timeout)

    # Distributes from intermediate_sif_account to each individual account
    def distribute_tokens_to_wallets(self, from_sif_account, tokens_to_transfer, amount_in_tokens, target_sif_accounts, amount_eth_gwei):
        send_amounts = [[amount_in_tokens * 10**t["decimals"], t["sif_denom"]] for t in tokens_to_transfer]
        if amount_eth_gwei > 0:
            send_amounts.append([amount_eth_gwei * eth.GWEI, self.ctx.ceth_symbol])
        progress_total = len(target_sif_accounts) * len(send_amounts)
        progress_current = 0
        for sif_acct in target_sif_accounts:
            remaining = send_amounts
            while remaining:
                batch_size = len(remaining)
                if (self.ctx.sifnode.max_send_batch_size > 0) and (batch_size > self.ctx.sifnode.max_send_batch_size):
                    batch_size = self.ctx.sifnode.max_send_batch_size
                batch = remaining[:batch_size]
                remaining = remaining[batch_size:]
                sif_balance_before = self.ctx.get_sifchain_balance(sif_acct)
                self.ctx.send_from_sifchain_to_sifchain(from_sif_account, sif_acct, batch)
                self.ctx.sifnode.wait_for_balance_change(sif_acct, sif_balance_before, min_changes=batch,
                    polling_time=2, timeout=0, change_timeout=self.wait_for_account_change_timeout)
                progress_current += batch_size
                log.debug("Distributing tokens to wallets: {:0.0f}% done".format((progress_current/progress_total) * 100))

    def export(self):
        return [{
            "symbol": token["symbol"],
            "name": token["name"],
            "decimals": token["decimals"]
        } for token in self.get_whitelisted_tokens() if ("ibc" not in token) and (token["symbol"] not in self.excluded_token_symbols)]

    def transfer(self, requested_tokens: Sequence[TokenDict], token_amount: int,
        target_sif_accounts: Sequence[cosmos.Address], eth_amount_gwei: int
    ):
        """
        It goes like this:
        1. Starting with assets.json of your choice, It will first compare the list of tokens to existing whitelist and deploy any new tokens (ones that have not yet been whitelisted)
        2. For each token in assets.json It will mint the given amount of all listed tokens to OPERATOR account
        3. It will do a single transaction across the bridge to move all tokens from OPERATOR to sif_broker_account
        4. It will distribute tokens from sif_broker_account to each of given target accounts
        The sif_broker_account and OPERATOR can be any Sifchain and Ethereum accounts, we might want to use something
        familiar so that any tokens that would get stuck in the case of interrupting the script can be recovered.
        """

        # TODO Add support for "rowan"

        n_accounts = len(target_sif_accounts)
        total_token_amount = token_amount * n_accounts
        total_eth_amount_gwei = eth_amount_gwei * n_accounts

        # Calculate how much rowan we need to fund intermediate account with. This is only an estimation at this point.
        # We need to take into account that we might need to break transfers in batches. The number of tokens is the
        # number of ERC20 tokens plus one for ETH, rounded up. 5 is a safety factor
        number_of_batches = 1 if self.ctx.sifnode.max_send_batch_size == 0 else (len(requested_tokens) + 1) // self.ctx.sifnode.max_send_batch_size + 1
        fund_rowan = [5 * test_utils.sifnode_funds_for_transfer_peggy1 * n_accounts * number_of_batches, "rowan"]
        log.debug("Estimated number of batches needed to transfer tokens from intermediate sif account to target sif wallet: {}".format(number_of_batches))
        log.debug("Estimated rowan funding needed for intermediate account: {}".format(fund_rowan))
        ether_faucet_account = self.ctx.operator
        sif_broker_account = self.ctx.create_sifchain_addr(fund_amounts=[fund_rowan])
        eth_broker_account = self.ctx.operator

        if (total_eth_amount_gwei > 0) and (ether_faucet_account != eth_broker_account):
            self.ctx.eth.send_eth(ether_faucet_account, eth_broker_account, total_eth_amount_gwei)

        log.info("Using eth_broker_account {}".format(eth_broker_account))
        log.info("Using sif_broker_account {}".format(sif_broker_account))

        # Check first that we have the key for ROWAN_SOURCE since the script uses it as an intermediate address
        keys = self.ctx.sifnode.keys_list()
        rowan_source_key = zero_or_one([k for k in keys if k["address"] == sif_broker_account])
        assert rowan_source_key is not None, "Need private key of broker account {} in sifnoded test keyring".format(sif_broker_account)

        existing_tokens = self.get_whitelisted_tokens()
        tokens_to_create = self.build_list_of_tokens_to_create(existing_tokens, requested_tokens)
        log.info("Existing tokens: {}".format(len(existing_tokens)))
        log.info("Requested tokens: {}".format(len(requested_tokens)))
        log.info("Tokens to create: {}".format(len(tokens_to_create)))

        new_tokens = self.create_new_tokens(tokens_to_create)
        existing_tokens.extend(new_tokens)

        # At this point, all tokens that we want to transfer should exist both on Ethereum blockchain as well as in
        # existing_tokens.
        tokens_to_transfer = [exactly_one(find_by_value(existing_tokens, "symbol", t["symbol"]))
            for t in requested_tokens]

        self.mint([t["address"] for t in tokens_to_transfer], total_token_amount, eth_broker_account)

        if (self.max_ethereum_batch_size > 0) and (len(tokens_to_transfer) > self.max_ethereum_batch_size):
            log.debug(f"Transferring {len(tokens_to_transfer)} tokens from ethereum to sifndde in batches of {self.max_ethereum_batch_size}...")
            remaining = tokens_to_transfer
            while remaining:
                batch = remaining[:self.max_ethereum_batch_size]
                remaining = remaining[self.max_ethereum_batch_size:]
                self.transfer_from_eth_to_sifnode(eth_broker_account, sif_broker_account, batch, total_token_amount, 0)
                log.debug(f"Batch completed, {len(remaining)} tokens remaining")
            # Transfer ETH separately
            log.debug("Thansfering ETH from ethereum to sifnode...")
            self.transfer_from_eth_to_sifnode(eth_broker_account, sif_broker_account, [], 0, total_eth_amount_gwei)
        else:
            log.debug(f"Transferring {len(tokens_to_transfer)} tokens from ethereum to sifnode in single batch...")
            self.transfer_from_eth_to_sifnode(eth_broker_account, sif_broker_account, tokens_to_transfer, total_token_amount, total_eth_amount_gwei)
        self.distribute_tokens_to_wallets(sif_broker_account, tokens_to_transfer, token_amount, target_sif_accounts, eth_amount_gwei)

        log.info("Done.")
        log.info("To see newly minted tokens in UI, you need to edit 'scripts/ibc/tokenregistry/generate-erc20-jsons.sh' "
            "and add any tokens that are not already there. Then cd into the directory and run './generate-erc20-jsons.sh devnet' "\
            "and commit the results in the sifchain-devnet-1 folder. @tim will pick up the PR and register it on "
            "devnet by running './register-one.sh' with the registry key. In the future this might be open for anybody "
            "to do on their own for devnet and testnet.")

    def transfer_eth(self, from_eth_addr: eth.Address, amount_gwei: int, target_sif_accounts: Iterable[cosmos.Address]):
        pending_txs = []
        for sif_acct in target_sif_accounts:
            txrcpt = self.ctx.tx_bridge_bank_lock_eth(from_eth_addr, sif_acct, amount_gwei * eth.GWEI)
            pending_txs.append(txrcpt)
        self.wait_for_all(pending_txs)


def run(*args):
    # This script should be run with SIFTOOL_ENV_FILE set to a file containing definitions for OPERATOR_ADDRESS,
    # ROWAN_SOURCE eth. Depending on if you're running it on Peggy1 or Peggy2 the format might be different.
    # See get_env_ctx() for details.
    assert not on_peggy2_branch, "Not supported yet on peggy2.0 branch"
    ctx = test_utils.get_env_ctx()
    script = InflateTokens(ctx)
    script.wait_for_account_change_timeout = 1800  # For Ropsten we need to wait for 50 blocks i.e. ~20 min = 1200 s
    cmd = args[0]
    args = args[1:]
    if cmd == "export":
        # Usage: inflate_tokens.py export assets.json
        ctx.cmd.write_text_file(args[0], json.dumps(script.export(), indent=4))
    elif cmd == "transfer":
        # Usage: inflate_tokens.py transfer assets.json token_amount accounts.json amount_eth_gwei
        assets_json_file, token_amount, accounts_json_file, amount_eth_gwei = args
        tokens = json.loads(ctx.cmd.read_text_file(assets_json_file))
        accounts = json.loads(ctx.cmd.read_text_file(accounts_json_file))
        script.transfer(tokens, int(token_amount), accounts, int(amount_eth_gwei))
    else:
        raise Exception("Invalid usage")


if __name__ == "__main__":
    import sys
    basic_logging_setup()
    run(*sys.argv[1:])
