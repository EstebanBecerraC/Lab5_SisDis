package main

import (
	"context"
	"fmt"
	registros "go-container/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	registros.UnimplementedSolicitud_InfoServer
}

func (s *server) Solicitar_Info(ctx context.Context, mensaje *registros.Mensaje) (*registros.Respuesta, error) {

	fulcrum_host := "dist022:50052"
	/// fulcrum_host := "localhost:50052"
	fmt.Println("Se recibio un mensaje. Respondiendo con Direccion: ", fulcrum_host)

	return &registros.Respuesta{
		Opciones: &registros.Respuesta_Direccion{
			Direccion: fulcrum_host,
		},
	}, nil
}

func main() {

	puerto := ":50051"
	lis, err := net.Listen("tcp", puerto)
	if err != nil {
		log.Fatalf("No se pudo iniciar el servidor: %v", err)
	}

	//se inicia la conexion grpc
	grpcServer := grpc.NewServer()
	registros.RegisterSolicitud_InfoServer(grpcServer, &server{})
	log.Println("Servidor iniciado en puerto", puerto)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("El servidor se detuvo: %v", err)
	}

}
