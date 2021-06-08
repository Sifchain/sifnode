# ui-stack

Take a sifnode binary and build a docker image that forms the basis of the [`ui-stack`](https://github.com/orgs/Sifchain/packages/container/package/sifnode%2Fui-stack) image that our frontend UI is based on. This is the primary way we keep our repositories in sync.

This code is managed by the frontend team.

## Scripts

| Script                       | Description                                                                                                                  |
| ---------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| `yarn stack`                 | Launch stack from scratch.                                                                                                   |
| `yarn start`                 | Start the stack from the saved snapshot                                                                                      |
| `yarn stack --save-snapshot` | Launch the stack then save the snapshot files                                                                                |
| `yarn stack --push`          | Take snapshot files that have already been saved and wrap them in a docker image and push it to the registry for consumption |

## Why is this set up like this?

We need to have a different deployment cadence for `ui` code compared to our `sifnode` or blockchain code. In order to do this we need to be able to run our frontend e2e and integration tests against background builds.

## Structure

### Docker registry

We tie our repos together using a docker image generated on this repo and consumed by the frontend repo. You can view the latest docker repository uploads here: https://github.com/orgs/Sifchain/packages/container/package/sifnode%2Fui-stack

You can pull the repo using the following although you should never have to as it is built in to the scripts on the ui repo:

```
# raw command for pulling the image
docker pull ghcr.io/sifchain/sifnode/ui-stack:develop
```

Example of running a particular ui-stack built from a commit in the UI repo:

```
yarn stack --tag 1fa110636f94aa2f1c87a30c131462c814293ba8
```

### GH Actions Workflow

We have a github action `./.github/workflows/ui-stack.yml` which uses the scripts located in `./ui` to:

1. Launch a single sifnode, ethereum node, and ebrelayer instance
1. Run migration scripts setting up the state required for our frontend e2e tests
1. Shutdown and save the state of the chains to the `./chains/snapshots` folder
1. Build a docker image to instantly restart the stack on demand
1. Tag the docker image with the current git commit hash and/or a stable tag if this is the tip of a blessed branch (ie.`develop` or `master`)
1. Push that docker container to the ghcr.io registry

### Chains folder

The chains folder holds folders for each individual service our frontend requires to run against locally. Ideally we would place all working data services we use here. Currently we are missing our support data services served from aws. As we increase the amount of blockchains we support they will each add a folder here. Currently we have:

| Folder   | Description                                       |
| -------- | ------------------------------------------------- |
| `/eth`   | our ethereum node we run in the docker container. |
| `/peggy` | ebrelayer                                         |
| `/sif`   | sifnode                                           |

Then we have a folder for snapshot data.

- `/snapshots`

#### Internal structure

If you drill down to each chain folder each one of the chains has a pattern to the way their script files are organized (kind of like an interface):

| Script        | Description                                                                                                                                                     |
| ------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `build.sh`    | Build the chain. This will prepare the chain for execution                                                                                                      |
| `config.sh`   | Provide unique information to other scripts within this folder                                                                                                  |
| `launch.sh`   | This will wipe the state and configure and launch the chain as if for the first time. genesis information will be set here. This calls `./start.sh` at the end. |
| `migrate.sh`  | This gets fired after the chain is running and sets up any state required for our e2e tests.                                                                    |
| `pause.sh`    | This will shutdown the chain.                                                                                                                                   |
| `revert.sh`   | This extracts saved state in `snapshots` for this chain to the appropriate db.                                                                                  |
| `snapshot.sh` | The archives state from the db folder to the snapshot.                                                                                                          |
| `start.sh`    | This starts the chain without configuring or resetting it.                                                                                                      |

if we had written these scripts in javascript each chain would probably be a class and they would all conform to an interface. We could do that but it is not very important right now. PRs welcome!

## Cookbook

Here is a recipe for how to work with the `ui` scripts in `sifnode`.

### I want to change the sifchain genesis information within the ui-stack.

It is rare to have to touch this but if you do, the following outlines how to go about it. You will want to do this because of a need within the e2e branch you will probably have a branch on the ui repo that requires a change to genesis amounts.

1. Firstly see if you can work out a way to avoid making the change here - if you can it would be good to simply do the change in the test you are writing on the frontend. If it takes too long or is needed in many tests or you need the change when developers are testing manually and locally then continue.
1. Create a new branch on the `sifnode` repo branching off `develop`.
1. Now you can make your changes which will probably be in one of the folowing places:

   1. The sifchain genesis code is located here: `./ui/chains/sif/launch.sh` here we specify genesis tokens and amounts. Set up your tokens and amounts in this file.
   1. The CLP pool genesis code is located here: `./ui/chains/sif/migrate.sh` here we specify the default pools that are created and included in the ui-stack image.
   1. Lastly the whitelist and burn limits for ethereum tokens is specified in the post migrate hook here: `./ui/chains/post_migrate.sh`

1. Test your changes

   1. Change these as you require, ensure you have no docker containers running on the ports required and test by running the `yarn stack` command in the `sifnode` (**backend**) repo. This will take a few minutes as it needs to setup and migrate all the chains.
   1. Once it is all setup you can head over to the ui repo and run the frontend end without a stack by runnning `yarn dev` this should connect to the running backend and you can test this works manually.

1. Once you are happy with the new ui-stack we need to save the state and push it to the remote registry.

   1. Firstly quit any of the above processes so we can write the process to disk.
   1. Next save the state to snapshot files by running `yarn stack --save-snapshot` this will take some time as it will setup the chains again and then populate the `./ui/chains/snapshots` folder with new files. Wait grab a coffee/read hackernews. Once it is done it will list a bunch of file names from leveldb and self exit.
   1. Save the snapshot files you just created to a git commit. Take note of the commithash by doing a `git rev-parse HEAD` you will need this commithash to use the stack in the frontend repo.
   1. Let's now push this snapshot to the docker repo using `yarn stack --push` this will wrap those snapshots in a docker image and push it to the docker repository.

1. Now we can run the frontend e2e tests against the repo.
   1. You can test the image you just published by running `yarn e2e --tag xxxxxxxx` where `xxx` is the git commit hash from the backend repo. This will run the e2e tests against the new `ui-stack` image. You may have to make tweaks or iterate on this if you are doing anything complex.

Some caveats and tips

- Every merge to `sifnode/develop` and `sifnode/master` publishes an image based off those branches.
- If you have a change that will break `develop` on the frontend this will cause problems best to ensure things are backwards compatible.
- If you absolutely must to break develop on the frontend besure to have a simultaneous PR that makes the frontend compatible ready to go.
- You can set your FE branch to use an arbitrary tag by using the `yarn stack --set-default-tag xxxx` command. Check in the file it changes to temporarily run your branch against any image tag. You want to ensure you do not merge this file to develop however which should be set to run against `develop`

Well Done! You have managed to update our ui-stack image!

## FAQ

#### Why are you saving state out to zips when you could simply run the chain and save the state when building the docker image?

Initially we were not using docker to create snapshots and when we added docker we thought about removing the zip binary but realized in the future we might want the ability to pass fixture tags for specific e2e tests that require different genesis conditions (note this is not yet built). Therefore it made sense to leave the archiving solution in place and just build around it with docker. In the future should we have special e2e tests that require specific genesis condition we could come back to this.

#### What happens when the backend has breaking changes for the frontend?

In this case the backend will merge their breaking change to develop. Our GHA will run building the stack based on the latest change. After this the sifchain-ui repo will be checked out and its tests applied to the built ui-stack within the github action. If this fails then the backend team know they have broken frontend.

#### How do we know the frontend that is about to be deployed is compatible with master?

Because of the nature of cosmos and blockchains in general we have to be more accommodating than normal on the frontend and need to handle backend versioning ourselves. This means there might be a time when we need to support two different backends based on environment. This is possible when there are breaking changes to one backend on mainnet and a separate API on devnet. We need the frontend to support both. We already run clean architecture and can setup a system to handle swapping out services to manage two different APIs should we require based on runtime configuration. For this reason as we deploy the frontend we test against the devnet, mainnet and testnet branches.
