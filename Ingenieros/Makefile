docker: imagen
	sudo docker run -it --name clienteingenieros cliente-ingenieros:latest

imagen: dependencias
	sudo docker build . -t cliente-ingenieros:latest

dependencias:
	go mod init go-container
	protoc --go_out=./proto --go-grpc_out=./proto registros.proto
	go mod tidy
	
clean:
	rm -f go.mod
	rm -f go.sum
	rm -f ./proto/registros_grpc.pb.go
	rm -f ./proto/registros.pb.go
	sudo docker stop clienteingenieros
	sudo docker rm clienteingenieros
	sudo docker rmi cliente-ingenieros:latest
	