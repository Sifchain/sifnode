mkdir -p configs logs gobin gocache
rm -rf configs/* logs/* 
python3 src/py/env_framework.py > docker-compose.yml
docker-compose up -d --force-recreate

runner="docker run -v $(pwd)/../..:/sifnode -v $(pwd)/configs:/configs -v $(pwd)/logs:/logs -v $(pwd)/gobin:/gobin -ti sifdocker:latest"
runner="docker exec -ti smartcontractrunner"
$runner bash -c "cd /sifnode/test/integration && python3 src/py/env_framework.py golang_build"
$runner bash -c "cd /sifnode/smart-contracts && yarn && cd /sifnode/test/integration && python3 src/py/env_framework.py deploy_contracts"
$runner bash -c "cd /sifnode/test/integration && python3 src/py/env_framework.py print_test_environment"

#elif "relayer" in component:
#elif component == "test_environment":
#elif component == "sifnodekeys":  # adds the validator keys to the current test keyring

