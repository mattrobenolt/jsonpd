GOPATH=$(realpath ../../../../)
GOBIN=$(GOPATH)/bin
GO=GOPATH=$(GOPATH) GOBIN=$(GOBIN) go
APP=jsonpd

OK_COLOR=\033[32;01m
NO_COLOR=\033[0m
BOLD=\033[1m

build:
	@echo "$(OK_COLOR)->$(NO_COLOR) Building $(BOLD)$(APP)$(NO_COLOR)"
	@echo "$(OK_COLOR)==>$(NO_COLOR) Installing dependencies"
	@$(GO) get -v -d ./...
	@echo "$(OK_COLOR)==>$(NO_COLOR) Compiling"
	@$(GO) install -v .

run: build
	@echo "$(OK_COLOR)==>$(NO_COLOR) Running"
	$(GOBIN)/$(APP) -b=:5000 -cb=cb

test:
	$(GO) test -v ./...

bench:
	$(GO) test -run=XXX -bench=. ./...

clean:
	rm -rf $(GOBIN)/*
	rm -rf $(GOPATH)/pkg/*

.PHONY: build run test bench clean
