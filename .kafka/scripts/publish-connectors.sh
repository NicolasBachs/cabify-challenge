#!/bin/bash
for file in /connectors/*.json; do
  echo "Loading connector $file"
  echo "$(cat $file) "
  response=$(curl -s -X POST -H 'Content-Type: application/json' --data "$(cat $file)" http://kafka-connect:8083/connectors)
  echo "Response: $response"
done