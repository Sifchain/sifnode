mkdir -p configs logs gobin gocache
rm -rf configs/* logs/* gobin/*
python3 src/py/env_framework.py > docker-compose.yml
docker-compose up -d --force-recreate --remove-orphans

runner="docker exec -d -ti smartcontractrunner"
for task in golang_build deploy_contracts print_test_environment
do
  $runner bash -c "cd /sifnode/test/integration && python3 src/py/env_framework.py $task"
done
# $runner bash -c "cd /sifnode/smart-contracts && yarn && cd /sifnode/test/integration && python3 src/py/env_framework.py deploy_contracts"
# $runner bash -c "cd /sifnode/test/integration && python3 src/py/env_framework.py print_test_environment"
