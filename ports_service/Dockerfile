FROM golang:1.20

# Port to run GRPC server
ENV PORTS_GRPC_PORT=4040
EXPOSE ${PORTS_GRPC_PORT}

WORKDIR /go/ports

COPY ./ports_service/ ./ports_service

WORKDIR /go/ports/ports_service

RUN go build -o ports_service ./cmd/ports/main.go

ENTRYPOINT ./ports_service