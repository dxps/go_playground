#!/bin/sh

echo "GET http://localhost:8000" | vegeta attack -duration=4s -rate=4 | vegeta report
