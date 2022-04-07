import shutil
import time
from typing import Mapping, List, Union, Optional
from siftool.common import *

ExecArgs = Mapping[str, Union[List[str], str, Mapping[str, str]]]


def buildcmd(args: Optional[str] = None, cwd: Optional[str] = None, env: Optional[Mapping[str, Optional[str]]] = None
) -> ExecArgs:
    return dict((("args", args),) +
        ((("cwd", cwd),) if cwd is not None else ()) +
        ((("env", env),) if env is not None else ())
    )


class Command:
    def execst(self, args, cwd=None, env=None, stdin=None, binary=False, pipe=True, check_exit=True):
        fd_stdout = subprocess.PIPE if pipe else None
        fd_stderr = subprocess.PIPE if pipe else None
        fd_stdin = subprocess.DEVNULL
        if stdin is not None:
            fd_stdin = subprocess.PIPE
            if type(stdin) == list:
                stdin = "".join([line + "\n" for line in stdin])
        proc = popen(args, env=env, cwd=cwd, stdin=fd_stdin, stdout=fd_stdout, stderr=fd_stderr, text=not binary)
        stdout_data, stderr_data = proc.communicate(input=stdin)
        assert pipe == (stdout_data is not None)
        assert pipe == (stderr_data is not None)
        if check_exit and (proc.returncode != 0):
            raise Exception("Command '{}' exited with returncode {}: {}".format(" ".join(args), proc.returncode, repr(stderr_data)))
        return proc.returncode, stdout_data, stderr_data

    # Default implementation of popen for environemnts to start long-lived processes
    def popen(self, args, log_file=None, **kwargs) -> subprocess.Popen:
        stdout = log_file or None
        stderr = log_file or None
        return popen(args, stdout=stdout, stderr=stderr, **kwargs)

    # Starts a process asynchronously (for sifnoded, hardhat, ebrelayer etc.)
    # The arguments should correspond to what buildcmd() returns.
    def spawn_asynchronous_process(self, exec_args: ExecArgs, log_file=None) -> subprocess.Popen:
        return self.popen(**exec_args, log_file=log_file)

    def ls(self, path):
        return os.listdir(path)

    def rm(self, path):
        if os.path.exists(path):
            os.remove(path)

    def mv(self, src, dst):
        os.rename(src, dst)

    def read_text_file(self, path):
        with open(path, "rt") as f:
            return f.read()  # TODO Convert to exec

    def write_text_file(self, path, s):
        with open(path, "wt") as f:
            f.write(s)

    def mkdir(self, path):
        os.makedirs(path, exist_ok=True)

    def rmdir(self, path):
        if os.path.exists(path):
            shutil.rmtree(path)  # TODO Convert to exec

    def rmf(self, path):
        if os.path.exists(path):
            if os.path.isdir(path):
                self.rmdir(path)
            else:
                self.rm(path)

    def copy_file(self, src, dst):
        shutil.copy(src, dst)

    def exists(self, path):
        return os.path.exists(path)

    def is_dir(self, path):
        return os.path.isdir(path) if self.exists(path) else False

    def find_files(self, path, filter=None):
        items = [os.path.join(path, name) for name in self.ls(path)]
        result = []
        for i in items:
            if self.is_dir(i):
                result.extend(self.find_files(i))
            else:
                if (filter is None) or filter(i):
                    result.append(i)
        return result

    def get_user_home(self, *paths):
        return os.path.join(os.environ["HOME"], *paths)

    def mktempdir(self, parent_dir=None):
        args = ["mktemp", "-d"] + (["-p", parent_dir] if parent_dir else [])
        return exactly_one(stdout_lines(self.execst(args)))

    def mktempfile(self, parent_dir=None):
        args = ["mktemp"] + (["-p", parent_dir] if parent_dir else [])
        return exactly_one(stdout_lines(self.execst(args)))

    def chmod(self, path, mode_str, recursive=False):
        args = ["chmod"] + (["-r"] if recursive else []) + [mode_str, path]
        self.execst(args)

    def pwd(self):
        return exactly_one(stdout_lines(self.execst(["pwd"])))

    def __tar_compression_option(self, tarfile):
        filename = os.path.basename(tarfile).lower()
        if filename.endswith(".tar"):
            return ""
        elif filename.endswith(".tar.gz"):
            return "z"
        else:
            raise ValueError(f"Unknown extension for tar file: {tarfile}")

    def tar_create(self, path, tarfile):
        comp = self.__tar_compression_option(tarfile)
        # tar on 9p filesystem reports "file shrank by ... bytes" and exits with errorcode 1
        tar_quirks = False
        if tar_quirks:
            tmpdir = self.mktempdir()
            try:
                shutil.copytree(path, tmpdir, dirs_exist_ok=True)
                self.execst(["tar", "cf" + comp, tarfile, "."], cwd=tmpdir)
            finally:
                self.rmdir(tmpdir)
        else:
            self.execst(["tar", "cf" + comp, tarfile, "."], cwd=path)

    def tar_extract(self, tarfile, path):
        comp = self.__tar_compression_option(tarfile)
        if not self.exists(path):
            self.mkdir(path)
        self.execst(["tar", "xf" + comp, tarfile], cwd=path)

    def wait_for_file(self, path):
        while not self.exists(path):
            time.sleep(1)

    def tcp_probe_connect(self, host, port):
        res = self.execst(["nc", "-z", host, str(port)], check_exit=False)
        return res[0] == 0

    def sha1_of_file(self, path):
        res = self.execst(["sha1sum", "-b", path])
        return stdout_lines(res)[0][:40]

    def download_url(self, url, output_file=None, output_dir=None):
        args = ["curl", "--location", "--silent", "--show-error", url] + \
            (["-O"] if not (output_dir or output_file) else []) + \
            (["-o", output_file] if (output_file and not output_dir) else [])
        self.execst(args, cwd=output_dir)
