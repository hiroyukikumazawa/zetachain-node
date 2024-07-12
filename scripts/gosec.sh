#!/usr/bin/env bash

docker run -it --rm -w /node -v "$(pwd):/node" ghcr.io/zeta-chain/gosec:2.21.0-zeta -exclude-generated -exclude-dir testutil ./...