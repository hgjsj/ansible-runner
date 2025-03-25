#!/bin/sh

export https_proxy=http://138.35.86.205:8080
export http_proxy=http://138.35.86.205:8080
export no_proxy=127.0.0.1,localhost,.env01.mcloud.entsvcs.net

cd /opt/ansible-runner

mkdir -p build_context/runner
cp -r /opt/ansible-sid-operator build_context/
cp -r /opt/ansible-runner/cert build_context/runner/
cp -r /opt/ansible-runner/configs build_context/runner/

# Build the project
docker build -t ansible-runner:$1 -f deployment/Dockerfile .

if $2="push"; then
    # Push the image to artifactory
    docker login -u fsun4 artifactory.platformdxc-mg.com/vpcp4s-docker

    docker tag ansible-runner:$1 artifactory.platformdxc-mg.com/vpcp4s-docker/ansible-runner:$1

    docker push artifactory.platformdxc-mg.com/vpcp4s-docker/ansible-runner:$1

    docker logout artifactory.platformdxc-mg.com/vpcp4s-docker
fi

rm -rf build_context