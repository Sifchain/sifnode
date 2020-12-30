# default to running tests

SHELL=/bin/bash -o pipefail

# where to mount the sifnode tree
BASEDIR=/sifnode

.PHONY: all
all: 
	rm -rf data/ ../ganachedb
	-$(MAKE) test
	$(MAKE) pulldata

.depend data: 
	mkdir -p $@

.depend/up: | .depend
	vagrant up
	touch $@

.depend/setup: .depend/up
	vagrant ssh -c "cd ${BASEDIR} && bash test/integration/setup-linux-environment.sh"
	touch $@

.depend/dockerrmi: $(shell find ../../../ -name *.go) | .depend
	vagrant ssh -c "docker rmi -f integration_sifnode1:latest"
	touch $@

.PHONY: integrationenv
integrationenv: .depend/setup .depend/dockerrmi data
	vagrant ssh -c "cd ${BASEDIR} && bash -x test/integration/start-integration-env.sh" 2>&1 | tee data/start-integration-env.txt

.PHONY: test-whitelistedvalidators
test-whitelistedvalidators: integrationenv
	vagrant ssh -c "cd ${BASEDIR} && bash -x test/integration/execute_integration_tests_whitelisted_validators.sh" 2>&1 | tee data/$@.sh

.PHONY: test-main
test-main: integrationenv
	vagrant ssh -c "cd ${BASEDIR} && bash -x test/integration/execute_integration_tests_main.sh" 2>&1 | tee data/$@.sh

.PHONY: test
test: test-main

# convenience targets

up: .depend/up

reload:
	vagrant reload

clean:
	vagrant destroy -f
	rm -rf .depend

pulldata: data
	vagrant ssh -c "bash /sifnode/test/integration/assemble_data.sh"
	tar cfz datafiles.tar.gz data/*
	cp datafiles.tar.gz datafiles.$(shell date +%m-%d-%H-%M-%S).tar.gz
	
cmds:
	@vagrant ssh -c "docker exec integration_sifnode1_1 cat /tmp/testrun.sh"

cmdsrm:
	vagrant ssh -c "docker exec integration_sifnode1_1 bash -c 'rm -f /tmp/testrun.sh'"

ssh:
	vagrant ssh -c "cd /sifnode && . test/integration/vagrantenv.sh && bash -i"

sshi:
	vagrant ssh -c "cd /sifnode && . test/integration/vagrantenv.sh && docker exec -ti integration_sifnode1_1 bash -c 'cd /sifnode && bash -i'"

sshg:
	vagrant ssh -c "cd /sifnode && . test/integration/vagrantenv.sh && docker exec -ti genesis_ganachecli_1 sh"

truffle:
	vagrant ssh -c "cd /sifnode && . test/integration/vagrantenv.sh && cd /sifnode/smart-contracts && truffle console --network develop --verbose-rpc"
