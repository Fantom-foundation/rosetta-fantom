.PHONY: deps build run lint run-mainnet-online run-mainnet-offline run-testnet-online \
	run-testnet-offline check-comments add-license check-license shorten-lines \
	spellcheck salus build-local format check-format update-tracer test coverage coverage-local \
	update-bootstrap-balances mocks

ADDLICENSE_INSTALL=go install github.com/google/addlicense@latest
ADDLICENSE_CMD=addlicense
ADDLICENCE_SCRIPT=${ADDLICENSE_CMD} -c "Fantom Foundation, Coinbase, Inc." -l "apache" -v
SPELLCHECK_CMD=go run github.com/client9/misspell/cmd/misspell
GOLINES_INSTALL=go install github.com/segmentio/golines@latest
GOLINES_CMD=golines
GOLINT_INSTALL=go get golang.org/x/lint/golint
GOLINT_CMD=golint
GOVERALLS_INSTALL=go install github.com/mattn/goveralls@latest
GOVERALLS_CMD=goveralls
GOIMPORTS_CMD=go run golang.org/x/tools/cmd/goimports
GO_PACKAGES=./services/... ./cmd/... ./configuration/... ./opera/...
GO_FOLDERS=$(shell echo ${GO_PACKAGES} | sed -e "s/\.\///g" | sed -e "s/\/\.\.\.//g")
TEST_SCRIPT=go test ${GO_PACKAGES}
LINT_SETTINGS=golint,misspell,gocyclo,gocritic,whitespace,goconst,gocognit,bodyclose,unconvert,lll,unparam
PWD=$(shell pwd)
NOFILE=100000

MAINNET_GENESIS=mainnet.g
MAINNET_GENESIS_HASH=704105c268a01093f18e896767086efa68b8045e
TESTNET_GENESIS=testnet.g
TESTNET_GENESIS_HASH=ba37d578249da67cb5744069cc54f49a6938030d

deps:
	go get ./...

test:
	${TEST_SCRIPT}

build:
	docker build -t rosetta-fantom:latest https://github.com/Fantom-foundation/rosetta-fantom.git

build-local:
	docker build -t rosetta-fantom:latest .

build-release:
	# make sure to always set version with vX.X.X
	docker build -t rosetta-fantom:$(version) .;
	docker save rosetta-fantom:$(version) | gzip > rosetta-fantom-$(version).tar.gz;

update-tracer:
	curl https://raw.githubusercontent.com/ethereum/go-ethereum/master/eth/tracers/js/internal/tracers/call_tracer_js.js -o fantom/call_tracer.js

run-mainnet-online:
	docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -v "${PWD}/opera-data-mainnet:/data" -e "MODE=ONLINE" -e "NETWORK=MAINNET" -e "GENESIS=${MAINNET_GENESIS}" -e "GENESISHASH=${MAINNET_GENESIS_HASH}" -e "PORT=8080" -p 8080:8080 -p 5050:5050 rosetta-fantom:latest

run-mainnet-online-var:
	docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -v "/var/opera/mainnet:/data" -e "MODE=ONLINE" -e "NETWORK=MAINNET" -e "GENESIS=${MAINNET_GENESIS}" -e "GENESISHASH=${MAINNET_GENESIS_HASH}" -e "PORT=8080" -p 8080:8080 -p 5050:5050 rosetta-fantom:latest

run-mainnet-offline:
	docker run -d --rm -e "MODE=OFFLINE" -e "NETWORK=MAINNET" -e "PORT=8081" -p 8081:8081 rosetta-fantom:latest

run-testnet-online:
	docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -v "${PWD}/opera-data-testnet:/data" -e "MODE=ONLINE" -e "NETWORK=TESTNET" -e "GENESIS=${TESTNET_GENESIS}" -e "GENESISHASH=${TESTNET_GENESIS_HASH}" -e "PORT=8080" -p 8080:8080 -p 5050:5050 rosetta-fantom:latest

run-testnet-online-var:
	docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -v "/var/opera/testnet:/data" -e "MODE=ONLINE" -e "NETWORK=TESTNET" -e "GENESIS=${TESTNET_GENESIS}" -e "GENESISHASH=${TESTNET_GENESIS_HASH}" -e "PORT=8080" -p 8080:8080 -p 5050:5050 rosetta-fantom:latest

run-testnet-offline:
	docker run -d --rm -e "MODE=OFFLINE" -e "NETWORK=TESTNET" -e "PORT=8081" -p 8081:8081 rosetta-fantom:latest

run-mainnet-remote:
	docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -e "MODE=ONLINE" -e "NETWORK=MAINNET" -e "PORT=8080" -e "OPERA=$(opera)" -p 8080:8080 rosetta-fantom:latest

run-testnet-remote:
	docker run -d --rm --ulimit "nofile=${NOFILE}:${NOFILE}" -e "MODE=ONLINE" -e "NETWORK=TESTNET" -e "PORT=8080" -e "OPERA=$(opera)" -p 8080:8080 rosetta-fantom:latest

check-comments:
	${GOLINT_INSTALL}
	${GOLINT_CMD} -set_exit_status ${GO_FOLDERS} .
	go mod tidy

lint: | check-comments
	golangci-lint run --timeout 2m0s -v -E ${LINT_SETTINGS},gomnd

add-license:
	${ADDLICENSE_INSTALL}
	${ADDLICENCE_SCRIPT} .

check-license:
	${ADDLICENSE_INSTALL}
	${ADDLICENCE_SCRIPT} -check .

shorten-lines:
	${GOLINES_INSTALL}
	${GOLINES_CMD} -w --shorten-comments ${GO_FOLDERS} .

format:
	gofmt -s -w -l .
	${GOIMPORTS_CMD} -w .

check-format:
	! gofmt -s -l . | read
	! ${GOIMPORTS_CMD} -l . | read

salus:
	docker run --rm -t -v ${PWD}:/home/repo coinbase/salus

spellcheck:
	${SPELLCHECK_CMD} -error .

coverage:
	${GOVERALLS_INSTALL}
	if [ "${COVERALLS_TOKEN}" ]; then ${TEST_SCRIPT} -coverprofile=c.out -covermode=count; ${GOVERALLS_CMD} -coverprofile=c.out -repotoken ${COVERALLS_TOKEN}; fi

coverage-local:
	${TEST_SCRIPT} -cover

mocks:
	rm -rf mocks;
	mockery --dir services --all --case underscore --outpkg services --output mocks/services;
	mockery --dir opera --all --case underscore --outpkg ethereum --output mocks/opera;
	${ADDLICENSE_INSTALL}
	${ADDLICENCE_SCRIPT} .;
