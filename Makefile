.PHONY: default build build-cmd test watch bundle bundle-code-mirror \
	firmware-upload firmware-download firmware-run

-include LocalEnv.mk

include mks/Firmware.mk

default: build

build: | build-cmd

build-cmd:
	mkdir -p dist/
	go build -o dist/rcauto ./cmd/rcauto

bundle: | bundle-code-mirror

bundle-code-mirror:
	cd ./bundles/codemirror && make bundle
	cp ./bundles/codemirror/dist/editor.bundle.js ./ui/assets/js/editor.bundle.js

test:
	go test ./...

watch:
	modd -f modd.conf

fetch-alpine-morph:
	curl -fsSL https://cdn.jsdelivr.net/npm/@alpinejs/morph@3.x.x/dist/cdn.min.js > ui/assets/js/alpine-morph.js
	curl -fsSL https://cdn.jsdelivr.net/npm/alpinejs@3.x.x/dist/cdn.min.js > ui/assets/js/alpine.js
