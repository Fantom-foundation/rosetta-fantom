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

if [ "$GENESIS" != "" ]; then

  # Download the genesis file
  echo "Downloading the genesis file"
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

else
  echo "Genesis file skipped"
fi

/app/rosetta-fantom run
