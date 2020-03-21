VERSION = v0.1
REPO = desqueeze
OWNER = etsxxx
BIN = $(REPO)
CURRENT_REVISION ?= $(shell git rev-parse --short HEAD)
LDFLAGS = -w -s -X 'main.version=$(VERSION)' -X 'main.gitcommit=$(CURRENT_REVISION)'

all: clean test build

test:
	go test ./...

build:
	go build -ldflags="$(LDFLAGS)" -trimpath -o bin/$(BIN) ./

cross: clean
	goxc -build-gcflags="-trimpath=$(shell pwd)" -build-ldflags="$(LDFLAGS)" -d=./dist

deploy: cross
	ghr -u $(OWNER) -r $(REPO) $(VERSION) ./dist/snapshot

clean:
	rm -rf bin dist

.PHONY: test build cross deploy clean
