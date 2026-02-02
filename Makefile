VAULT_ADDR := https://vault.appier.us
VAULT_KEY_PATH := secret/project/recommendation

DEV_NAME ?= $(shell whoami | sed -e "s/\./-/g")
DOCKER_DEV_REPO := asia-docker.pkg.dev/appier-docker/docker-ai-rec-asia/rec-vendor-api-dev
DOCKER_TAG := $(DEV_NAME)

CHART_DIR := ./deploy/rec-vendor-api
RELEASE_NAME := rec-vendor-api-dev-$(DEV_NAME)

DEV_CLUSTER := gke_appier-k8s-ai-rec_asia-east1_nelson
DEV_NAMESPACE := rec

REQ_EXECUTABLES := helm kubectl vault consul-template kubectx

DOCKER_RF_REPO := asia-docker.pkg.dev/appier-docker/docker-ai-rec-asia/qa/system_test_robot
RF_TAG := v1.0.28
RF_CONFIG := /rec-vendor-api/tests/system_tests/API/res/valueset.dat
RF_TEST_FOLDER := /rec-vendor-api/tests/system_tests/API/testsuite
RF_REPORT_DIR := ./tests/results

.PHONY: all
all: docker-build docker-push


.PHONY: install
install: deploy-dev


.PHONY: clean
clean: delete-dev


.PHONY: validate-vendors-config
validate-vendors-config: config-dev
	go run ./cmd/validate-config $(CHART_DIR)/secrets/vendors.yaml


.PHONY: check-environment
check-environment:
	@if [ -z "$(DEV_NAME)" ]; then \
		echo "Error: DEV_NAME not defined."; \
		exit 1; \
	fi
	$(eval K := $(foreach exec,$(REQ_EXECUTABLES),\
		$(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH, please install it or set PATH variable"))) || true)
	@echo "Passed requirement test"


.PHONY: config-dev
config-dev:
	mkdir -p secrets
	mkdir -p $(CHART_DIR)/secrets
	vault kv get -address $(VAULT_ADDR) --field=private_key secret/project/recommendation/ssh_key/ai-rec-common > secrets/ai-rec-common-key && chmod 600 secrets/ai-rec-common-key

	cp ./config-template/nginx.conf $(CHART_DIR)/secrets/

	consul-template -once -vault-addr $(VAULT_ADDR) \
			-template "./config-template/vendors.yaml:$(CHART_DIR)/secrets/vendors.yaml"
	consul-template -once -vault-addr $(VAULT_ADDR) \
			-template "./config-template/config-dev.yaml:$(CHART_DIR)/secrets/config.yaml"


.PHONY: install-tool
install-tool:
	go install go.uber.org/mock/mockgen@v0.4.0
	go install golang.org/x/tools/cmd/goimports@v0.41.0
	brew install golangci-lint

#############  Local ################
.PHONY: local
local: config-dev
	go run ./cmd/rec-vendor-api/server.go -c $(CHART_DIR)/secrets/config.yaml

# TODO: to be removed with gin retirement
.PHONY: local-grpc
local-grpc: config-dev
	go run ./cmd/rec-vendor-api/server.go -c $(CHART_DIR)/secrets/config.yaml -t grpc

# TODO: to be removed with grpc retirement
.PHONY: local-gin
local-gin: config-dev
	go run ./cmd/rec-vendor-api/server.go -c $(CHART_DIR)/secrets/config.yaml -t gin


#############  Testing  #############
.PHONY: generate
generate:
	go generate ./...


.PHONY: test
test:
	go test -v -cover -race ./...


####################### Pre-Commit Check ##################

.PHONY: fmt
fmt:
	@echo "==> Tidying imports and simplifying format..."
	@go mod tidy
	@goimports -w .
	@gofmt -s -w .


.PHONY: fmt-check
fmt-check:
	@echo "==> Checking if files are formatted and imports are tidied..."
	@if [ -n "$$(gofmt -s -l .)" ] || [ -n "$$(goimports -l .)" ]; then \
		echo "Format check failed! Run 'make fmt' to fix the following files:"; \
		gofmt -s -l .; \
		goimports -l .; \
		exit 1; \
	fi
	@echo "Format check passed"


.PHONY: lint-check
lint-check:
	@echo "==> Running static analysis..."
	@golangci-lint run ./...


.PHONY: pre-commit-check
pre-commit-check: fmt-check lint-check test
	@echo "Success! All checks (Fmt/Lint/Test) passed"


#############  Docker related  #############

.PHONY: docker-build
docker-build: check-environment config-dev
	DOCKER_BUILDKIT=1 docker build . -f ./Dockerfile -t $(DOCKER_DEV_REPO):$(DOCKER_TAG) --ssh ai-rec-common=secrets/ai-rec-common-key


.PHONY: docker-push
docker-push: check-environment
	docker push $(DOCKER_DEV_REPO):$(DOCKER_TAG)


#############  Helm related  #############

.PHONY: deploy-dev
deploy-dev: check-environment config-dev
	kubectx $(DEV_CLUSTER)
	helm upgrade $(RELEASE_NAME) \
		--install  \
		--namespace $(DEV_NAMESPACE) \
		--values ./deploy/rec-vendor-api/values-dev.yaml \
		--set image.tag=$(DOCKER_TAG) \
		--set image.repository=$(DOCKER_DEV_REPO) \
		$(CHART_DIR)
	kubectl rollout restart deployment $(RELEASE_NAME) -n $(DEV_NAMESPACE)


.PHONY: delete-dev
delete-dev: check-environment
	kubectx $(DEV_CLUSTER)
	helm delete $(RELEASE_NAME) --namespace $(DEV_NAMESPACE)


# TODO: change port-forward to the gateway server when gin server is retired
.PHONY: portforward-dev
portforward-dev:
	kubectx $(DEV_CLUSTER)
	kubectl port-forward svc/$(RELEASE_NAME) 8080:80 -n $(DEV_NAMESPACE)


.PHONY: run-e2e
run-e2e:
	rm -f -r ${RF_REPORT_DIR}
	-docker run --rm \
		-it --add-host=host.docker.internal:host-gateway -l ${RF_TAG} -v $(shell pwd):/rec-vendor-api ${DOCKER_RF_REPO}:${RF_TAG} \
		bash -c "cd /rec-vendor-api; mkdir ${RF_REPORT_DIR}; \
		pabot --pabotlib --resourcefile ${RF_CONFIG} --quiet -v ENV:dev -d ${RF_REPORT_DIR}/rec-vendor-api \
		-i RAT ${RF_TEST_FOLDER}/api_rec_vendor_rat.robot"
