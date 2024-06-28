FROM golang:latest

WORKDIR /ClientesIngenieros

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY proto/registros_grpc.pb.go ./proto/
COPY proto/registros.pb.go ./proto/
COPY ingenieros.go .

RUN go build -o ci ingenieros.go

ENTRYPOINT ["/ClientesIngenieros/ci"]