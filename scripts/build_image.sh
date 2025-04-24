#!/bin/sh
cd /opt/ansible-runner

mkdir -p build_context/runner
cp -r /opt/ansible-sid-operator build_context/
cp -r /opt/ansible-runner/cert build_context/runner/
mkdir -p build_context/runner/configs
cp -r /opt/ansible-runner/configs/config.docker.yaml build_context/runner/configs/config.yaml

# Build the project
docker build -t ansible-runner:$1 -t ansible-runner:latest -f deployment/Dockerfile .

if [ "$2" = "push" ]; then
    # Push the image to artifactory
    docker login -u fsun4 artifactory.platformdxc-mg.com/vpcp4s-docker

    docker tag ansible-runner:$1 artifactory.platformdxc-mg.com/vpcp4s-docker/ansible-runner:$1

    docker push artifactory.platformdxc-mg.com/vpcp4s-docker/ansible-runner:$1

    docker push artifactory.platformdxc-mg.com/vpcp4s-docker/ansible-runner:latest

    docker rmi artifactory.platformdxc-mg.com/vpcp4s-docker/ansible-runner:$1

    docker rmi artifactory.platformdxc-mg.com/vpcp4s-docker/ansible-runner:latest

    docker logout artifactory.platformdxc-mg.com/vpcp4s-docker
fi

rm -rf build_context