from siftool import test_utils, cosmos, command


def get_ctx():
    return test_utils.get_env_ctx()


def test_many_balances():
    ctx = get_ctx()
    _test_many_balances(ctx, "sif16xkjwvvgg5ua48mg8kqy366xq36rz3yexf9nst")


def _test_many_balances(ctx: test_utils.EnvCtx, sif_addr: cosmos.Address):
    balances = ctx.get_sifchain_balance_long(sif_addr, height=3350)
    return
