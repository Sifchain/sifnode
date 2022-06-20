import pytest
import siftool_path
from siftool import test_utils


@pytest.fixture(scope="function")
def ctx(request):
    yield from test_utils.pytest_ctx_fixture(request)


@pytest.fixture(autouse=True)
def test_wrapper_fixture():
    test_utils.pytest_test_wrapper_fixture()
