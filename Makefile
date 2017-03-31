all:
	echo 'Provide a target: cssminifier clean'

vendor:
	gb vendor fetch github.com/boltdb/bolt

fmt:
	find src/ -name '*.go' -exec go fmt {} ';'

minify:
	curl -X POST -s --data-urlencode 'input@static/s/js/ready.js' https://javascript-minifier.com/raw > static/s/js/ready.min.js
	curl -X POST -s --data-urlencode 'input@static/s/css/style.css' https://cssminifier.com/raw > static/s/css/style.min.css

test:
	curl -X POST -s --data-urlencode 'input@test/calc.css'          http://localhost:9804/raw > test/calc.min.css
	curl -X POST -s --data-urlencode 'input@test/import.css'        http://localhost:9804/raw > test/import.min.css
	curl -X POST -s --data-urlencode 'input@test/ok.css'            http://localhost:9804/raw > test/ok.min.css
	curl -X POST -s --data-urlencode 'input@test/bootstrap.css'     http://localhost:9804/raw > test/bootstrap.min.css
	curl -X POST -s --data-urlencode 'input@test/causes-error.css'  http://localhost:9804/raw > test/causes-error.min.css
	curl -X POST -s --data-urlencode 'input@test/infinite-loop.css' http://localhost:9804/raw > test/infinite-loop.min.css

test-remote:
	curl -X POST -s --data-urlencode 'input@test/calc.css'          http://cssminifier.com/raw > test/calc.min.css
	curl -X POST -s --data-urlencode 'input@test/import.css'        http://cssminifier.com/raw > test/import.min.css
	curl -X POST -s --data-urlencode 'input@test/ok.css'            http://cssminifier.com/raw > test/ok.min.css
	curl -X POST -s --data-urlencode 'input@test/bootstrap.css'     http://cssminifier.com/raw > test/bootstrap.min.css
	curl -X POST -s --data-urlencode 'input@test/causes-error.css'  http://cssminifier.com/raw > test/causes-error.min.css
	curl -X POST -s --data-urlencode 'input@test/infinite-loop.css' http://cssminifier.com/raw > test/infinite-loop.min.css

build: fmt
	gb build all

cssminifier: build
	./bin/cssminifier

test:
	gb test -v

clean:
	rm -rf bin/ pkg/

.PHONY: cssminifier
