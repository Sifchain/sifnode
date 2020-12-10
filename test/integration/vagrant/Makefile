# default to running tests

# where to mount the sifnode tree
BASEDIR=/sifnode

all: test

.depend: 
	mkdir -p $@

.depend/up: | .depend
	vagrant up
	touch $@

.depend/setup: .depend/up
	vagrant ssh -c "cd ${BASEDIR} && bash test/integration/setup-linux-environment.sh"
	touch $@

.PHONY: test
test: .depend/setup
	vagrant ssh -c "cd ${BASEDIR} && bash -x test/integration/start-integration-env.sh"

# convenience targets

up: .depend/up

reload:
	vagrant reload

clean:
	vagrant destroy -f
	rm -rf .depend

cmds:
	vagrant ssh -c "docker exec integration_sifnode1_1 cat /tmp/testrun.sh"
