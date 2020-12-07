#!/bin/sh

##
## Passing "-routes" flag triggers the generation of router documentation.
##

echo "Generating the routes documentation to 'routes.md' file ..." 

go run main.go -routes 1> routes.md

