#!/usr/bin/env bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

mkdir "$DIR"/../configs/certs -p
docker run -v "$DIR"/../configs/certs:/export -v "$DIR":/scripts alpine /scripts/certs.sh
