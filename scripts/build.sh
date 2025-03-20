#!/bin/sh

cd /opt/ansible-runner

# pull latest code the project
echo Pull latest code from github
git pull origin main 

# Build the project
echo Building the project
go build -o ansible-runner

# Run the API server in background
echo Running the API server
./ansible-runner server server -p 8081 &

# Run the worker node in background
echo Running the worker node
./ansible-runner worker &