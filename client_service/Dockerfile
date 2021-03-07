FROM golang:1.16

# Port to run GRPC server
ENV HTTP_PORT=8080
EXPOSE ${HTTP_PORT}

ENV PORTS_ADDRESS=ports_service:4040
ENV PORTS_JSON=/go/client/ports.json

#RUN go env -w GO111MODULE=auto
#ENV GO111MODULE=auto

COPY ./client_service/ ./client

WORKDIR /go/client

RUN go build -o client_service ./cmd/client/main.go

ENTRYPOINT ./client_service