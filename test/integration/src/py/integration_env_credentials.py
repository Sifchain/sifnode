import burn_lock_functions
from burn_lock_functions import SifchaincliCredentials
from test_utilities import get_required_env_var, get_shell_output


def sifchain_cli_credentials_for_test(key: str) -> SifchaincliCredentials:
    """Returns SifchaincliCredentials for the test keyring with from_key set to key"""
    return SifchaincliCredentials(
        keyring_passphrase="",
        keyring_backend="test",
        from_key=key,
        sifnodecli_homedir=f"""{get_required_env_var("HOME")}/.sifnodecli"""
    )


def create_new_sifaddr_and_credentials() -> (str, SifchaincliCredentials):
    new_account_key = get_shell_output("uuidgen")
    credentials = sifchain_cli_credentials_for_test(new_account_key)
    new_addr = burn_lock_functions.create_new_sifaddr(credentials=credentials, keyname=new_account_key)
    credentials.from_key = new_addr["name"]
    return new_addr["address"], credentials,
