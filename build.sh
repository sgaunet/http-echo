#!/usr/bin/env bash

docker build . -t sgaunet/http-echo:latest
rc=$?

if [ "$rc" != "0" ]
then
    echo "Build Failed"
    exit 1
fi

docker push sgaunet/http-echo:latest
