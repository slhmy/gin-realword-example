ifndef COMMIT_HASH
	COMMIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "nohash-$(shell date +%Y%m%d)")
endif

ifndef RUNTIME_ENV
	RUNTIME_ENV := "development"
endif

WEBSITE_API_V1_SWAGGER := api/v1_swagger.json
WEBSITE_API_V1_CLIENT_DIR := website/src/api/v1

.PHONY: gen_clients
gen_clients:
	pnpm install -g @openapitools/openapi-generator-cli
	openapi-generator-cli generate \
        -i /local/${WEBSITE_API_V1_SWAGGER} \
        -g typescript-axios \
        -o /local/${WEBSITE_API_V1_CLIENT_DIR}
	sudo chown -R $(id -u):$(id -g) ${WEBSITE_API_V1_CLIENT_DIR}

.PHONY: gen_docs
gen_docs:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag fmt \
		--dir "internal"
	swag init \
		--parseDependency \
		--exclude ./website \
		--instanceName v1 \
		--dir "./internal/routers/api/v1,./internal/models" \
		--generalInfo "routes.go" \
		--outputTypes json \
		--output api

.PHONY: mod
mod:
	go mod tidy

.PHONY: fmt
fmt:
	go install mvdan.cc/gofumpt@latest
	go install github.com/segmentio/golines@latest
	golines -w -m 120 .
	gofumpt -w .

.PHONY: lint
lint: gen_docs
	@which golangci-lint > /dev/null || (echo "golangci-lint is not installed."; exit 1)
	golangci-lint run \
		--timeout 5m
	cd website; pnpm install; pnpm lint

.PHONY: build
build: mod
	@echo "Building with commit hash: ${COMMIT_HASH}"
	go build -v -ldflags "-X main.lastCommitHash=${COMMIT_HASH}" -o build ./...

.PHONY: build_website
build_website:
	@echo "RUNTIME_ENV: ${RUNTIME_ENV}"
	cd website; pnpm install; pnpm build

.PHONY: dev_website
dev_website:
	@echo "RUNTIME_ENV: ${RUNTIME_ENV}"
	cd website; pnpm install; pnpm dev

.PHONY: web
web: build
	PROJECT_ROOT=${pwd} ./build/web

.PHONY: run
run: build
	${MAKE} -j4 web
