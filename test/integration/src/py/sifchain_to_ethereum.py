import json
import logging
import sys

import burn_lock_functions


def test_nothing():
    return True


if __name__ == "__main__":
    rc = burn_lock_functions.process_args(sys.argv, __file__)

    request = rc.transfer_request

    logging.info(f"transfer_request_json: {json.dumps(request.__dict__)}")

    credentials = rc.credentials

    transfer_result = burn_lock_functions.transfer_sifchain_to_ethereum(request, credentials)

    final_balance = transfer_result["ethereum_ending_balance"]

    result = json.dumps({
        "final_balance": final_balance,
        "final_balance_10_18": float(final_balance) / (10 ** 18),
        "transfer_request": request.__dict__,
        "logfile": args.logfile[0],
        "steps": transfer_result
    })

    logging.info(f"transfer_result_json: {result}")

    print(result)
