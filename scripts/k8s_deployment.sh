#!/bin/sh

cd /opt/ansible-runner
cp -f configs/config.docker.yaml deplyment/base/config.yaml
cp -rf cert deplyment/base/

# Build the project
kubectl apply -k deployment/base

# cleanup
rm -rf deployment/base/config.yaml
rm -rf deployment/base/cert
