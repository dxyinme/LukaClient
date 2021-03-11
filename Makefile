
PROJECT=LukaClient
GOPATH ?= $(shell go env GOPATH)
CURDIR := $(shell pwd)

FAIL_ON_STDOUT := awk '{ print } END { if (NR > 0) { exit 1 } }'

GO        := GO111MODULE=on go
GOBUILD   := $(GO) build
GOTEST    := $(GO) test -p 4

FILES     := $$(find . -name "*.go")

default: cli

cli:
	@echo "generate cli"
	$(GOBUILD) -o bin/client_cli main/cli.go

benchmark:
	@echo "generate benchmark"
	$(GOBUILD) -o bin/client_benchmark main/benchmark.go
	@cp -r test/ bin/

cli-core:
	@echo "generate cli-core"
	$(GOBUILD) -o bin/cli_core main/cli-core.go
	@cp -r ClientExample/ bin/

gui:
	@echo "generate gui-client"
	$(GOBUILD) -o bin/gui main/gui.go
	@cp -r ClientExample/ bin/
	@mkdir -p bin/privateKey/
	@mkdir -p bin/SaveTmp/

fmt:
	@echo "gofmt (simplify)"
	@gofmt -s -l -w $(FILES) 2>&1 | $(FAIL_ON_STDOUT)
