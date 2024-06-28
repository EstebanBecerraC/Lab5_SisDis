docker: imagen
	docker run -it --name clienteingenieros cliente-ingenieros:latest

imagen: dependencias
	docker build . -t cliente-ingenieros:latest

dependencias:
	go mod init go-container
	protoc --go_out=./proto --go-grpc_out=./proto registros.proto
	go mod tidy
	
clean:
	rm -f go.mod
	rm -f go.sum
	rm -f ./proto/registros_grpc.pb.go
	rm -f ./proto/registros.pb.go
	docker stop clienteingenieros
	docker rm clienteingenieros
	docker rmi cliente-ingenieros:latest
	