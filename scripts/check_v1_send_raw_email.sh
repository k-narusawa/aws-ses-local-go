#!/bin/bash

AWS_ENDPOINT_URL="http://localhost:8080"

aws ses send-raw-email \
  --raw-message file://./scripts/message.json \
  --endpoint-url "http://localhost:8080"

if [ $? -eq 0 ]; then
  echo "Test succeeded: SES mock responded correctly."
else
  echo "Test failed: SES mock did not respond as expected."
  exit 1
fi
