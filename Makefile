# TODO: Remove .ONESHELL

SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX = >

# Default - top level rule is what gets ran when you run just `make`
build: wtfd
.PHONY: build

run: tmp/.js.sentinel
> rm -f tmp/.gen.sentinel
> pushd internal
> packr2 clean
> popd
> go run ./cmd/wtfd.go
.PHONY: run

test: tmp/.test.sentinel
.PHONY: test

clean: tmp/.check-deps.sentinel
> rm -rf tmp
> rm -f wtfd
> rm -rf html/static/*.js html/static/*.js.map
> pushd internal
> packr2 clean
> popd
.PHONY: clean

js-run: tmp/.js-deps.sentinel
> pushd frontend
> yarn run start
> popd
.PHONY: js-run

wtfd: $(shell find . -name '*.go') tmp/.gen.sentinel $(shell find html -name '*.html')
> go build ./cmd/wtfd.go


.PHONY: fmt
fmt:
> find -type f -name "*.go" -exec go fmt {} \;


tmp/.test.sentinel: $(shell find . -name '*.go') tmp/.check-deps.sentinel
> mkdir -p $(@D)
> go test ./...
> touch $@

# go generate
tmp/.gen.sentinel: tmp/.js.sentinel $(shell find html/static/ -type f | grep -v '.js$$' | grep -v '.map$$')
> mkdir -p $(@D)
> go generate ./...
> touch $@

tmp/.js.sentinel: $(shell find frontend/src -type f) tmp/.js-deps.sentinel $(shell find  frontend/ -name 'webpack.*.js')
> mkdir -p $(@D)
> pushd frontend
> yarn run build
> popd
> rm -f tmp/.gen.sentinel
> pushd internal
> packr2 clean
> popd
> touch $@

tmp/.js-deps.sentinel: tmp/.check-deps.sentinel
> mkdir -p $(@D)
> pushd frontend
> yarn
> popd
> touch $@

tmp/.check-deps.sentinel:
> mkdir -p $(@D)
> which go || (echo "go not installed"; exit 1)
> which packr2 || (echo "packr2 not installed, run 'go get -u github.com/gobuffalo/packr/v2/packr2' please"; exit 1)
> which yarn || (echo "yarn not installed"; exit 1)
> touch $@
