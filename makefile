# Check to see if we can use ash, in Alpine images, or default to BASH.
SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)


# ==============================================================================
# Define dependencies

GOLANG          := golang:1.22.2
ALPINE          := alpine:3.18
KIND            := kindest/node:v1.27.3
POSTGRES        := postgres:15.4
VAULT           := hashicorp/vault:1.15
GRAFANA         := grafana/grafana:10.1.0
PROMETHEUS      := prom/prometheus:v2.47.0
TEMPO           := grafana/tempo:2.2.0
LOKI            := grafana/loki:2.9.0
PROMTAIL        := grafana/promtail:2.9.0

KIND_CLUSTER    := go-starter-cluster
NAMESPACE       := sales-system
APP             := sales
BASE_IMAGE_NAME := golang/service
SERVICE_NAME    := sales-api
VERSION         := 0.0.1
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME)-metrics:$(VERSION)

# VERSION       := "0.0.1-$(shell git rev-parse --short HEAD)"

# ==============================================================================
# Install dependencies

dev-gotooling:
	go install github.com/divan/expvarmon@latest
	go install github.com/rakyll/hey@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install golang.org/x/tools/cmd/goimports@latest

dev-brew:
	brew update
	brew tap hashicorp/tap
	brew list kind || brew install kind
	brew list kubectl || brew install kubectl
	brew list kustomize || brew install kustomize
	brew list pgcli || brew install pgcli
	brew list watch || brew instal watch
	brew list hey || brew install hey
	brew list vault || brew install vault

dev-docker:
	docker pull $(GOLANG)
	docker pull $(ALPINE)
	docker pull $(KIND)
	docker pull $(POSTGRES)
	docker pull $(VAULT)
	docker pull $(GRAFANA)
	docker pull $(PROMETHEUS)
	docker pull $(TEMPO)
	docker pull $(LOKI)
	docker pull $(PROMTAIL)

# ==============================================================================
# Building containers

all: service metrics

service:
	docker build \
		-f zarf/docker/dockerfile.service \
		-t $(SERVICE_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

metrics:
	docker build \
		-f zarf/docker/dockerfile.metrics \
		-t $(METRICS_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Running from within k8s/kind

dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

	kind load docker-image $(VAULT) --name $(KIND_CLUSTER)
	kind load docker-image $(POSTGRES) --name $(KIND_CLUSTER)
	kind load docker-image $(GRAFANA) --name $(KIND_CLUSTER)
	kind load docker-image $(PROMETHEUS) --name $(KIND_CLUSTER)
	kind load docker-image $(TEMPO) --name $(KIND_CLUSTER)
	kind load docker-image $(LOKI) --name $(KIND_CLUSTER)
	kind load docker-image $(PROMTAIL) --name $(KIND_CLUSTER)

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

# ------------------------------------------------------------------------------

dev-load:
	cd zarf/k8s/dev/sales; kustomize edit set image service-image=$(SERVICE_IMAGE)
	kind load docker-image $(SERVICE_IMAGE) --name $(KIND_CLUSTER)

	cd zarf/k8s/dev/sales; kustomize edit set image metrics-image=$(METRICS_IMAGE)
	kind load docker-image $(METRICS_IMAGE) --name $(KIND_CLUSTER)

dev-apply:
	kustomize build zarf/k8s/dev/vault | kubectl apply -f -

	kustomize build zarf/k8s/dev/database | kubectl apply -f -
	kubectl rollout status --namespace=$(NAMESPACE) --watch --timeout=120s sts/database

	kustomize build zarf/k8s/dev/grafana | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=grafana --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/prometheus | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=prometheus --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/tempo | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=tempo --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/loki | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=loki --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/promtail | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=promtail --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/sales | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(APP) --timeout=120s --for=condition=Ready

dev-restart:
	kubectl rollout restart deployment $(APP) --namespace=$(NAMESPACE)

dev-update: all dev-load dev-restart

dev-update-apply: all dev-load dev-apply

# ------------------------------------------------------------------------------

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run app/tooling/logfmt/main.go -service=$(SERVICE_NAME)

dev-logs-init:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) -f --tail=100 -c init-vault-system
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) -f --tail=100 -c init-vault-loadkeys
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) -f --tail=100 -c init-migrate
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) -f --tail=100 -c init-seed

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-describe:
	kubectl describe nodes
	kubectl describe svc

dev-describe-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) $(APP)

dev-describe-sales:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(APP)

# ------------------------------------------------------------------------------

dev-logs-vault:
	kubectl logs --namespace=$(NAMESPACE) -l app=vault --all-containers=true -f --tail=100

dev-logs-db:
	kubectl logs --namespace=$(NAMESPACE) -l app=database --all-containers=true -f --tail=100

dev-logs-grafana:
	kubectl logs --namespace=$(NAMESPACE) -l app=grafana --all-containers=true -f --tail=100

dev-logs-tempo:
	kubectl logs --namespace=$(NAMESPACE) -l app=tempo --all-containers=true -f --tail=100

dev-logs-loki:
	kubectl logs --namespace=$(NAMESPACE) -l app=loki --all-containers=true -f --tail=100

dev-logs-promtail:
	kubectl logs --namespace=$(NAMESPACE) -l app=promtail --all-containers=true -f --tail=100

# ------------------------------------------------------------------------------

dev-services-delete:
	kustomize build zarf/k8s/dev/sales | kubectl delete -f -
	kustomize build zarf/k8s/dev/grafana | kubectl delete -f -
	kustomize build zarf/k8s/dev/tempo | kubectl delete -f -
	kustomize build zarf/k8s/dev/loki | kubectl delete -f -
	kustomize build zarf/k8s/dev/promtail | kubectl delete -f -
	kustomize build zarf/k8s/dev/database | kubectl delete -f -

dev-describe-replicaset:
	kubectl get rs
	kubectl describe rs --namespace=$(NAMESPACE) -l app=$(APP)

dev-events:
	kubectl get ev --sort-by metadata.creationTimestamp

dev-events-warn:
	kubectl get ev --field-selector type=Warning --sort-by metadata.creationTimestamp

dev-shell:
	kubectl exec --namespace=$(NAMESPACE) -it $(shell kubectl get pods --namespace=$(NAMESPACE) | grep sales | cut -c1-26) --container sales-api -- /bin/sh

dev-database-restart:
	kubectl rollout restart statefulset database --namespace=$(NAMESPACE)

# ==============================================================================
# Administration

migrate:
	go run app/tooling/sales-admin/main.go migrate

seed: migrate
	go run app/tooling/sales-admin/main.go seed

vault:
	go run app/tooling/sales-admin/main.go vault

pgcli:
	pgcli postgresql://postgres:postgres@localhost

liveness:
	curl -il http://localhost:3000/v1/liveness

readiness:
	curl -il http://localhost:3000/v1/readiness

token-gen:
	go run app/tooling/sales-admin/main.go gentoken 5cf37266-3473-4006-984f-9325122678b7 54bb2165-71e1-41a6-af3e-7da4a0e1e2c1

# ==============================================================================
# Metrics and Tracing

metrics-view-sc:
	expvarmon -ports="localhost:4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

metrics-view:
	expvarmon -ports="localhost:3001" -endpoint="/metrics" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

grafana:
	open -a "Google Chrome" http://localhost:3100/

# ==============================================================================
# Running tests within the local computer

test-race:
	CGO_ENABLED=1 go test -race -count=1 ./...

test-only:
	CGO_ENABLED=0 go test -count=1 ./...

lint:
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...

vuln-check:
	govulncheck ./...

test: test-only lint vuln-check

test-race: test-race lint vuln-check

# make docs ARGS="-out json"
# make docs ARGS="-out html"
docs:
	go run app/tooling/docs/main.go --browser $(ARGS)

docs-debug:
	go run app/tooling/docs/main.go $(ARGS)

# ==============================================================================
# Hitting endpoints

token:
	curl -il --user "admin@example.com:gophers" http://localhost:3000/v1/users/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1

# export TOKEN="COPY TOKEN STRING FROM LAST CALL"

users:
	curl -il -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users?page=1&rows=2

user-get:
	curl -il -X GET -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/users/d88b88d9-8bf4-40b3-a700-179424ca40f5


curl-update-user:
	curl -il -X PUT -H "Authorization: Bearer ${TOKEN}" -H 'Content-Type: application/json' -d '{"name":"User Gopher Editado 2"}' http://localhost:3000/v1/users/45b5fbd3-755f-4379-8f07-a58d4a30fa2f

curl-create-product:
	curl -il -X POST -H "Authorization: Bearer ${TOKEN}" -H 'Content-Type: application/json' -d '{"userId":"45b5fbd3-755f-4379-8f07-a58d4a30fa2f","name":"rocambole","cost": 22.80,"quantity": 2}' http://localhost:3000/v1/products

curl-user-get-summary:
	curl -il -X GET -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/usersummary

load:
	hey -m GET -c 100 -n 10000 -H "Authorization: Bearer ${TOKEN}" "http://localhost:3000/v1/users?page=1&rows=2"

otel-test:
	curl -il -H "Traceparent: 00-918dd5ecf264712262b68cf2ef8b5239-896d90f23f69f006-01" --user "admin@example.com:gophers" http://localhost:3000/v1/users/token/54bb2165-71e1-41a6-af3e-7da4a0e1e2c1

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-list:
	go list -m -u -mod=readonly all

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

list:
	go list -mod=mod all

# ==============================================================================
# Class Stuff

run:
	go run app/services/sales-api/main.go | go run app/tooling/logfmt/main.go

run-help:
	go run app/services/sales-api/main.go --help | go run app/tooling/logfmt/main.go

curl:
	curl -il http://localhost:3000/v1/hack

curl-auth:
	curl -il -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/hackauth

load-hack:
	hey -m GET -c 100 -n 100000 "http://localhost:3000/v1/hack"

admin:
	go run app/tooling/sales-admin/main.go

ready:
	curl -il http://localhost:3000/v1/readiness

live:
	curl -il http://localhost:3000/v1/liveness

curl-create:
	curl -il -X POST -H 'Content-Type: application/json' -d '{"name":"diego","email":"diego@testmail.com","roles":["ADMIN"],"department":"IT","password":"123changeme","passwordConfirm":"123changeme"}' http://localhost:3000/v1/users
