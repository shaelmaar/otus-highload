#!/bin/bash
set -e

echo "" > /dev/tcp/localhost/5000 || exit 1
echo "" > /dev/tcp/localhost/5001 || exit 1
