import re
import os
import logging
import subprocess
import string
import random
import yaml
import urllib.request


log = logging.getLogger(__name__)

ANY_ADDR = "0.0.0.0"


def stdout(res):
    return res[1]

def stderr(res):
    return res[2]

def stdout_lines(res):
    return stdout(res).splitlines()

def joinlines(lines):
    return "".join([x + os.linesep for x in lines])

def zero_or_one(items):
    if len(items) == 0:
        return None
    elif len(items) > 1:
        raise ValueError("Multiple items")
    else:
        return items[0]

def exactly_one(items):
    if len(items) == 0:
        raise ValueError("Zero items")
    elif len(items) > 1:
        raise ValueError("Multiple items")
    else:
        return items[0]

def find_by_value(list_of_dicts, field, value):
    return [t for t in list_of_dicts if t[field] == value]

def random_string(length):
    chars = string.ascii_letters + string.digits
    return "".join([chars[random.randrange(len(chars))] for _ in range(length)])

def project_dir(*paths):
    return os.path.abspath(os.path.join(os.path.normpath(os.path.join(os.path.dirname(__file__), *([os.path.pardir]*3))), *paths))

def yaml_load(s):
    return yaml.load(s, Loader=yaml.SafeLoader)

# TODO Move to sifchain.py
def sif_format_amount(amount, denom):
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
def popen(args, cwd=None, env=None, text=None, stdin=None, stdout=None, stderr=None):
    if env:
        env = dict_merge(os.environ, env)
    logging.debug(f"popen(): args={repr(args)}, cwd={repr(cwd)}")
    return subprocess.Popen(args, cwd=cwd, env=env, stdin=stdin, stdout=stdout, stderr=stderr, text=text)

def dict_merge(*dicts, override=True):
    result = {}
    for d in dicts:
        for k, v in d.items():
            if override or (k not in result):
                result[k] = v
    return result

def format_as_shell_env_vars(env, export=True):
    # TODO escaping/quoting, e.g. shlex.quote(v)
    return ["{}{}=\"{}\"".format("export " if export else "", k, v) for k, v in env.items()]

def basic_logging_setup():
    import sys
    # logging.basicConfig(stream=sys.stdout, level=logging.WARNING, format="%(name)s: %(message)s")
    logging.basicConfig(stream=sys.stdout, level=logging.DEBUG, format="%(asctime)s [%(levelname)-8s] %(name)s: %(message)s")
    # logging.getLogger(__name__).setLevel(logging.DEBUG)
    logging.getLogger("eth").setLevel(logging.WARNING)
    logging.getLogger("websockets").setLevel(logging.WARNING)
    logging.getLogger("web3").setLevel(logging.WARNING)

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


on_peggy2_branch = not os.path.exists(project_dir("smart-contracts", "truffle-config.js"))
