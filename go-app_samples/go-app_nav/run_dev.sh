#!/bin/sh

find . -name "*.go" | entr -r make dev

