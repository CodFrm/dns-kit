
check-cago:
ifneq ($(which cago),)
	go get -u github.com/codfrm/cago
endif

check-mockgen:
ifneq ($(which mockgen),)
	go get -u github.com/golang/mock/mockgen
endif

check-golangci-lint:
ifneq ($(which golangci-lint),)
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
endif

check-goconvey:
ifneq ($(which goconvey),)
	go get -u github.com/smartystreets/goconvey
endif

swagger: check-cago
	cago gen swag

lint: check-golangci-lint
	golangci-lint run

lint-fix: check-golangci-lint
	golangci-lint run --fix

test: lint
	go test -v ./...

coverage.out cover:
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

html-cover: coverage.out
	go tool cover -html=coverage.out
	go tool cover -func=coverage.out

generate: check-mockgen swagger
	go generate ./... -x

goconvey: check-goconvey
	goconvey

GOOS=linux
GOARCH=amd64
APP_NAME=dns-kit
APP_VERSION=1.0.0
LD_FLAGS=-w -s -X github.com/codfrm/cago/configs.Version=${APP_VERSION}

SUFFIX=
ifeq ($(GOOS),windows)
	SUFFIX=.exe
endif

build:
	# 构建前端
	cd frontend && yarn && yarn build
	# 构建后端
	go mod tidy
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o "bin/${APP_NAME}${BINARY_SUFFIX}" -trimpath -ldflags "${LD_FLAGS}" ./cmd/app

docker:
	docker build --platform=$(GOOS)/$(GOARCH) -t $(APP_NAME):$(APP_VERSION) .

docker-push:
	docker tag $(APP_NAME):$(APP_VERSION) codfrm/$(APP_NAME):$(APP_VERSION)
	docker push codfrm/$(APP_NAME):$(APP_VERSION)
	docker tag $(APP_NAME):$(APP_VERSION) codfrm/$(APP_NAME):latest
	docker push codfrm/$(APP_NAME):latest
