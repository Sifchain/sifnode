from siftool import test_utils, cosmos, command


def get_ctx():
    return test_utils.get_env_ctx()


def test_many_balances():
    ctx = get_ctx()
    _test_many_balances(ctx, "sif12qykekulyevag8f7egccayg7x8a98udhf3fn26")


def _test_many_balances(ctx: test_utils.EnvCtx, sif_addr: cosmos.Address):
    balances = ctx.get_sifchain_balance_long(sif_addr)
    return
