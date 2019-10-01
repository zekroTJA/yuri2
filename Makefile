### NAMES AND LOCS ############################
APPNAME  = yuri
PACKAGE  = github.com/zekroTJA/yuri2
LDPAKAGE = static
CONFIG   = $(CURDIR)/config/private.config.yml
BINPATH  = $(CURDIR)/bin
###############################################

### EXECUTABLES ###############################
GO     = go
DEP    = dep
GOLINT = golint
GREP   = grep
CLOC   = cloc
###############################################

# ---------------------------------------------

BIN = $(BINPATH)/$(APPNAME)

TAG        = $(shell git describe --tags)
COMMIT     = $(shell git rev-parse HEAD)


ifneq ($(GOOS),)
	BIN := $(BIN)_$(GOOS)
endif

ifneq ($(GOARCH),)
	BIN := $(BIN)_$(GOARCH)
endif

ifneq ($(TAG),)
	BIN := $(BIN)_$(TAG)
endif

ifeq ($(OS),Windows_NT)
	ifeq ($(GOOS),)
		BIN := $(BIN).exe
	endif
endif

ifeq ($(GOOS),windows)
	BIN := $(BIN).exe
endif


PHONY = _make
_make: deps build cleanup

PHONY += build
build: $(BIN) 

PHONY += deps
deps:
	$(DEP) ensure -v

$(BIN):
	$(GO) build  \
		-v -o $@ -ldflags "\
			-X $(PACKAGE)/$(LDPAKAGE).AppVersion=$(TAG) \
			-X $(PACKAGE)/$(LDPAKAGE).AppCommit=$(COMMIT) \
			-X $(PACKAGE)/$(LDPAKAGE).Release=TRUE" \
		$(CURDIR)/cmd/$(APPNAME)/*.go

PHONY += test
test:
	$(GO) test -v -cover ./...

PHONY += lint
lint:
	$(GOLINT) ./... | $(GREP) -v vendor || true

PHONY += run
run:
	$(GO) run -v \
		$(CURDIR)/cmd/$(APPNAME)/*.go -c $(CONFIG)

PHONY += cleanup
cleanup:

PHONY += cloc
cloc:
	$(CLOC) \
		--exclude-dir=vendor,docs,public \
		--exclude-lang=JSON,Markdown,YAML,XML,TOML,Sass ./

PHONY += release
release:
	bash $(CURDIR)/scripts/bundle-release.sh

PHONY += runfe
runfe:
	cd web &&\
		ng serve --port=8081

PHONY += help
help:
	@echo "Available targets:"
	@echo "  #        - creates binary in ./bin"
	@echo "  cleanup  - tidy up temporary stuff created by build or scripts"
	@echo "  cloc     - tidy up temporary stuff created by build or scripts"
	@echo "  deps     - ensure dependencies are installed"
	@echo "  lint     - run linters (golint)"
	@echo "  run      - debug run app (go run) with test config"
	@echo "  test     - run tests (go test)"
	@echo ""
	@echo "Cross Compiling:"
	@echo "  (env GOOS=linux GOARCH=arm make)"
	@echo ""
	@echo "Use different configs for run:"
	@echo "  make CONF=./myCustomConfig.yml run"
	@echo ""


.PHONY: $(PHONY)
