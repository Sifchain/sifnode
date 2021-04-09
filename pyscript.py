#!/usr/bin/env python
import json
import urllib3
http = urllib3.PoolManager()
import subprocess
print("Starting to Pull Secrets")
result = subprocess.Popen(["kubectl exec --kubeconfig=./kubeconfig -n vault -it vault-0 -- vault kv get -format json kv-v2/us/governance"], stdout=subprocess.PIPE, shell=True)
output,error = result.communicate()
vars_return = json.loads(output.decode('utf-8'))["data"]["data"]
print("Opening temporary secrets file for writing secrets")
temp_secrets = open("tmp_secrets", "w")
for var in vars_return:
    temp_secrets.write('export {key}=\'{values}\' \n'.format(key=var, values=vars_return[var]))
temp_secrets.close()
print("secrets written.")
