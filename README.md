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

## Requester Strategy and Tracker Strategy

We use macros (placeholders) in our URL templates for dynamic replacement. At runtime, these macros get swapped out for real data, making the request API and tracking URLs dynamic and easy to maintain.

### Supported Request URL Macros

| Macro             | Description                          | Example Replacement                    |
|-------------------|--------------------------------------|----------------------------------------|
| `{width}`         | Image width (integer)                | `1200`                                 |
| `{height}`        | Image height (integer)               | `600`                                  |
| `{user_id_lower}` | User ID in lowercase                 | `57846b41-0290-40c5-9e96-88d17f59eac5` |

### Supported Tracking URL Macros

| Macro                   | Description                              | Example Replacement                                          |
|-------------------------|------------------------------------------|--------------------------------------------------------------|
| `{product_url}`         | Product URL string                       | `https://ads-partners.example.com/image2/uuid1234`           |
| `{encoded_product_url}` | Encoded Product URL (URL-encoded)        | `https%3A%2F%2Fads-partners.example.com%2Fimage2%2Fuuid1234` |
| `{click_id_base64}`     | Click ID encoded in base64               | `Y2xpY2tJRA`                                                 |
| `{user_id_lower}`       | User ID in lowercase                     | `57846b41-0290-40c5-9e96-88d17f59eac5`                       |
