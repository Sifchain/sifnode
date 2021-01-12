import burn_lock_functions
from burn_lock_functions import SifchaincliCredentials
from test_utilities import get_required_env_var, get_shell_output


def sifchain_cli_credentials_for_test(key: str):
    return SifchaincliCredentials(
        keyring_passphrase="",
        keyring_backend="test",
        from_key=key,
        sifnodecli_homedir=f"""{get_required_env_var("HOME")}/.sifnodecli"""
    )
