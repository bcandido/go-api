#!/usr/bin/env bash

echo -e "==== building images\n"

echo -e "==== build go-api"
docker build -t go-api/go-api .

echo -e "\n==== build postgres"
docker build -t go-api/postgres deploy/docker/postgres/

echo -e "\n==== images built"