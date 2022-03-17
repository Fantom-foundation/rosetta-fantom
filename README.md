<h3 align="center">
   Rosetta Fantom
</h3>
<p align="center">
  <a href="https://github.com/Fantom-foundation/rosetta-fantom/blob/master/LICENSE.txt"><img src="https://img.shields.io/github/license/Fantom-foundation/rosetta-fantom.svg" /></a>
</p>

<p align="center"><b>
ROSETTA-FANTOM IS CONSIDERED <a href="https://en.wikipedia.org/wiki/Software_release_life_cycle#Alpha">ALPHA SOFTWARE</a>.
USE AT YOUR OWN RISK! COINBASE ASSUMES NO RESPONSIBILITY OR LIABILITY IF THERE IS A BUG IN THIS IMPLEMENTATION.
</b></p>

## Overview
`rosetta-fantom` provides an implementation of the Rosetta API for Fantom Opera in Golang.
If you haven't heard of the Rosetta API, you can find more information [here](https://rosetta-api.org).

## Features
* Comprehensive tracking of all FTM balance changes
* Stateless, offline, curve-based transaction construction (with address checksum validation)
* Idempotent access to all transaction traces and receipts

### Recommended OS Settings
To increase the load `rosetta-fantom` can handle, it is recommended to tune your OS
settings to allow for more connections. On a linux-based OS, you can run the following
commands ([source](http://www.tweaked.io/guide/kernel)):
```text
sysctl -w net.ipv4.tcp_tw_reuse=1
sysctl -w net.core.rmem_max=16777216
sysctl -w net.core.wmem_max=16777216
sysctl -w net.ipv4.tcp_max_syn_backlog=10000
sysctl -w net.core.somaxconn=10000
sysctl -p (when done)
```

You should also modify your open file settings to `100000`. This can be done on a linux-based OS
with the command: `ulimit -n 100000`.

## Usage
As specified in the [Rosetta API Principles](https://www.rosetta-api.org/docs/automated_deployment.html),
all Rosetta implementations must be deployable via Docker and support running via either an
[`online` or `offline` mode](https://www.rosetta-api.org/docs/node_deployment.html#multiple-modes).

**YOU MUST INSTALL DOCKER FOR THE FOLLOWING INSTRUCTIONS TO WORK. YOU CAN DOWNLOAD
DOCKER [HERE](https://www.docker.com/get-started).**

### Install
Running the following commands will create a Docker image called `rosetta-fantom:latest`.

#### From GitHub
To download the pre-built Docker image from the latest release, run:
```text
curl -sSfL https://raw.githubusercontent.com/coinbase/rosetta-fantom/master/install.sh | sh -s
```

_Do not try to install rosetta-fantom using GitHub Packages!_


#### From Source
After cloning this repository, run:
```text
make build-local
```

### Run
Running the following commands will start a Docker container in
[detached mode](https://docs.docker.com/engine/reference/run/#detached--d) with
a data directory at `<working directory>/opera-data` with a genesis file and the Rosetta API accessible
at port `8080`.

#### Configuration Environment Variables
* `MODE` (required) - Determines if Rosetta can make outbound connections. Options: `ONLINE` or `OFFLINE`.
* `NETWORK` (required) - Network to launch and/or communicate with. Options: `MAINNET` or `TESTNET`.
* `PORT`(required) - Which port to use for Rosetta.
* `OPERA` (optional) - Point to a remote `opera` node instead of initializing one
* `SKIP_ADMIN` (optional, default: `FALSE`) - Instruct Rosetta to not use the `opera` `admin` RPC calls. This is typically disabled by hosted blockchain node services.

#### Mainnet:Online
```text
docker run -d --rm --ulimit "nofile=100000:100000" -v "$(pwd)/opera-data:/data" -e "MODE=ONLINE" -e "NETWORK=MAINNET" -e "PORT=8080" -p 8080:8080 -p 30303:30303 rosetta-fantom:latest
```
_If you cloned the repository, you can run `make run-mainnet-online`._

#### Mainnet:Online (Remote)
```text
docker run -d --rm --ulimit "nofile=100000:100000" -e "MODE=ONLINE" -e "NETWORK=MAINNET" -e "PORT=8080" -e "OPERA=<NODE URL>" -p 8080:8080 -p 30303:30303 rosetta-fantom:latest
```
_If you cloned the repository, you can run `make run-mainnet-remote opera=<NODE URL>`._

#### Mainnet:Offline
```text
docker run -d --rm -e "MODE=OFFLINE" -e "NETWORK=MAINNET" -e "PORT=8081" -p 8081:8081 rosetta-fantom:latest
```
_If you cloned the repository, you can run `make run-mainnet-offline`._

#### Testnet:Online
```text
docker run -d --rm --ulimit "nofile=100000:100000" -v "$(pwd)/opera-data:/data" -e "MODE=ONLINE" -e "NETWORK=TESTNET" -e "PORT=8080" -p 8080:8080 -p 30303:30303 rosetta-fantom:latest
```
_If you cloned the repository, you can run `make run-testnet-online`._

#### Testnet:Online (Remote)
```text
docker run -d --rm --ulimit "nofile=100000:100000" -e "MODE=ONLINE" -e "NETWORK=TESTNET" -e "PORT=8080" -e "OPERA=<NODE URL>" -p 8080:8080 -p 30303:30303 rosetta-fantom:latest
```
_If you cloned the repository, you can run `make run-testnet-remote opera=<NODE URL>`._

#### Testnet:Offline
```text
docker run -d --rm -e "MODE=OFFLINE" -e "NETWORK=TESTNET" -e "PORT=8081" -p 8081:8081 rosetta-fantom:latest
```
_If you cloned the repository, you can run `make run-testnet-offline`._

## Testing with rosetta-cli
To validate `rosetta-fantom`, [install `rosetta-cli`](https://github.com/coinbase/rosetta-cli#install)
and run one of the following commands:
* `rosetta-cli check:data --configuration-file rosetta-cli-conf/testnet/config.json` - This command validates that the Data API implementation is correct using the Opera `testnet` node. It also ensures that the implementation does not miss any balance-changing operations.
* `rosetta-cli check:construction --configuration-file rosetta-cli-conf/testnet/config.json` - This command validates the Construction API implementation. It also verifies transaction construction, signing, and submissions to the `testnet` network.
* `rosetta-cli check:data --configuration-file rosetta-cli-conf/mainnet/config.json` - This command validates that the Data API implementation is correct using the Opera `mainnet` node. It also ensures that the implementation does not miss any balance-changing operations.

## Issues
Interested in helping fix issues in this repository? You can find to-dos in the [Issues](https://github.com/Fantom-foundation/rosetta-fantom/issues) section.

## Development
* `make deps` to install dependencies
* `make test` to run tests
* `make lint` to lint the source code
* `make salus` to check for security concerns
* `make build-local` to build a Docker image from the local context
* `make coverage-local` to generate a coverage report

## License
This project is available open source under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0).

This project is based on [Rosetta Ethereum](https://github.com/coinbase/rosetta-ethereum), a reference implementation of the Rosetta API by Coinbase.
