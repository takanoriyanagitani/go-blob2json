#!/bin/sh

tinygo \
	build \
	-o ./blob2json.wasm \
	-target=wasip1 \
	-opt=z \
	-no-debug \
	./cmd/blob2json/main.go
