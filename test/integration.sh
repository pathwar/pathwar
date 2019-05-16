#!/bin/bash -xe
# depends on: httpie, jq

TARGET="${TARGET:-http://localhost:8000}"

http --check-status "${TARGET}"/api/ping
token=$(http --check-status "${TARGET}"/api/authenticate username=integration | jq -r .token)
test "$(http --check-status "${TARGET}"/api/user-session | jq -r .message)" = "no token provided"
test "$(http --check-status "${TARGET}"/api/user-session "Authorization:$token" | jq -r .username)" = "integration"
