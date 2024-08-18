#!/bin/bash

EMAIL_SOURCE="test@example.com"
EMAIL_DESTINATION="recipient@example.com"
EMAIL_SUBJECT="Test Email"
EMAIL_BODY="This is a test email."

AWS_ENDPOINT_URL="http://localhost:8080"

aws sesv2 send-email \
--from-email-address "${EMAIL_SOURCE}" \
--destination "ToAddresses=${EMAIL_DESTINATION}" \
--content "Simple={Subject={Data=${EMAIL_SUBJECT},Charset=utf-8},Body={Text={Data=${EMAIL_BODY},Charset=utf-8}}}" \
--endpoint-url "${AWS_ENDPOINT_URL}" > /dev/null

if [ $? -eq 0 ]; then
  echo "Test succeeded: SES mock responded correctly."
else
  echo "Test failed: SES mock did not respond as expected."
  exit 1
fi
