FROM golang:latest

WORKDIR /Broker

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY proto/registros_grpc.pb.go ./proto/
COPY proto/registros.pb.go ./proto/
COPY broker.go .

RUN go build -o bk broker.go

EXPOSE 50051

ENTRYPOINT ["/Broker/bk"]