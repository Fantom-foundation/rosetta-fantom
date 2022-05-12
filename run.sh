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
  export OPERA_ARGS="$MAINNET_OPERA_ARGS"
  GENESIS="$MAINNET_GENESIS"
  GENESISHASH="$MAINNET_GENESIS_HASH"
  SNAPSHOT="$MAINNET_SNAPSHOT"
  SNAPSHOTMD5="$MAINNET_SNAPSHOT_MD5"
elif [ "$NETWORK" == "TESTNET" ]; then
  export OPERA_ARGS="$TESTNET_OPERA_ARGS"
  GENESIS="$TESTNET_GENESIS"
  GENESISHASH="$TESTNET_GENESIS_HASH"
else
  echo "Unrecognized NETWORK variable!"
  exit 53
fi

if [ "$MODE" == "ONLINE" ]; then

  # Download the genesis file if not exists
  echo "Downloading the genesis file $GENESIS if not exists"
  test -f "/data/$GENESIS" || axel -n 10 -o "/data/$GENESIS" "https://opera.fantom.network/$GENESIS" || ERRCODE=$?
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

  # Use snapshot if available and the database is not initialized yet
  if [ "$SNAPSHOT" != "" ] && [ ! -d "/data/chaindata" ]; then

    # Download the snapshot
    echo "Downloading the snapshot archive $SNAPSHOT if not exists"
    test -f "/data/$SNAPSHOT" || axel -n 10 -o "/data/$SNAPSHOT" "https://download.fantom.network/$SNAPSHOT" || ERRCODE=$?
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
    echo "Extracting the snapshot archive"
    tar --extract --file="/data/$SNAPSHOT" --strip-components=1 --directory="/data/" || ERRCODE=$?
    if [ $ERRCODE != 0 ]; then
      echo "Failed to extract the snapshot file /data/$SNAPSHOT ($ERRCODE)"
      exit 56
    fi

    echo "Extracted, removing the snapshot archive now"
    rm -f "/data/$SNAPSHOT"
  fi

else
  echo "Offline mode - skipping genesis/snapshot file"
fi

exec /app/rosetta-fantom run
