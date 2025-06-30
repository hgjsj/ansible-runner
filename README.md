# Ansible Runner

Ansible Runner is a backend service built with Go to manage and execute Ansible playbooks via a RESTful API. It is designed to be run in a containerized environment and can be easily deployed to Kubernetes.

The system uses a worker-based architecture with Redis and Asynq for asynchronous job processing and MongoDB for persistent storage of job status and results.

## Features

-   **RESTful API**: Provides endpoints to manage playbooks, inventories, and execution jobs.
-   **Asynchronous Job Execution**: Long-running Ansible tasks are handled by background workers, preventing API timeouts.
-   **Persistent Storage**: Job definitions, status, and output are stored in MongoDB.
-   **Containerized**: Comes with a `Dockerfile` for easy container image building.
-   **Kubernetes Ready**: Includes deployment manifests for Kubernetes, managed with Kustomize.
-   **Flexible Configuration**: Application behavior can be configured via YAML files for different environments.

## Architecture Overview

The application consists of two main components:

1.  **Server**: An API server that exposes endpoints for users to submit and query jobs. It acts as a producer, creating tasks and pushing them to a Redis queue.
2.  **Worker**: A background process that consumes tasks from the Redis queue. It is responsible for executing the Ansible playbooks and updating the job status in MongoDB.

This architecture decouples the API from the execution logic, ensuring the system is scalable and resilient.

```
+---------+       +------------------+       +-------+
|         |------>|   Redis (Asynq)  |<------|       |
|  Server |       |   (Task Queue)   |       | Worker|
|  (API)  |       +------------------+       |       |
|         |                                  |       |
|         |<-------------------------------->|       |
+---------+       +------------------+       +-------+
     ^            |     MongoDB      |            ^
     |            |  (Data Store)    |            |
     +------------|                  |------------+
                  +------------------+
```

## Getting Started

### Prerequisites

-   Go (version 1.18+)
-   Docker
-   Kubernetes (e.g., Minikube, Kind)
-   `kubectl`

### Configuration

1.  Copy the sample configuration file:
    ```sh
    cp configs/config.yaml configs/config.local.yaml
    ```
2.  Edit `configs/config.local.yaml` to match your local environment settings (e.g., database connection strings).

### Build

You can build the application binaries using the provided script:

```sh
./scripts/build.sh
```

This will create `server` and `worker` executables in the `bin/` directory.

### Running Locally

1.  **Start Server**:
    ```sh
    ./bin/server
    ```
2.  **Start Worker**:
    ```sh
    ./bin/worker
    ```

## API Endpoints

The following are some of the core API endpoints (details may vary):

-   `POST /api/v1/playbooks`: Upload a new Ansible playbook.
-   `POST /api/v1/inventories`: Create or update an inventory.
-   `POST /api/v1/jobs`: Create a new job to run a playbook against an inventory.
-   `GET /api/v1/jobs/{job_id}`: Get the status and output of a specific job.

## Deployment

This project is designed to be deployed on Kubernetes.

### Build the Docker Image

Use the build script to create a Docker image for the application:

```sh
# Make sure to be logged into your container registry
./scripts/build_image.sh <your-registry>/ansible-runner:latest
```

### Deploy to Kubernetes

The `deployment/` directory contains all the necessary Kubernetes manifests, managed by Kustomize.

You can deploy the application using the deployment script:

```sh
# This script will apply the kustomization for the 'dev' overlay
./scripts/k8s_deployment.sh
```

This will deploy the `server`, `worker`, and all required dependencies like `MongoDB` and `Redis` to your cluster.

## Project Structure

```
.
├── configs/         # Application configuration files
├── deployment/      # Docker & Kubernetes deployment files
├── internal/        # Core application code
│   ├── api/         # HTTP API handlers and routing
│   ├── cmd/         # Application entrypoints (server, worker)
│   ├── db/          # Database interaction logic
│   ├── models/      # Data models and structures
│   ├── services/    # Business logic
│   ├── tasks/       # Asynchronous task definitions
│   └── utils/       # Utility functions
├── scripts/         # Build and deployment scripts
├── go.mod           # Go module definition
└── main.go          # Main application entrypoint
```
