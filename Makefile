PREFIX			  ?= registry.cn-hangzhou.aliyuncs.com/rdc-incubator
TAG				  ?= $(shell date +%s)
SHADOW_IMAGE	  =  kt-connect-shadow
SHADOW_BASE_IMAGE =  shadow-base
ROUTER_IMAGE	  =  kt-connect-router
NAVIGATOR_IMAGE	  =  kt-connect-navigator

# run mod tidy
mod:
	go mod tidy -compat=1.17

# run unit test
test:
	mkdir -p artifacts/report/coverage
	go test -v -cover -coverprofile c.out.tmp ./...
	cat c.out.tmp | grep -v "_mock.go" > c.out
	go tool cover -html=c.out -o artifacts/report/coverage/index.html

# build kt project
compile:
	goreleaser --snapshot --skip-publish --rm-dist

# check the style
check:
	go vet ./pkg/... ./cmd/...

# build ktctl
ktctl:
	GOARCH=amd64 GOOS=linux go build -o artifacts/ktctl/ktctl-linux ./cmd/ktctl
	GOARCH=amd64 GOOS=darwin go build -o artifacts/ktctl/ktctl-darwin ./cmd/ktctl
	GOARCH=amd64 GOOS=windows go build -o artifacts/ktctl/ktctl-windows ./cmd/ktctl

# build this image before shadow
shadow-base:
	docker build -t $(PREFIX)/$(SHADOW_BASE_IMAGE):$(TAG) -f build/docker/shadow/Dockerfile_base .

# build shadow image
shadow:
	GOARCH=amd64 GOOS=linux go build -gcflags "all=-N -l" -o artifacts/shadow/shadow-linux-amd64 cmd/shadow/main.go
	docker build -t $(PREFIX)/$(SHADOW_IMAGE):$(TAG) -f build/docker/shadow/Dockerfile .

# shadow with dlv
shadow-dlv:
	make shadow TAG=latest
	scripts/build-shadow-dlv

# shadow for local debug
shadow-local:
	go build -gcflags "all=-N -l" -o artifacts/shadow/shadow-local cmd/shadow/main.go

# build router image
router:
	GOARCH=amd64 GOOS=linux go build -gcflags "all=-N -l" -o artifacts/router/router-linux-amd64 cmd/router/main.go
	docker build -t $(PREFIX)/$(ROUTER_IMAGE):$(TAG) -f build/docker/router/Dockerfile .

# build this image before navigator
navigator-base:
	docker build -t $(PREFIX)/$(NAVIGATOR_BASE_IMAGE):$(TAG) -f build/docker/navigator/Dockerfile_base .

# build navigator image
navigator:
	GOARCH=amd64 GOOS=linux go build -gcflags "all=-N -l" -o artifacts/navigator/navigator-linux-amd64 cmd/navigator/main.go
	docker build -t $(PREFIX)/$(NAVIGATOR_IMAGE):$(TAG) -f build/docker/navigator/Dockerfile .

# navigator for local debug
navigator-local:
	go build -gcflags "all=-N -l" -o artifacts/navigator/navigator cmd/navigator/main.go

# clean up workspace
clean:
	rm -fr artifacts dist
