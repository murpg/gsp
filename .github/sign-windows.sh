#!/usr/bin/env bash

set -o errexit

BINARY_PATH=$1

if [[ $BINARY_PATH != *"gsp_windows_amd64"* ]]; then
  echo "=> $1 is not windows binary, skipping..."
  exit 0
fi

(osslsigncode sign -h sha512 \
  -certs "$SIGNING_CERT"
  -key "$SIGNING_KEY"
  -n "GSP - Git Simple Packager"
  -i "https://websbygeorge.com"
  -t "http://timestamp.comodoca.com/authenticode"
  -in "$BINARY_PATH"
  -out "$BINARY_PATH"_signed)
if [[ $? -ne 0 ]]; then
  echo "Could not sign $BINARY_PATH"
  exit 1
fi

# replace unsigned binary with the signed one
mv ${BINARY_PATH}_signed ${BINARY_PATH}
