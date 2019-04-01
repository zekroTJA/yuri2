GO     = go
DEP    = dep
GOLINT = golint
GREP   = grep

TESTCONF = $(CURDIR)/config/private.config.yml

BINNAME = shinpuru
BIN = $(CURDIR)/bin/yuri_server

ifeq ($(OS),Windows_NT)
	BIN := $(BIN).exe
endif

PACKAGE = github.com/zekroTJA/yuri2/

TAG    = $(shell git describe --tags)
COMMIT = $(shell git rev-parse HEAD)

$(BIN): deps
	@echo TEST $@
	$(GO) build \
		-v -o $@ -ldflags "\
			-X $(PACKAGE)/static.AppVersion=$(TAG) \
			-X $(PACKAGE)/static.AppCommit=$(COMMIT) \
			-X $(PACKAGE)/static.Release=TRUE" \
		$(CURDIR)/cmd/yuri

PHONY = deps
deps:
	$(DEP) ensure -v

PHONY += run
run:
	$(GO) run -v $(CURDIR)/cmd/yuri \
		-c $(TESTCONF)

PHONY += lint
lint:
	$(GOLINT) ./... | $(GREP) -v vendor

.PHONY: $(PHONY)