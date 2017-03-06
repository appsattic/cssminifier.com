all:
	echo 'Provide a target: cssminifier clean'

vendor:
	gb vendor fetch github.com/boltdb/bolt

fmt:
	find src/ -name '*.go' -exec go fmt {} ';'

build: fmt
	gb build all

cssminifier: build
	./bin/cssminifier

test:
	gb test -v

clean:
	rm -rf bin/ pkg/

.PHONY: cssminifier
