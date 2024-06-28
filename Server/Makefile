docker: imagen
	sudo docker run --name servidorfulcrum -p 50052:50051 servidor-fulcrum:latest

imagen: dependencias
	sudo docker build . -t servidor-fulcrum:latest

dependencias:
	go mod init go-container
	protoc --go_out=./proto --go-grpc_out=./proto registros.proto
	go mod tidy
	
clean:
	rm -r Sectores
	mkdir Sectores
	rm -f go.mod
	rm -f go.mod
	rm -f go.sum
	rm -f ./proto/registros_grpc.pb.go
	rm -f ./proto/registros.pb.go
	sudo docker stop servidorfulcrum
	sudo docker rm servidorfulcrum
	sudo docker rmi servidor-fulcrum:latest
	