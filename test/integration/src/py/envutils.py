import os
import subprocess


def load_env_file(file: str = "../../vagrantenv.sh"):
    sub = subprocess.run(f"bash -c \"source {file} && env | sort\"", shell=True, capture_output=True)
    stdout_string = sub.stdout.decode("utf-8").rstrip()
    lines = stdout_string.split("\n")
    for line in lines:
        elements = line.split("=", 2)
        k = elements[0]
        v = elements[1]
        os.environ[k] = v
