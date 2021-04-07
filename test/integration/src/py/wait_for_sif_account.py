import sys

from test_utilities import wait_for_sif_account

print(f"sysargv: {sys.argv}")
wait_for_sif_account(sys.argv[2], "")
