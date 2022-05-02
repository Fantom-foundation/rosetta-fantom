#!/bin/bash
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

# Script to be started in the rosetta-fantom docker container
# Downloads and verifies the genesis file and starts the rosetta-fantom (which start opera if needed)
ERRCODE=0

echo "Running with network $NETWORK in $MODE mode"

if [ "$NETWORK" == "MAINNET" ]; then
  GENESIS=mainnet.g
  GENESISHASH=704105c268a01093f18e896767086efa68b8045e

  # pruned snapshot for faster testing
  SNAPSHOT=opera_pruned_20apr22.tgz
  SNAPSHOTMD5=6b142110281f31c831c3182070687db2
elif [ "$NETWORK" == "TESTNET" ]; then
  GENESIS=testnet.g
  GENESISHASH=ba37d578249da67cb5744069cc54f49a6938030d
else
  echo "Unrecognized NETWORK variable!"
  exit 53
fi

if [ "$MODE" == "ONLINE" ]; then

  # Download the genesis file if not exists
  echo "Downloading the genesis file $GENESIS if not exists"
  test -f "/data/$GENESIS" || wget -O "/data/$GENESIS" "https://opera.fantom.network/$GENESIS" || ERRCODE=$?
  if [ $ERRCODE != 0 ]; then
    echo "Failed to download the genesis file ($ERRCODE)"
    exit 51
  fi

  # Check the genesis file checksum
  echo "Checking the genesis file checksum"
  echo "$GENESISHASH  /data/$GENESIS" | sha1sum -c - || ERRCODE=$?
  if [ $ERRCODE != 0 ]; then
    echo "Invalid checksum of the genesis file /data/$GENESIS (not equal to $GENESISHASH)"
    exit 52
  fi

  # Use snapshot if available and not initialized yet
  if [ "$SNAPSHOT" != "" ] && [ ! -d "/data/chaindata" ]; then

    # Download the snapshot
    echo "Downloading the snapshot archive $SNAPSHOT"
    test -f "/data/$SNAPSHOT" || wget -O "/data/$SNAPSHOT" "https://download.fantom.network/$SNAPSHOT" || ERRCODE=$?
    if [ $ERRCODE != 0 ]; then
      echo "Failed to download the snapshot file $SNAPSHOT ($ERRCODE)"
      exit 54
    fi

    # Check the snapshot archive checksum
    echo "Checking the snapshot archive checksum"
    echo "$SNAPSHOTMD5  /data/$SNAPSHOT" | md5sum -c - || ERRCODE=$?
    if [ $ERRCODE != 0 ]; then
      echo "Invalid checksum of the snapshot file /data/$SNAPSHOT (not equal to $SNAPSHOTMD5)"
      exit 55
    fi

    # Extract the .opera/chaindata from the archive into /data/chaindata
    tar --extract --file="/data/$SNAPSHOT" --strip-components=1 --directory="/data/" || ERRCODE=$?
    if [ $ERRCODE != 0 ]; then
          echo "Failed to extract the snapshot file /data/$SNAPSHOT"
          exit 56
        fi
  fi

else
  echo "Not online mode - skipping genesis file"
fi

/app/rosetta-fantom run
