import logging
import pytest
from integrator import cmd

@pytest.fixture(autouse=True)
def around():
    logging.info("Before...")
    logging.info("Before (debug)...")
    yield
    logging.info("After...")

def test_jure1():
    print("Print something")
    logging.debug("Debug message")
    logging.info("Info message")
    logging.warning("Warning message")
    logging.error("Error message")
    logging.info("Running test_jure1()...")
    cmd.hello()
    assert False
    # test_ebrelayer_restart.py::test_ethereum_transactions_with_offline_relayer
