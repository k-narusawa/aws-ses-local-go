#!/bin/bash

EMAIL_SOURCE="test@example.com"
EMAIL_DESTINATION="recipient@example.com"
EMAIL_SUBJECT="Test Email"
EMAIL_BODY="This is a test email."

AWS_ENDPOINT_URL="http://localhost:8080"

aws ses send-email \
  --endpoint-url $AWS_ENDPOINT_URL \
  --to $EMAIL_DESTINATION \
  --from $EMAIL_SOURCE \
  --subject "$EMAIL_SUBJECT" \
  --text "$EMAIL_BODY" > /dev/null

if [ $? -eq 0 ]; then
  echo "Test succeeded: SES mock responded correctly."
else
  echo "Test failed: SES mock did not respond as expected."
  exit 1
fi
