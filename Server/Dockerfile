FROM golang:latest

WORKDIR /ServidorFulcrum

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY proto/registros_grpc.pb.go ./proto/
COPY proto/registros.pb.go ./proto/
COPY fulcrum.go .

RUN go build -o sf fulcrum.go

EXPOSE 50052

ENTRYPOINT ["/ServidorFulcrum/sf"]