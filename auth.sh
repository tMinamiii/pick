#!/bin/bash
PAYLOAD='{\"consumer_key\": \"$1\", \"code\":\"$2\"}'
PAYLOAD=$(eval echo $PAYLOAD)
echo $PAYLOAD
ACCESS_TOKEN=$(curl -s -H "Content-Type: application/json; charset=UTF-8" -d "${PAYLOAD}" -X POST https://getpocket.com/v3/oauth/authorize)
echo $ACCESS_TOKEN
