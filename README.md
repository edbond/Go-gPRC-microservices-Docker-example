# Ports

The project contains an example of microservices written in Go.
The purpose of the system is to store, update and return marine port objects.
There are 2 microservices (Client and Ports) that communicate using gRPC.
Client service exposes the HTTP interface to add and list ports.

## Ports service

Ports service is used to store ports.
Port can be inserted or updated in DB.
Ports service is a gRPC server.
The repository is an interface that includes several functions to upsert (insert or update),
list all ports.
There is an in-memory implementation of the repository that stores all data in a map.

## Client service

Client service on startup reads ports from JSON file and calls upsert on ports service to
store ports.
Client service provides HTTP API to get a list of all ports.

## Folders structure

There are 3 folders:

- `ports` contains definitions of a Port structure, ports JSON file.
- `ports_service` contains the implementation of Ports Service, repository and gRPC server.
- `client_service` contains the implementation of Client Service.

## Generate structures and gRPC client and server from protobuf 

Script `gen.sh` generates protobuf structures, functions, and gRPC client and server code.
`ports.proto` file contains a description of the Ports service.

## Tests

To run tests run the script `test.sh`. The script will run tests in all 3 folders.

## Running

Client and Ports services contain dockerfiles to build and run them in containers.
`docker-compose.yml` file contains instructions on how to launch services in docker.

Run `docker-compose up --build` to start both services.

```
Starting ports_ports_service_1 ... done
Starting ports_client_service_1 ... done
Attaching to ports_ports_service_1, ports_client_service_1
ports_service_1   | time="2021-02-04T20:17:35Z" level=info msg="Starting Ports GRPC Server on port 4040" Server="Ports GRPC Server"
client_service_1  | time="2021-02-04T20:17:37Z" level=info msg="Total ports loaded: 1632, failed: 0" Server=Client
client_service_1  | time="2021-02-04T20:17:37Z" level=info msg="🆙 Starting server at port 8080" Server=Client
```

Press `Ctrl+C` to stop services gracefully.

