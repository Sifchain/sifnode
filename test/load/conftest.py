import pytest
import siftool_path
import siftool.test_utils


@pytest.fixture(scope="function")
def ctx(request):
    yield from siftool.test_utils.pytest_ctx_fixture(request)


@pytest.fixture(autouse=True)
def test_wrapper_fixture():
    siftool.test_utils.pytest_test_wrapper_fixture()
