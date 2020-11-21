#!/bin/bash

go run server.go &

while ! nc -z localhost 9999; do
  sleep 0.1
done

curl http://localhost:9999/profiles/admin
curl http://localhost:9999/profiles/foo
pkill -TERM -P $!
