package main

import (
	"bufio"
	"context"
	"fmt"
	registros "go-container/proto"
	"log"
	"os"

	"google.golang.org/grpc"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ingrese el Nombre del Sector: ")
	sector, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
	}

	fmt.Println("Ingrese el Nombre de la Base: ")
	base, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
	}

	mensaje := &registros.Mensaje{
		Comando:      5,
		NombreSector: sector,
		NombreBase:   base,
	}

	/// Manda el Mensaje
	host := "dist021:50051"
	/// host := "localhost:50051"

	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se pudo conectar: %v", err)
	}
	defer conn.Close()

	cliente1 := registros.NewSolicitud_InfoClient(conn)

	fmt.Println("Enviando Mensaje a :", host)

	respuesta1, err := cliente1.Solicitar_Info(context.Background(), mensaje)
	if err != nil {
		log.Fatalf("Error al llamar a Solicitar_Info: %v", err)
	}

	switch respuesta1.Opciones.(type) {
	case *registros.Respuesta_Direccion:

		host = respuesta1.GetDireccion()
		conn, err := grpc.Dial(host, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("No se pudo conectar: %v", err)
		}
		defer conn.Close()

		cliente2 := registros.NewSolicitud_InfoClient(conn)

		fmt.Println("Enviando Mensaje a :", host)

		respuesta2, err := cliente2.Solicitar_Info(context.Background(), mensaje)
		if err != nil {
			log.Fatalf("Error al llamar a Solicitar_Info: %v", err)
		}

		switch respuesta2.Opciones.(type) {
		case *registros.Respuesta_DetalleConLista:
			log.Printf("Respuesta del Servidor: %v", respuesta2.GetDetalleConLista())

		default:
			log.Println("Respuesta inesperada")
		}

	case *registros.Respuesta_DetalleConLista:

		log.Printf("Respuesta de vector: %v", respuesta1.GetDetalleConLista())

	default:
		log.Println("Respuesta inesperada")
	}

}
