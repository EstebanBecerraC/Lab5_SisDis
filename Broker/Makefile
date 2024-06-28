docker: imagen
	sudo docker run --name broker -p 50051:50051 broker:latest

imagen: dependencias
	sudo docker build . -t broker:latest

dependencias:
	go mod init go-container
	protoc --go_out=./proto --go-grpc_out=./proto registros.proto
	go mod tidy
	
clean:
	rm -f go.mod
	rm -f go.sum
	rm -f ./proto/registros_grpc.pb.go
	rm -f ./proto/registros.pb.go
	sudo docker stop broker
	sudo docker rm broker
	sudo docker rmi broker:latest
	