# Minecarft Manager

## Build tags

`js` build tag is used to exclude any dependency on `gameserver` package & related docker client SDK, when build wasm package.

`windows` build tag is used to have alternate way of initializing docker client when debugging on windows OS

`unix` build tag is used for linux compilation where we need access to `/var/run/docker.sock` socket.

## Debugging commands

To support dynamic build, to be as close as HMR

```pwsh
wgo -xdir tmp -file .go -file .js -file .css  make run
```

## Docker image build commands

Image build command

```pwsh
docker build . --tag apogee-dev/mc-manager:local -f Dockerfile
```

## Docker test commands

Test image build command

```pwsh
docker build . --tag apogee-dev/mc-manager-tests:local -f Dockerfile.test
```

Test run command

```pwsh
docker run --volume=/var/run/docker.sock:/var/run/docker.sock apogee-dev/mc-manager-tests:local
```
