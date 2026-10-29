#!/usr/bin/env python3
"""Scaffold Terraform clusters and deploy sifnode to AWS EKS or Azure AKS."""

import argparse
import ipaddress
import json
import os
from pathlib import Path
import re
import subprocess
import sys


DEPLOY = Path(__file__).resolve().parent
ROOT = DEPLOY.parent


def run(command, **kwargs):
    return subprocess.run(command, check=True, **kwargs)


def name(value):
    if not re.fullmatch(r"[a-z][a-z0-9-]{1,29}[a-z0-9]", value):
        raise argparse.ArgumentTypeError("use 3-31 lowercase letters, digits or hyphens")
    return value


def cidr(value):
    try:
        network = ipaddress.ip_network(value, strict=True)
    except ValueError as error:
        raise argparse.ArgumentTypeError(str(error)) from error
    if network.version != 4 or network.prefixlen == 0:
        raise argparse.ArgumentTypeError("use an IPv4 administrator CIDR narrower than /0")
    return str(network)


def private_write(path, text):
    # O_EXCL also prevents following an existing symlink when writing kubeconfig.
    fd = os.open(path, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600)
    with os.fdopen(fd, "w") as stream:
        stream.write(text)


def scaffold(args, directory):
    if not args.admin_cidr:
        raise ValueError("scaffold requires at least one --admin-cidr for the Kubernetes API")
    directory.mkdir(parents=True, mode=0o700, exist_ok=False)
    provider = args.provider
    template = (DEPLOY / "terraform" / "template" / provider / "cluster.tf.tpl").read_text()
    module_path = DEPLOY / "terraform" / "providers" / provider
    private_write(directory / "main.tf", template.replace("__MODULE_SOURCE__", json.dumps(str(module_path))))
    variables = {
        "cluster_name": f"sifchain-{provider}-{args.name}",
        "api_server_authorized_ip_ranges": args.admin_cidr,
        "kubernetes_version": args.kubernetes_version,
    }
    variables["location" if provider == "azure" else "region"] = args.location
    private_write(directory / "terraform.tfvars.json", json.dumps(variables, indent=2) + "\n")
    private_write(directory / "environment.json", json.dumps({"provider": provider, "name": args.name}) + "\n")
    print(f"Scaffolded {provider} in {directory}. Review terraform.tfvars.json, then run init and plan.")


def check_environment(args, directory):
    environment = json.loads((directory / "environment.json").read_text())
    if environment != {"provider": args.provider, "name": args.name}:
        raise ValueError("provider/name do not match this scaffold; use the original environment")


def kubeconfig(directory):
    destination = directory / "kubeconfig"
    result = run(["terraform", f"-chdir={directory}", "output", "-raw", "kubectl_config"],
                 stdout=subprocess.PIPE, text=True)
    if not result.stdout.strip():
        raise ValueError("Terraform returned an empty kubeconfig; apply the cluster first")
    # Replace atomically with permissions fixed before writing. Never print the config.
    temporary = directory / "kubeconfig.tmp"
    private_write(temporary, result.stdout)
    os.replace(temporary, destination)
    print(f"Wrote {destination} (mode 0600)")


def parser():
    result = argparse.ArgumentParser(description=__doc__)
    result.add_argument("command", choices=["scaffold", "init", "validate", "plan", "apply",
                                            "kubeconfig", "node", "uninstall", "destroy"])
    result.add_argument("--provider", choices=["aws", "azure"], required=True)
    result.add_argument("--name", type=name, required=True, help="environment name, e.g. testnet")
    result.add_argument("--directory", type=Path, help="override the default .live environment path")
    result.add_argument("--location", help="Azure location or AWS region (scaffold only)")
    result.add_argument("--admin-cidr", action="append", type=cidr,
                        help="administrator IPv4 CIDR, repeatable (scaffold only)")
    result.add_argument("--kubernetes-version", help="supported version; omitted uses cloud default")
    result.add_argument("--values", type=Path, help="Helm values file (required for node)")
    result.add_argument("--namespace", type=name, default="sifnode")
    result.add_argument("--release", type=name, default="sifnode")
    return result


def main(argv=None):
    args = parser().parse_args(argv)
    directory = (args.directory or ROOT / ".live" / f"sifchain-{args.provider}-{args.name}").resolve()
    args.location = args.location or ("eastus" if args.provider == "azure" else "us-west-2")
    try:
        if args.command == "scaffold":
            scaffold(args, directory)
            return 0
        check_environment(args, directory)
        if args.command in ("init", "validate", "plan", "apply", "destroy"):
            # Keep Terraform's interactive confirmation for apply/destroy.
            run(["terraform", f"-chdir={directory}", args.command])
        elif args.command == "kubeconfig":
            kubeconfig(directory)
        else:
            config = directory / "kubeconfig"
            if not config.is_file():
                raise ValueError("run kubeconfig first; this command uses only the environment's cluster")
            command = ["helm", "--kubeconfig", str(config)]
            if args.command == "node":
                if not args.values or not args.values.is_file():
                    raise ValueError("node requires --values pointing to a Helm values file")
                command += ["upgrade", "--install", args.release, str(DEPLOY / "helm" / "sifnode"),
                            "--namespace", args.namespace, "--create-namespace", "--values", str(args.values),
                            # Provider is set last so a copied values file cannot select the wrong cloud.
                            "--set-string", f"cloudProvider={args.provider}", "--wait", "--timeout", "15m"]
            else:
                command += ["uninstall", args.release, "--namespace", args.namespace, "--wait"]
            run(command)
        return 0
    except (OSError, ValueError, subprocess.CalledProcessError) as error:
        print(f"error: {error}", file=sys.stderr)
        return 1


if __name__ == "__main__":
    sys.exit(main())
