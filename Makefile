PROJECTNAME=$(shell basename "$(PWD)")
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

build:
	@echo "start build..."
	mkdir -p bin && GOARCH=amd64 GOOS=linux GOBIN=$(GOBIN) go build -o $(GOBIN)/$(PROJECTNAME)

lint:
	bash run_lints.sh