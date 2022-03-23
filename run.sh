#!/bin/bash
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
