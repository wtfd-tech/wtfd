SHELL := bash
.SHELLFLAGS := -euo pipefail -c
.MAKEFLAGS += --warn-undefined-variables --no-builtin-rules
.ONESHELL:

all: .js.sentinel .html.sentinel .static.sentinel
run: .html.sentinel .static.sentinel
	cd js
	yarn run start

.js.sentinel: .build.sentinel .jsdeps.sentinel $(shell find js -type f)
	pushd js
	yarn run build
	popd
	touch .js.sentinel

.jsdeps.sentinel: js/package.json
	pushd js
	yarn
	popd
	touch .jsdeps.sentinel

.html.sentinel: .build.sentinel $(shell find html -type f)
	cp -vr html/. build/
	touch .html.sentinel

.static.sentinel: .build.sentinel $(shell find static -type f)
	mkdir -p build/static
	cp -vr static/. build/static
	touch .static.sentinel

.build.sentinel:
	mkdir -p build
	touch .build.sentinel

.PHONY: clean
clean:
	rm -rf build
	rm .*.sentinel
