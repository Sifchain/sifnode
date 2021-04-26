mkdir -p configs logs gobin gocache
rm -rf configs/* logs/*
python3 src/py/env_framework.py > docker-compose.yml
docker-compose up -d --force-recreate --remove-orphans

runner="docker exec -ti smartcontractrunner"
for task in golang_build deploy_contracts print_test_environment
do
  $runner bash -c "cd /sifnode/test/integration && python3 src/py/env_framework.py $task"
done
