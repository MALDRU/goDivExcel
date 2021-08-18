package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"log"
	"os"
	"strings"
)

type ctrArchivo struct {
	numeroArchivosGen int
	nombreArchivoXlsx string
	nombreArchivoTxt string
	archivoXlsx *excelize.File
	archivoTxt *os.File
	programaActual string
	planActual string
}

func (x *ctrArchivo) iniciar(tipo string) {
	// pedir datos por consola
	hojaArchivoDatos, celdaPrograma, celdaPlanEstudio := x.obtenerDatosConsola()

	// extraer filas
	filas := x.obtenerFilasArchivo(tipo, hojaArchivoDatos)

	// programa
	indxCeldaPrograma, _, err := excelize.CellNameToCoordinates(fmt.Sprintf("%s1", celdaPrograma))
	if err != nil {
		log.Fatalf("Error en la letra digitada, %v\n", err)
	}
	indxCeldaPrograma = indxCeldaPrograma - 1

	// plan de estudios
	indxCeldaPlanEstudio, _, err := excelize.CellNameToCoordinates(fmt.Sprintf("%s1", celdaPlanEstudio))
	if err != nil {
		log.Fatalf("Error en la letra digitada, %v\n", err)
	}
	indxCeldaPlanEstudio = indxCeldaPlanEstudio - 1

	hojaPlantilla := configuracion[tipoPlantilla].(map[string]interface{})["nombre_hoja"].(string)
	filaInicial := int(configuracion[tipoPlantilla].(map[string]interface{})["fila_inicial"].(float64))
	celdaInicial := int(configuracion[tipoPlantilla].(map[string]interface{})["celda_inicial"].(float64))

	celdaFinal := indxCeldaPrograma - 1

	c := filaInicial
	x.numeroArchivosGen = 0
	numFilas := len(*filas) - 1
	contFilas := 0

	for i, fila := range *filas {
		if i == 0 {
			x.programaActual = fila[indxCeldaPrograma]
			x.planActual = fila[indxCeldaPlanEstudio]

			if ignorarEncabezado {
				continue
			}
			// crear archivos
			x.crearArchivos()
		}
		contFilas++
		if fila[indxCeldaPrograma] != "" {
			if fila[indxCeldaPrograma] != x.programaActual || fila[indxCeldaPlanEstudio] != x.planActual {
				c = filaInicial
				x.programaActual = fila[indxCeldaPrograma]
				x.planActual = fila[indxCeldaPlanEstudio]
				x.guardar()
				x.crearArchivos()
			}

			if c != filaInicial {
				_, err = x.archivoTxt.WriteString("\n")
				if err != nil {
					log.Fatalf("Error al escribir en archivo txt %v", err)
				}
			}

			for ind, valorCelda := range fila {
				if ind <= celdaFinal {
					ref, _ := excelize.CoordinatesToCellName(ind+celdaInicial, c)
					err := x.archivoXlsx.SetCellStr(hojaPlantilla, ref, strings.TrimSpace(valorCelda))
					if err != nil {
						log.Fatalf("Error al escribir en celda ctrArchivo: %v", err)
					}
					_, err = x.archivoTxt.WriteString(fmt.Sprintf("%s", strings.TrimSpace(valorCelda)))
					if err != nil {
						log.Fatalf("Error al escribir en txt: %v", err)
					}
					if ind != celdaFinal {
						_, err = x.archivoTxt.WriteString("\t")
						if err != nil {
							log.Fatalf("Error al escribir en txt: %v", err)
						}
					}
				}
				ind++
			}
			c++
			if i == numFilas {
				x.guardar()
			}
		}
	}
	log.Printf("Se generaron: %d archivos, con un total de: %d filas", x.numeroArchivosGen, contFilas)
}

func (*ctrArchivo) abir(rutaArchivo string) (*excelize.File, error) {
	f, err := excelize.OpenFile(rutaArchivo)
	return f, err
}

func (x *ctrArchivo) obtenerDatosConsola () (string, string, string) {
	// nombre de la hoja
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese el nombre de la hoja (Sheet1): ")
	hojaArchivoDatos, _ := reader.ReadString('\n')
	hojaArchivoDatos = strings.TrimSpace(hojaArchivoDatos)
	if hojaArchivoDatos == "" {
		hojaArchivoDatos = "Sheet1"
	}

	// columna del programa y del plan de estudios
	fmt.Print("Ingrese la letra de la columna que contiene el programa: ")
	celdaPrograma, _ := reader.ReadString('\n')
	celdaPrograma = strings.TrimSpace(celdaPrograma)


	fmt.Print("Ingrese la letra de la columna que contiene el plan de estudios: ")
	celdaPlanEstudio, _ := reader.ReadString('\n')
	celdaPlanEstudio = strings.TrimSpace(celdaPlanEstudio)

	return hojaArchivoDatos, celdaPrograma, celdaPlanEstudio
}

func (x *ctrArchivo) obtenerFilasArchivo (tipo, hojaArchivoDatos string) (*[][]string) {
	var (
		filas = make([][]string, 0)
	)

	if tipo == ".csv" {
		archivoCsv, err := os.Open(archivoDatos)
		if err != nil {
			log.Fatalf("Error al leer el archivo de datos, %v\n", err)
		}
		defer func() {_ = archivoCsv.Close()}()

		filas, err = csv.NewReader(archivoCsv).ReadAll()
		if err != nil {
			log.Fatalf("Error al extraer las filas del archivo, %v\n", err)
		}

	} else {
		archivoXlsx, err := x.abir(archivoDatos)
		if err != nil {
			log.Fatalf("Error al leer el archivo de datos, %v\n", err)
		}
		// obtener filas
		filas, err = archivoXlsx.GetRows(hojaArchivoDatos)
		if err != nil {
			log.Fatalf("Error al extraer las filas del archivo, %v\n", err)
		}
	}

	if len(filas) <= 0 {
		log.Fatal("El archivo de datos esta vacio")
	}

	return &filas
}

func (x *ctrArchivo) generarNombresArchivos (programaActual, planActual string) (string, string) {
	return fmt.Sprintf(fmt.Sprintf("%s.xlsx", formatoSalida), programaActual, planActual), fmt.Sprintf(fmt.Sprintf("%s.txt", formatoSalida), programaActual, planActual)
}

func (x *ctrArchivo) guardar() {
	if x.nombreArchivoXlsx == "" {
		return
	}
	err := x.archivoXlsx.SaveAs(fmt.Sprintf("%s/%s", rutaSalida, x.nombreArchivoXlsx))
	if err != nil {
		log.Fatalf("Error al guardar archivo XLSX: %v", err)
	}

	// verificar si esta vacio
	dtll, err := x.archivoTxt.Stat()
	if err != nil {
		log.Fatalf("Error al obtener informacion de archivo TXT: %v", err)
	}

	if dtll.Size() == 0 {
		log.Fatalf("Error: No relleno datos en archivo: %s", x.nombreArchivoXlsx)
	}

	err = x.archivoTxt.Close()
	if err != nil {
		log.Fatalf("Error al cerrar archivo TXT: %v", err)
	}

	x.numeroArchivosGen = x.numeroArchivosGen + 2
	log.Println(x.nombreArchivoXlsx)
	log.Println(x.nombreArchivoTxt)
}

func (x *ctrArchivo) crearArchivos() {
	var err error
	x.nombreArchivoXlsx, x.nombreArchivoTxt = x.generarNombresArchivos(x.programaActual, x.planActual)
	x.archivoXlsx, err = excelize.OpenFile(fmt.Sprintf("%s/%s", configuracion["ruta_plantillas"], configuracion[tipoPlantilla].(map[string]interface{})["plantilla"].(string)))
	if err != nil {
		log.Fatalf("Error al crear archivo excel, %v\n", err)
	}
	x.archivoTxt, err = os.Create(fmt.Sprintf("%s/%s", rutaSalida, x.nombreArchivoTxt))
	if err != nil {
		log.Fatalf("Error al crear archivo txt, %v\n", err)
	}
}