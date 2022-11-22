import re
import os
import logging
import subprocess
import string
import random
import time
import yaml
import urllib.request
from typing import Optional, Mapping, Sequence, Set, IO, Union, Iterable, List, Any, Callable, Dict


ANY_ADDR = "0.0.0.0"
LOCALHOST = "127.0.0.1"
JsonObj = Any
JsonDict = Dict[str, JsonObj]


def stdout(res):
    return res[1]

def stderr(res):
    return res[2]

def stdout_lines(res):
    return stdout(res).splitlines()

def joinlines(lines):
    return "".join([x + os.linesep for x in lines])

def zero_or_one(items: Sequence[Any]) -> Any:
    if len(items) == 0:
        return None
    elif len(items) > 1:
        raise ValueError("Multiple items")
    else:
        return items[0]

def exactly_one(items: Union[Sequence[Any], Set[Any]]) -> Any:
    if len(items) == 0:
        raise ValueError("Zero items")
    elif len(items) > 1:
        raise ValueError("Multiple items")
    else:
        return next(iter(items))

def find_by_value(list_of_dicts, field, value):
    return [t for t in list_of_dicts if t[field] == value]

def random_string(length):
    chars = string.ascii_letters + string.digits
    return "".join([chars[random.randrange(len(chars))] for _ in range(length)])

# Choose m out of n in random order
def random_choice(m: int, n: int, rnd: Optional[random.Random] = None):
    rnd = rnd if rnd is not None else random
    a = [x for x in range(n)]
    result = []
    for i in range(m):
        idx = rnd.randrange(len(a))
        result.append(a[idx])
        a.pop(idx)
    return result

def project_dir(*paths):
    return os.path.abspath(os.path.join(os.path.normpath(os.path.join(os.path.dirname(__file__), *([os.path.pardir]*5))), *paths))

def yaml_load(s):
    return yaml.load(s, Loader=yaml.SafeLoader)

# TODO Move to sifchain.py
# TODO Refactoring in progress. This should be moved to sifchain.py and only used for gas (float amount) + renamed
def sif_format_amount(amount: Union[int, float], denom: str) -> str:
    return "{}{}".format(amount, denom)

def http_get(url):
    with urllib.request.urlopen(url) as r:
        return r.read()

# Not used yet
def mkcmd(args, env=None, cwd=None, stdin=None):
    result = {"args": args}
    if env is not None:
        result["env"] = env
    if cwd is not None:
        result["cwd"] = cwd
    if stdin is not None:
        result["stdin"] = stdin
    return result

# stdin will always be redirected to the returned process' stdin.
# If pipe, the stdout and stderr will be redirected and available as stdout and stderr of the returned object.
# If not pipe, the stdout and stderr will not be redirected and will inherit sys.stdout and sys.stderr.
def popen(args: Sequence[str], cwd: Optional[str] = None, env: Optional[Mapping[str, str]] = None,
    text: Optional[bool] = None, stdin: Union[str, int, IO, None] = None, stdout: Optional[IO] = None,
    stderr: Optional[IO] = None, disable_log: bool = False
) -> subprocess.Popen:
    if env:
        env = dict_merge(os.environ, env)
    if not disable_log:
        __log.debug(f"popen(): args={repr(args)}, cwd={repr(cwd)}")
    return subprocess.Popen(args, cwd=cwd, env=env, stdin=stdin, stdout=stdout, stderr=stderr, text=text)

def dict_merge(*dicts, override=True):
    result = {}
    for d in dicts:
        for k, v in d.items():
            if override or (k not in result):
                result[k] = v
    return result

def flatten(items: Iterable[Iterable]) -> List:
    return [item for sublist in items for item in sublist]

def format_as_shell_env_vars(env, export=True):
    # TODO escaping/quoting, e.g. shlex.quote(v)
    return ["{}{}=\"{}\"".format("export " if export else "", k, v) for k, v in env.items()]

def disable_noisy_loggers():
    logging.getLogger("eth").setLevel(logging.WARNING)
    logging.getLogger("websockets").setLevel(logging.WARNING)
    logging.getLogger("web3").setLevel(logging.WARNING)
    logging.getLogger("asyncio").setLevel(logging.WARNING)
    logging.getLogger("eth_hash").setLevel(logging.WARNING)

def basic_logging_setup():
    import sys
    # logging.basicConfig(stream=sys.stdout, level=logging.WARNING, format="%(name)s: %(message)s")
    logging.basicConfig(stream=sys.stdout, level=logging.DEBUG, format="%(asctime)s [%(levelname)-8s] %(name)s: %(message)s")
    # logging.getLogger(__name__).setLevel(logging.DEBUG)
    disable_noisy_loggers()

def siftool_logger(name: Optional[str] = None):
    # Shortening "siftool.eth" to "eth" results in a name clash with "eth" dependencies for which we want to disable logging...
    # if name is not None:
    #     name = name[name.rfind(".") + 1:]
    return logging.getLogger(name)

# Recursively transforms template strings containing "${VALUE}". Example:
# >>> template_transform("You are ${what}!", {"what": "${how} late", "how": "very"})
# 'You are very late!'
# Warning: if you use cyclic definitions, this will loop forever.
def template_transform(s, d):
    p = re.compile("^(.*?)(\\${(.*?)})(.*)$")
    while True:
        m = p.match(s)
        if not m:
            return s
        s = s[0:m.start(2)] + d[m[3]] + s[m.end(2):]

def wait_for_enter_key_pressed():
    try:
        input("Press ENTER to exit...")
    except EOFError:
        log = logging.getLogger(__name__)
        log.error("Cannot wait for ENTER keypress since standard input is closed. Instead, this program will now wait "
            "for 100 years and you will have to kill it manually. If you get this message when running in recent "
            "versions of pycharm, enable 'Emulate terminal in output console' in run configuration.")
        time.sleep(3155760000)

def retry(function: Callable, sleep_time: Optional[int] = 5, retries: Optional[int] = 0,
    log: Optional[logging.Logger] = None
) -> Callable:
    def wrapper(*args, **kwargs):
        retries_left = retries
        while True:
            try:
                return function(*args, **kwargs)
            except Exception as e:
                if retries_left == 0:
                    raise e
                if log is not None:
                    log.debug("Retriable exception for {}: args: {}, kwargs: {}, exception: {}".format(repr(function), repr(args), repr(kwargs), repr(e)))
                if sleep_time > 0:
                    time.sleep(sleep_time)
                    retries_left -= 1
                continue
    return wrapper


on_peggy2_branch = not os.path.exists(project_dir("smart-contracts", "truffle-config.js"))

in_github_ci = (os.environ.get("CI") == "true") and os.environ.get("GITHUB_REPOSITORY") and os.environ.get("GITHUB_RUN_ID")

# Make log variable private since it's this module is commonly imported as "*"
__log = siftool_logger(__name__)
