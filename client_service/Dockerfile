FROM golang

# Port to run GRPC server
ENV HTTP_PORT=8080
EXPOSE ${HTTP_PORT}

ENV PORTS_ADDRESS=ports_service:4040
ENV PORTS_JSON=/go/ports/ports/ports.json

WORKDIR /go/ports

COPY ./client_service/ ./client_service
COPY ./ports/ ./ports

WORKDIR /go/ports/client_service

RUN go build -o client_service ./cmd/client/main.go

ENTRYPOINT ./client_service