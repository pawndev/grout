#!/bin/sh
CUR_DIR="$(dirname "$0")"
cd "$CUR_DIR"/grout || exit 1

export CFW=CROSSMIX
export LD_LIBRARY_PATH=lib:"$LD_LIBRARY_PATH"

./grout
