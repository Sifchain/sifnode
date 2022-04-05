import logging
import pytest
import siftool_path
from siftool import test_utils


@pytest.fixture(scope="function")
def ctx(request):
    # To pass the "snapshot_name" as a parameter with value "foo" from test, annotate the test function like this:
    # @pytest.mark.snapshot_name("foo")
    snapshot_name = request.node.get_closest_marker("snapshot_name")
    if snapshot_name is not None:
        snapshot_name = snapshot_name.args[0]
    logging.error("Context setup: snapshot_name={}".format(repr(snapshot_name)))
    with test_utils.get_test_env_ctx() as ctx:
        yield ctx
        logging.debug("Test context cleanup")
