import shutil
import time
from common import *


def buildcmd(args, cwd=None, env=None):
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
    def popen(self, args, log_file=None, **kwargs):
        stdout = log_file or None
        stderr = log_file or None
        return popen(args, stdout=stdout, stderr=stderr, **kwargs)

    # Starts a process asynchronously (for sifnoded, hardhat, ebrelayer etc.)
    # The arguments should correspond to what buildcmd() returns.
    def spawn_asynchronous_process(self, exec_args, log_file=None):
        return self.popen(**exec_args, log_file=log_file)

    def rm(self, path):
        if os.path.exists(path):
            os.remove(path)

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

    def get_user_home(self, *paths):
        return os.path.join(os.environ["HOME"], *paths)

    def mktempdir(self):
        return exactly_one(stdout_lines(self.execst(["mktemp", "-d"])))

    def mktempfile(self):
        return exactly_one(stdout_lines(self.execst(["mktemp"])))

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
        tar_quirks = True
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
