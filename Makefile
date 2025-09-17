VAULT_ADDR := https://vault.appier.us
VAULT_KEY_PATH := secret/project/recommendation

DOCKER_DEV_REPO := asia-docker.pkg.dev/appier-docker/docker-ai-rec-asia/rec-vendor-api-dev
DOCKER_TAG := $(DEV_NAME)

CHART_DIR := ./deploy/rec-vendor-api
RELEASE_NAME := rec-vendor-api-dev-$(DEV_NAME)

DEV_CLUSTER := gke_appier-k8s-ai-rec_asia-east1_nelson
DEV_NAMESPACE := rec

REQ_EXECUTABLES := helm kubectl vault consul-template kubectx

.PHONY: all
all: docker-build docker-push


.PHONY: install
install: deploy-dev


.PHONY: clean
clean: delete-dev


.PHONY: pre-commit-check
pre-commit-check: generate test
	golangci-lint run


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
			-template "./config-template/vendors.yaml:$(CHART_DIR)/secrets/vendors.yaml" \
			-template "./config-template/config-dev.yaml:$(CHART_DIR)/secrets/config.yaml"


.PHONY: install-tool
install-tool:
	go install go.uber.org/mock/mockgen@v0.4.0
	brew install golangci-lint


#############  Testing  #############
.PHONY: generate
generate:
	go generate ./...


.PHONY: test
test:
	go test -v -cover -race ./...


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


.PHONY: portforward-dev
portforward-dev:
	kubectx $(DEV_CLUSTER)
	kubectl port-forward svc/$(RELEASE_NAME) 8080:80 -n $(DEV_NAMESPACE)
