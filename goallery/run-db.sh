#!/bin/sh

docker-compose up -d

echo 
echo 

docker ps --format 'table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}' -f name=goallery_db

echo 
echo 

sleep 2

docker logs goallery_db

