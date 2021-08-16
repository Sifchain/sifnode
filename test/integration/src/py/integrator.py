import logging
import pytest


class Integrator:
    def hello(self):
        print("Hello!")
        logging.debug("Hello [debug]")

cmd = Integrator()

@pytest.fixture
def integration_env():
    return cmd
