#!/usr/bin/env bash

set -o errexit

BINARY_PATH=$1

if [[ $BINARY_PATH != *"windows_amd64"* ]]; then
  echo "=> $1 is not windows binary, skipping..."
  exit 0
else
  echo "=> signing $1..."
fi

# create cert files
CERT_FILE=cert.pem
if [ ! -f "$CERT_FILE" ]; then
  echo "$SIGNING_CERT" > "$CERT_FILE"
fi
KEY_FILE=key.pem
if [ ! -f "$KEY_FILE" ]; then
  echo "$SIGNING_KEY" > "$KEY_FILE"
fi

osslsigncode sign -h sha512 -certs cert.pem -key key.pem -n "GSP - Git Simple Packager" -i "https://websbygeorge.com" -t "http://timestamp.comodoca.com/authenticode" -in $BINARY_PATH -out signed.exe
if [[ $? -ne 0 ]]; then
  echo "Could not sign $BINARY_PATH"
  exit 1
fi

echo "Signing process was successful: ${BINARY_PATH}"

# replace unsigned binary with the signed one
mv signed.exe ${BINARY_PATH}
