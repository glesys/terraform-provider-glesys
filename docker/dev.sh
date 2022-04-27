#!/bin/sh

echo ""
read -p "Username:" username

echo ""
read -s -p "Secret: " password

export GLESYS_USERID="${username}"
export GLESYS_TOKEN="${password}"

docker run --rm -it --entrypoint /bin/sh \
  -e GLESYS_USERID -e GLESYS_TOKEN \
  -v "$PWD":/app \
  -w /app \
  cgr.dev/chainguard/go:latest-dev
