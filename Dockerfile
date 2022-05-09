# Copyright 2020 Coinbase, Inc.
# Copyright 2022 Fantom Foundation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Compile Opera
FROM golang:1.18 as opera-builder

# VERSION: go-opera release/txtracing/1.1.0-rc.5
RUN git clone https://github.com/Fantom-foundation/go-opera \
  && cd go-opera \
  && git -c advice.detachedHead=false checkout 004e7a6e206210b440f7a510845ca6df106fedeb

RUN cd go-opera \
  && make

RUN mkdir -p /app \
  && mv go-opera/build/opera /app/opera \
  && rm -rf go-opera

# Compile rosetta-fantom
FROM golang:1.18 as rosetta-builder

RUN git clone https://github.com/Fantom-foundation/rosetta-fantom src \
  && cd src \
  && git -c advice.detachedHead=false checkout 3382336e6b35cc673440dd29a3baf7c0912dcc8d

RUN cd src \
  && go build

RUN mkdir -p /app \
  && mv src/rosetta-fantom /app/rosetta-fantom \
  && mkdir /app/fantom \
  && mv src/fantom/call_tracer.js /app/fantom/call_tracer.js \
  && mv src/fantom/opera.toml /app/fantom/opera.toml \
  && mv src/run.sh /app/run.sh \
  && rm -rf src

## Build Final Image
FROM ubuntu:20.04

RUN apt-get update && apt-get install -y ca-certificates wget axel && update-ca-certificates

RUN mkdir -p /app \
  && chown -R nobody:nogroup /app \
  && mkdir -p /data \
  && chown -R nobody:nogroup /data

WORKDIR /app

# Copy binary from opera-builder
COPY --from=opera-builder /app/opera /app/opera

# Copy binary and assets from rosetta-builder
COPY --from=rosetta-builder /app/fantom /app/fantom
COPY --from=rosetta-builder /app/rosetta-fantom /app/rosetta-fantom
COPY --from=rosetta-builder /app/run.sh /app/run.sh

# Set permissions for everything added to /app
RUN chmod -R 755 /app/*

CMD ["/bin/bash", "/app/run.sh"]
