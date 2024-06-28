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

type comandos func() *registros.Mensaje

func Agregar() *registros.Mensaje {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ingrese el Nombre del Sector:")
	sector, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
	}

	fmt.Println("Ingrese el Nombre de la Base a Agregar: ")
	base, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
	}

	var valor int
	fmt.Println("Ingrese la Cantidad de Enemigos en la Base: ")
	_, err = fmt.Scanf("%d", &valor)
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
	}

	valor32 := int32(valor)

	mensaje := &registros.Mensaje{
		Comando:      1,
		NombreSector: sector,
		NombreBase:   base,
		Opcional: &registros.Mensaje_Valor{
			Valor: valor32,
		},
	}

	return mensaje
}

func Renombrar() *registros.Mensaje {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ingrese el Nombre del Sector:")
	sector, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
	}

	fmt.Println("Ingrese el Nombre de la Base a Modificar: ")
	base, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
	}

	fmt.Println("Ingrese el Nuevo Nombre para la Base: ")
	nuevo_nombre, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
	}

	mensaje := &registros.Mensaje{
		Comando:      2,
		NombreSector: sector,
		NombreBase:   base,
		Opcional: &registros.Mensaje_NuevoNombre{
			NuevoNombre: nuevo_nombre,
		},
	}

	return mensaje
}

func Actualizar() *registros.Mensaje {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ingrese el Nombre del Sector:")
	sector, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
	}

	fmt.Println("Ingrese el Nombre de la Base a Actualizar: ")
	base, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
	}

	var valor int
	fmt.Println("Ingrese la Cantidad de Enemigos en la Base: ")
	_, err = fmt.Scanf("%d", &valor)
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
	}

	valor32 := int32(valor)

	mensaje := &registros.Mensaje{
		Comando:      3,
		NombreSector: sector,
		NombreBase:   base,
		Opcional: &registros.Mensaje_Valor{
			Valor: valor32,
		},
	}

	return mensaje
}

func Borrar() *registros.Mensaje {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Ingrese el Nombre del Sector:")
	sector, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
	}

	fmt.Println("Ingrese el Nombre de la Base a Eliminar: ")
	base, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
	}

	mensaje := &registros.Mensaje{
		Comando:      4,
		NombreSector: sector,
		NombreBase:   base,
	}

	return mensaje
}

func main() {

	acciones := []comandos{Agregar, Renombrar, Actualizar, Borrar}

	var modo int

	fmt.Println("Eliga el comando a ingresar\n1 - Agregar Base\n2 - Renombrar Base\n3 - Actualizar Valor de una Base\n4 - Borrar Base\n\n0 - Cerrar Consola")
	_, err := fmt.Scanf("%d", &modo)
	if err != nil {
		fmt.Println("Error al leer la entrada:", err)
		return
	}

	if modo == 0 {

		return
	}

	mensaje := acciones[modo-1]()

	/// Manda el Mensaje
	host := "dist021:50051"
	///host := ":50051"
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
		case *registros.Respuesta_Vector:
			log.Printf("Respuesta del Servidor: %v", respuesta2.GetVector())

		default:
			log.Println("Respuesta inesperada")
		}

	case *registros.Respuesta_Vector:

		log.Printf("Respuesta de vector: %v", respuesta1.GetVector())

	default:
		log.Println("Respuesta inesperada")
	}

}
