#!/bin/sh

export https_proxy=http://138.35.86.205:8080
export http_proxy=http://138.35.86.205:8080
export no_proxy=127.0.0.1,localhost,.env01.mcloud.entsvcs.net

cd /opt/ansible-runner

mkdir build_context/runner
cp -r /opt/ansible-sid-operator build_context/
cp -r /opt/ansible-runner/cert build_context/runner/
cp -r /opt/ansible-runner/configs build_context/runner/

# Build the project
docker build -t ansible-runner:$1 -f deployment/Dockerfile .

rm -rf build_context