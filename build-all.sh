#!/usr/bin/env bash
set -e
export GIT_COMMIT=$(git rev-parse --short HEAD)

if [[ -z "$NDK_ROOT" ]]; then
   echo "NDK_ROOT environmental variable must be set to the path of your Android NDK."
   echo "GO cannot build binaries to run on Android without NDK build tools."
   exit -1
fi

make all

echo "Build Finished"
