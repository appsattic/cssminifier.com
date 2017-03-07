all:
	echo 'Provide a target: cssminifier clean'

vendor:
	gb vendor fetch github.com/boltdb/bolt

fmt:
	find src/ -name '*.go' -exec go fmt {} ';'

assets:
	curl -X POST -s --data-urlencode 'input@static/s/js/ready.js' https://javascript-minifier.com/raw > static/s/js/ready.min.js
	curl -X POST -s --data-urlencode 'input@static/s/css/style.css' https://cssminifier.com/raw > static/s/css/style.min.css

build: fmt
	gb build all

cssminifier: build
	./bin/cssminifier

test:
	gb test -v

clean:
	rm -rf bin/ pkg/

.PHONY: cssminifier
