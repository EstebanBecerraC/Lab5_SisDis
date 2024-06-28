package main

import (
	"bufio"
	"context"
	"fmt"
	registros "go-container/proto"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
)

var ( //variables globales
	sectores map[string]map[string]int
	relojes  map[string][2][3]int
	/// ticker         = time.NewTicker(30 * time.Second)
)

type server struct {
	listaSectores map[string]map[string]int
	listaRelojes  map[string][2][3]int
	registros.UnimplementedSolicitud_InfoServer
}

///func Consistencia() {
///
///	for range ticker.C {
///		/// Realizar consistencia entre los sv
///	}
///
///}

type comandos func(context.Context, *registros.Mensaje) (*registros.Respuesta, error)

// / Agregar una nueva base a un sector
func (s *server) Agregar(ctx context.Context, mensaje *registros.Mensaje) (*registros.Respuesta, error) {

	/// Se guardan los datos en variables
	sector := mensaje.NombreSector
	base := mensaje.NombreBase
	var valor32 int32

	switch msg := mensaje.Opcional.(type) {
	case *registros.Mensaje_Valor:
		valor32 = msg.Valor
	default:
	}

	valor := int(valor32)

	/// Se agregan los datos en el sector
	if s.listaSectores[sector] == nil {
		/// Si no esta el sector, se crea
		s.listaSectores[sector] = make(map[string]int)
	}

	s.listaSectores[sector][base] = valor

	/// Se modifica el relog del sector
	if _, exists := s.listaRelojes[sector]; !exists {
		/// Se crea un relog en caso de no existir
		s.listaRelojes[sector] = [2][3]int{
			{0, 0, 0},
			{0, 0, 0},
		}
	}

	relojes := s.listaRelojes[sector]
	relojes[0][0]++
	s.listaRelojes[sector] = relojes

	/// Se abre el documento .txt del sector para escribir la nueva base
	path := "Sectores/"
	formato := ".txt"
	nombreDocumento := path + sector + formato
	escribir := sector + " " + base + " " + strconv.Itoa(valor) + "\n"

	/// Abrir el archivo para escritura y crear si no existe
	archivo1, err := os.OpenFile(nombreDocumento, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error al abrir o crear el archivo:", err)
	}

	defer archivo1.Close()

	/// Se escribe
	_, err = archivo1.WriteString(escribir)
	if err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
	}

	/// Abrir el log para escribir el cambio
	escribirLog := "AgregarBase " + sector + " " + base + " " + strconv.Itoa(valor) + "\n"
	fmt.Println(escribirLog)

	nombreLog := "LogFulcrum1" + formato
	archivoLog, err := os.OpenFile(nombreLog, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error al abrir o crear el archivo:", err)
	}

	defer archivoLog.Close()

	_, err = archivoLog.WriteString(escribirLog)
	if err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
	}

	/// Se arma la respuesta del servidor
	array := s.listaRelojes[sector][0]

	slice := make([]int32, 3)
	for i, v := range array {
		slice[i] = int32(v)
	}

	respuesta := &registros.Vector{
		Vector: slice,
	}

	return &registros.Respuesta{
		Opciones: &registros.Respuesta_Vector{
			Vector: respuesta,
		},
	}, nil
}

// / Renombrar una base
func (s *server) Renombrar(ctx context.Context, mensaje *registros.Mensaje) (*registros.Respuesta, error) {

	/// Se guardan los datos en variables
	sector := mensaje.NombreSector
	base := mensaje.NombreBase
	var nuevo_nombre string

	switch msg := mensaje.Opcional.(type) {
	case *registros.Mensaje_NuevoNombre:
		nuevo_nombre = msg.NuevoNombre
	default:
	}

	tmp_valor := s.listaSectores[sector][base]

	s.listaSectores[sector][nuevo_nombre] = tmp_valor

	delete(s.listaSectores[sector], base)

	/// Modificar el documento
	path := "Sectores/"
	formato := ".txt"
	nombreDocumento := path + sector + formato
	linea_vieja := sector + " " + base + " " + strconv.Itoa(tmp_valor) + "\n"
	linea_nueva := sector + " " + nuevo_nombre + " " + strconv.Itoa(tmp_valor) + "\n"

	archivo1, err := os.OpenFile(nombreDocumento, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
	}

	defer archivo1.Close()

	/// Se busca la linea para reemplazarla
	scanner := bufio.NewScanner(archivo1)
	var lineas []string
	cambiado := false

	for scanner.Scan() {
		linea := scanner.Text()
		if linea == linea_vieja && !cambiado {
			lineas = append(lineas, linea_nueva)
			cambiado = true
		} else {
			lineas = append(lineas, linea)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
	}

	/// Se vuelve a escribir el archivo
	if err := archivo1.Truncate(0); err != nil {
		fmt.Println("Error al truncar el archivo:", err)
	}

	if _, err := archivo1.Seek(0, 0); err != nil {
		fmt.Println("Error al situar el offset al inicio del archivo:", err)
	}

	writer := bufio.NewWriter(archivo1)
	for _, line := range lineas {
		_, err := fmt.Fprintln(writer, line)
		if err != nil {
			fmt.Println("Error al escribir en el archivo:", err)
		}
	}

	// Vaciar el buffer del writer
	if err := writer.Flush(); err != nil {
		fmt.Println("Error al vaciar el buffer:", err)
	}

	/// Abrir el log para escribir el cambio
	escribirLog := "RenombrarBase " + sector + " " + base + " " + nuevo_nombre + "\n"
	fmt.Println(escribirLog)

	nombreLog := "LogFulcrum1" + formato
	archivoLog, err := os.OpenFile(nombreLog, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error al abrir o crear el archivo:", err)
	}

	defer archivoLog.Close()

	_, err = archivoLog.WriteString(escribirLog)
	if err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
	}

	/// Se actualiza el reloj
	reloj := s.listaRelojes[sector]
	reloj[0][0]++
	s.listaRelojes[sector] = reloj

	array := s.listaRelojes[sector][0]

	slice := make([]int32, 3)
	for i, v := range array {
		slice[i] = int32(v)
	}

	respuesta := &registros.Vector{
		Vector: slice,
	}

	return &registros.Respuesta{
		Opciones: &registros.Respuesta_Vector{
			Vector: respuesta,
		},
	}, nil
}

// / Actualizar cantidad de enemigos
func (s *server) Actualizar(ctx context.Context, mensaje *registros.Mensaje) (*registros.Respuesta, error) {

	/// Se guardan los datos en variables
	sector := mensaje.NombreSector
	base := mensaje.NombreBase
	var valor32 int32

	switch msg := mensaje.Opcional.(type) {
	case *registros.Mensaje_Valor:
		valor32 = msg.Valor
	default:
	}

	valor_nuevo := int(valor32)

	valor_viejo := s.listaSectores[sector][base]

	s.listaSectores[sector][base] = valor_nuevo

	/// Modificar el documento
	path := "Sectores/"
	formato := ".txt"
	nombreDocumento := path + sector + formato
	linea_vieja := sector + " " + base + " " + strconv.Itoa(valor_viejo) + "\n"
	linea_nueva := sector + " " + base + " " + strconv.Itoa(valor_nuevo) + "\n"

	archivo1, err := os.OpenFile(nombreDocumento, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
	}

	defer archivo1.Close()

	/// Se busca la linea para reemplazarla
	scanner := bufio.NewScanner(archivo1)
	var lineas []string
	cambiado := false

	for scanner.Scan() {
		linea := scanner.Text()
		if linea == linea_vieja && !cambiado {
			lineas = append(lineas, linea_nueva)
			cambiado = true
		} else {
			lineas = append(lineas, linea)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
	}

	/// Se vuelve a escribir el archivo
	if err := archivo1.Truncate(0); err != nil {
		fmt.Println("Error al truncar el archivo:", err)
	}

	if _, err := archivo1.Seek(0, 0); err != nil {
		fmt.Println("Error al situar el offset al inicio del archivo:", err)
	}

	writer := bufio.NewWriter(archivo1)
	for _, line := range lineas {
		_, err := fmt.Fprintln(writer, line)
		if err != nil {
			fmt.Println("Error al escribir en el archivo:", err)
		}
	}

	// Vaciar el buffer del writer
	if err := writer.Flush(); err != nil {
		fmt.Println("Error al vaciar el buffer:", err)
	}

	/// Abrir el log para escribir el cambio
	escribirLog := "ActualizarValor " + sector + " " + base + " " + strconv.Itoa(valor_nuevo) + "\n"
	fmt.Println(escribirLog)

	nombreLog := "LogFulcrum1" + formato
	archivoLog, err := os.OpenFile(nombreLog, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error al abrir o crear el archivo:", err)
	}

	defer archivoLog.Close()

	_, err = archivoLog.WriteString(escribirLog)
	if err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
	}

	/// Se actualiza el reloj
	reloj := s.listaRelojes[sector]
	reloj[0][0]++
	s.listaRelojes[sector] = reloj

	array := s.listaRelojes[sector][0]

	slice := make([]int32, 3)
	for i, v := range array {
		slice[i] = int32(v)
	}

	respuesta := &registros.Vector{
		Vector: slice,
	}

	return &registros.Respuesta{
		Opciones: &registros.Respuesta_Vector{
			Vector: respuesta,
		},
	}, nil
}

// / Borrar Base
func (s *server) Borrar(ctx context.Context, mensaje *registros.Mensaje) (*registros.Respuesta, error) {

	/// Se guardan los datos en variables
	sector := mensaje.NombreSector
	base := mensaje.NombreBase

	valor_viejo := s.listaSectores[sector][base]

	delete(s.listaSectores[sector], base)

	/// Modificar el documento
	path := "Sectores/"
	formato := ".txt"
	nombreDocumento := path + sector + formato
	linea_vieja := sector + " " + base + " " + strconv.Itoa(valor_viejo) + "\n"

	archivo1, err := os.OpenFile(nombreDocumento, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
	}

	defer archivo1.Close()

	/// Se busca la linea para reemplazarla
	scanner := bufio.NewScanner(archivo1)
	var lineas []string
	cambiado := false

	for scanner.Scan() {
		linea := scanner.Text()
		if linea == linea_vieja && !cambiado {
			cambiado = true
		} else {
			lineas = append(lineas, linea)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error al leer el archivo:", err)
	}

	/// Se vuelve a escribir el archivo
	if err := archivo1.Truncate(0); err != nil {
		fmt.Println("Error al truncar el archivo:", err)
	}

	if _, err := archivo1.Seek(0, 0); err != nil {
		fmt.Println("Error al situar el offset al inicio del archivo:", err)
	}

	writer := bufio.NewWriter(archivo1)
	for _, line := range lineas {
		_, err := fmt.Fprintln(writer, line)
		if err != nil {
			fmt.Println("Error al escribir en el archivo:", err)
		}
	}

	// Vaciar el buffer del writer
	if err := writer.Flush(); err != nil {
		fmt.Println("Error al vaciar el buffer:", err)
	}

	/// Abrir el log para escribir el cambio
	escribirLog := "BorrarBase " + sector + " " + base + "\n"
	fmt.Println(escribirLog)

	nombreLog := "LogFulcrum1" + formato
	archivoLog, err := os.OpenFile(nombreLog, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error al abrir o crear el archivo:", err)
	}

	defer archivoLog.Close()

	_, err = archivoLog.WriteString(escribirLog)
	if err != nil {
		fmt.Println("Error al escribir en el archivo:", err)
	}

	/// Se actualiza el reloj
	reloj := s.listaRelojes[sector]
	reloj[0][0]++
	s.listaRelojes[sector] = reloj

	array := s.listaRelojes[sector][0]

	slice := make([]int32, 3)
	for i, v := range array {
		slice[i] = int32(v)
	}

	respuesta := &registros.Vector{
		Vector: slice,
	}

	return &registros.Respuesta{
		Opciones: &registros.Respuesta_Vector{
			Vector: respuesta,
		},
	}, nil
}

func (s *server) GetEnemigos(ctx context.Context, mensaje *registros.Mensaje) (*registros.Respuesta, error) {

	sector := mensaje.NombreSector
	base := mensaje.NombreBase

	valor := s.listaSectores[sector][base]
	valor32 := int32(valor)

	vector := []int32{int32(s.listaRelojes[base][0][0]), int32(s.listaRelojes[base][0][1]), int32(s.listaRelojes[base][0][2])}

	res := sector + " " + base + " " + strconv.Itoa(valor) + "\n"
	fmt.Println("Mandando el Siguiente Mensaje al Comandante :")
	fmt.Println("   ", res)

	respuestaKais := &registros.Res_Kais{
		NombreSector: sector,
		NombreBase:   base,
		Valor:        valor32,
		Vector:       vector,
	}

	// Crear y retornar la respuesta que contiene el mensaje Res_Kais
	return &registros.Respuesta{
		Opciones: &registros.Respuesta_DetalleConLista{
			DetalleConLista: respuestaKais,
		},
	}, nil
}

// / funcion que recibe el mensaje
func (s *server) Solicitar_Info(ctx context.Context, mensaje *registros.Mensaje) (*registros.Respuesta, error) {

	fmt.Println("Se recibio un mensaje, Procesando ...")
	acciones := []comandos{s.Agregar, s.Renombrar, s.Actualizar, s.Borrar, s.GetEnemigos}
	modo := mensaje.Comando
	return acciones[modo-1](ctx, mensaje)

}

func main() {

	sectores = make(map[string]map[string]int) /// la primera llave es el sector, la 2da es la base y el int es la cantidad de enemigos
	relojes = make(map[string][2][3]int)       /// Es un mapa que al ussar el sector como llave, te lleva a una lista de largo 2, las cuales contienen, cada una, una lista de largo 3 en donde se guardan el valor de los vectores. La primera lista de largo 3 es el vector actual y la 2da es de la ultima vez que se coordinaron los servidores fulcrum

	/// go Consistencia()

	servidor := &server{
		listaSectores: sectores,
		listaRelojes:  relojes,
	}

	puerto := ":50051"
	lis, err := net.Listen("tcp", puerto)
	if err != nil {
		log.Fatalf("No se pudo iniciar el servidor: %v", err)
	}

	//se inicia la conexion grpc
	grpcServer := grpc.NewServer()
	registros.RegisterSolicitud_InfoServer(grpcServer, servidor)
	log.Println("Servidor iniciado en puerto", puerto)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("El servidor se detuvo: %v", err)
	}

}
