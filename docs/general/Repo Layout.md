# Repo Layout

Proposed repo layout as follows;

```
./api -> Public API code. ie graphQL
./app -> Cosmos APP generated code
./build -> Terraform / Kubes / Dockerimage build files
./cmd -> golang std cmd binary entry points
./docs -> Project implantation documentation
./test -> E2E blackbox integration test suites
./ui -> vuejs frontend
./x -> Cosmos sifchain modules
./Makefile -> Make cmds for driving everything from dev -> test -> build -> deploy
./config.yml -> Cosmos/Starport config file
./docker-compose.yml -> local development environment setup for all components
```
