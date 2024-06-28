FROM golang:latest

WORKDIR /ClientesComandante

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY proto/registros_grpc.pb.go ./proto/
COPY proto/registros.pb.go ./proto/
COPY kais.go .

RUN go build -o cc kais.go

ENTRYPOINT ["/ClientesComandante/cc"]
