#!/bin/sh

export https_proxy=http://138.35.86.205:8080
export http_proxy=http://138.35.86.205:8080
export no_proxy=127.0.0.1,localhost,.env01.mcloud.entsvcs.net

cd /opt/ansible-runner
# pull latest code the project
git pull origin main 

# Build the project
docker build -t ansible-runner:$1 -f deploy/Dockerfile .