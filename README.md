# rec-vendor-api

Table of contents
=================
- [rec-vendor-api](#rec-vendor-api)
- [Table of contents](#table-of-contents)
  - [Prerequisite](#prerequisite)
  - [Development](#development)
    - [Install modules](#install-modules)
    - [Run pre-commit check](#run-pre-commit-check)
    - [Test on dev cluster](#test-on-dev-cluster)
  - [Configuration](#configuration)
    - [TS Team Vendor Configuration Guide](#ts-team-vendor-configuration-guide)
  - [Requester Strategy and Tracker Strategy](#requester-strategy-and-tracker-strategy)
    - [Supported URL Macros](#supported-url-macros)


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
# Publish docker image & deploy to dev cluster
make all install

# Port-forward api
make portforward-dev

# Open this on browser, and content should show up:
http://localhost:8080/healthz

# Manual test to verify the correct image size:
./scripts/manual_test_all.sh
```

We provide the QA E2E Robotframework tool in the Makefile for DEV pre-testing purposes.
Make sure your DEV vendor-api service has been port-forwarded at port **_8080_**, and run:

```shell
make run-e2e
```


## Configuration

### TS Team Vendor Configuration Guide

For Technical Support (TS) team members who need to modify vendor configurations, please refer to our comprehensive guide:

**[Vendor Configuration Guide for TS Team](https://appier.atlassian.net/wiki/spaces/AI/pages/4584833092/Vendor+Configuration+Guide+for+TS+Team)**

## Requester Strategy and Tracker Strategy

We use macros (placeholders) in our URL templates for dynamic replacement. At runtime, these macros get swapped out for real data, making the request API and tracking URLs dynamic and easy to maintain.

### Supported URL Macros

| Macro                  | Description                                      | Example Replacement                                |
| ---------------------- | ------------------------------------------------ | -------------------------------------------------- |
| `{width}`              | Image width (integer)                            | `1200`                                             |
| `{height}`             | Image height (integer)                           | `600`                                              |
| `{user_id_lower}`      | User ID in lowercase                             | `57846b41-0290-40c5-9e96-88d17f59eac5`             |
| `{user_id_case_by_os}` | User ID in lowercase (aos) or in uppercase (ios) | `57846b41-0290-40c5-9e96-88d17f59eac5`             |
| `{click_id_base64}`    | Click ID encoded in base64                       | `Y2xpY2tJRA`                                       |
| `{web_host}`           | Site domain for web or empty for app             | `testabc.com`                                      |
| `{bundle_id}`          | App bundle ID or empty for web                   | `com.coupang.mobile`                               |
| `{adtype}`             | The impression ad type                           | `value=2(banner) and 3(native)`                    |
| `{partner_id}`         | Partner ID                                       | `kakao_kr`                                         |
| `{subid}`              | Sub ID for coupang partners                      | `650alldb2`                                        |
| `{keeta_campaign_id}`  | Campaign ID of Keeta                             | `1901910420462051330`                              |
| `{click_id}`           | Raw click ID                                     | `oSRKfG7nRAy0wgPAg3gN8`                            |
| `{client_ip}`          | User's IP address                                | `182.239.90.0`                                     |
| `{latitude}`           | User's geo latitude                              | `22.3200`                                          |
| `{longitude}`          | User's geo longitude                             | `114.1800`                                         |
| `{product_url}`        | Product URL string                               | `https://ads-partners.example.com/image2/uuid1234` |
