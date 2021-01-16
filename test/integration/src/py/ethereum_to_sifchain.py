import json
import logging
import sys

import burn_lock_functions

if __name__ == "__main__":
    arg_parser = burn_lock_functions.transfer_argument_parser()
    args = burn_lock_functions.add_credentials_arguments(arg_parser).parse_args()
    burn_lock_functions.configure_logging(args)

    logging.debug(f"command line arguments: {sys.argv} {args}")

    request = burn_lock_functions.EthereumToSifchainTransferRequest.from_args(args)

    logging.info(f"transferrequestjson: {json.dumps(request.__dict__)}")

    transfer_result = burn_lock_functions.transfer_ethereum_to_sifchain(request, max_seconds=30)
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
