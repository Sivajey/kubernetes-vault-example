#!/bin/bash

CLIENT_IMAGE="${CLIENT_IMAGE:-time-client}"
SERVER_IMAGE="${SERVER_IMAGE:-time-server}"
VERSION="${VERSION:-1.0.0}"

docker build -t "${CLIENT_IMAGE}:${VERSION}-01" 01-initial-example/client
docker build -t "${CLIENT_IMAGE}:${VERSION}-02" 02-using-secrets-from-vault/client/
docker build -t "${CLIENT_IMAGE}:${VERSION}-03" 03-using-certificates-from-vault/client/
docker build -t "${SERVER_IMAGE}:${VERSION}-01" 01-initial-example/server
docker build -t "${SERVER_IMAGE}:${VERSION}-02" 02-using-secrets-from-vault/server
docker build -t "${SERVER_IMAGE}:${VERSION}-03" 03-using-certificates-from-vault/server/

docker push "${CLIENT_IMAGE}:${VERSION}-01"
docker push "${CLIENT_IMAGE}:${VERSION}-02"
docker push "${CLIENT_IMAGE}:${VERSION}-03"
docker push "${SERVER_IMAGE}:${VERSION}-01"
docker push "${SERVER_IMAGE}:${VERSION}-02"
docker push "${SERVER_IMAGE}:${VERSION}-03"
