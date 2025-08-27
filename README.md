# rec-vendor-api

Table of contents
=================
* [Prerequisite](#prerequisite)
* [Development](#development)
    * [Install modules](#install-modules)
    * [Run pre-commit check](#run-pre-commit-check)
    * [Test on dev cluster](#test-on-dev-cluster)


## Prerequisite

* [go](https://formulae.brew.sh/formula/go)
* [gin-swagger](https://github.com/swaggo/gin-swagger)
* [golangci-lint](https://golangci-lint.run/)
* [gomock](https://github.com/uber-go/mock)

```shell
make install-tool
```

## Development

### Install modules

```shell
# since we use private module, need to setup related env variable
export GOPRIVATE=github.com/plaxieappier

# copy config to ~/.gitconfig
cat .gitconfig >> ~/.gitconfig

# install modules
go mod download
```

### Run pre-commit check

```shell
# run unit test & linter
make pre-commit-check
```

### Test on dev cluster

```shell
# setup dev name
export DEV_NAME=$YOUR_NAME

# Publish docker image & deploy to dev cluster
make all install

# Port-forward api
make portforward-dev

# Open this on browser, and content should show up:
http://localhost:8080/healthz
```
