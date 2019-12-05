#!/bin/bash

# go clean -modcache
go mod download && \
go test -v ./... && \
docker-compose up -d && \
echo "* Wait sample data to be generated..." && \
while [ "$(docker-compose ps | grep datagen | grep Up | wc -l | awk '{print $1}')" == "1" ]; do
	sleep 1
done

echo "* Run Go ksql client..."
go run main.go

docker-compose down


