# Setup

Create the base docker image:

cd test/integration/vagrant && make sifdocker

Start the containers:

cd test/integration && bash tst.sh

That creates a test/integration/configs directory with all the configuration json.
