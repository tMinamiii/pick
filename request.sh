#!/bin/bash
BODY='{\"consumer_key\": \"$1\", \"redirect_uri\":\"pocketapp1234:authorizationFinished\"}'
BODY=$(eval echo $BODY)
CODE=$(curl -s -H "Content-Type: application/json; charset=UTF-8" -d "${BODY}" -X POST https://getpocket.com/v3/oauth/request)
VAR=($(echo $CODE | tr -s '=' ' '))
URL="https://getpocket.com/auth/authorize?request_token=${VAR[1]}&redirect_uri=pocketapp1234:authorizationFinished"
echo $URL
