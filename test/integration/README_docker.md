# Setup

Create the base docker image:

```
cd test/integration/vagrant && make sifdocker
```

Start the containers:

```
cd test/integration && bash tst.sh
```

That creates a test/integration/configs directory with all the configuration json.

# Testing

You can keep rerunning the tests by using the last line of the tst.sh file:

```
time docker exec -ti testrunner bash -x smalltest.sh
```
