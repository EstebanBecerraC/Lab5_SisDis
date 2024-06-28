docker: imagen
	docker run --name broker -p 50051:50051 broker:latest

imagen: dependencias
	docker build . -t broker:latest

dependencias:
	go mod init go-container
	protoc --go_out=./proto --go-grpc_out=./proto registros.proto
	go mod tidy
	
clean:
	rm -f go.mod
	rm -f go.sum
	rm -f ./proto/registros_grpc.pb.go
	rm -f ./proto/registros.pb.go
	docker stop broker
	docker rm broker
	docker rmi broker:latest
	