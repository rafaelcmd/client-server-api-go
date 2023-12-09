# Dockerized Client-Server Application

## Building Container Images

To build the container images, use the following commands:

```bash
docker build -t server:latest -f Dockerfile.Server .
docker build -t client:latest -f Dockerfile.Client .
```

## Creating a Docker Network

To create a Docker network for communication between server and client containers:

```bash
docker network create client-server-network
```

## Running the Containers

To run the server container:

```bash
docker run --name server -p 8080:8080 --network client-server-network server
```

To run the client container:

```bash
docker run --name client --network client-server-network --volume ./file:/app/file client
```

## Running the Client Multiple Times

If you need to run the client container more than once:

```bash
docker start client
```