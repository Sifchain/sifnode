import json
import logging
import sys

import burn_lock_functions
from test_utilities import get_required_env_var

if __name__ == "__main__":
    args = burn_lock_functions.transfer_argument_parser().parse_args()
    burn_lock_functions.configure_logging(args)

    logging.debug(f"command line arguments: {sys.argv} {args}")

    request = burn_lock_functions.args_to_EthereumToSifchainTransferRequest(args)

    logging.info(f"transferrequestjson: {json.dumps(request.__dict__)}")

    credentials = burn_lock_functions.SifchaincliCredentials(
        get_required_env_var("OWNER_PASSWORD"),
        from_key="user1",
        homedir=get_required_env_var("CHAINDIR") + "/.sifnodecli"
    )

    transfer_result = burn_lock_functions.transfer_sifchain_to_sifchain(request, credentials)
    logging.debug(f"transfer_result is: {transfer_result}")
    final_balance = transfer_result["sifchain_ending_balance"]

    result = json.dumps({
        "final_balance": final_balance,
        "final_balance_10_18": float(final_balance) / (10 ** 18),
        "transfer_request": request.__dict__,
        "logfile": args.logfile[0],
        "steps": transfer_result
    })

    logging.info(f"transferresultjson: {result}")

    print(result)
