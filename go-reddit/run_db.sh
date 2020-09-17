#!/bin/sh

CONTAINER_NAME=go-reddit-db

docker-compose up -d

echo
echo

docker ps --format 'table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' -f name=${CONTAINER_NAME}

echo
echo

sleep 5

docker logs ${CONTAINER_NAME}
