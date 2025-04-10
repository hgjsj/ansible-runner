#!/bin/sh

cd /opt/ansible-runner
cp -f configs/config.docker.yaml deployment/base/config.yaml
cp -rf cert deployment/base/

# Build the project
kubectl apply -k deployment/base

# cleanup
rm -rf deployment/base/config.yaml
rm -rf deployment/base/cert
