package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	tipoPlantilla string
	configuracion map[string]interface{}
	rutaSalida    string
	formatoSalida string
	prefijoSalida string
)

func main() {
	inicializar()
}

func inicializar() {
	// inicializar cli
	cli{}.inicializar()

	// validar ext
	tipo := filepath.Ext(archivoDatos)
	if tipo != ".ctrArchivo" && tipo != ".csv" {
		log.Fatalln("El archivo de origen de datos debe ser: .csv o .ctrArchivo")
	}

	// cargar configuracion
	js, err := os.Open(archivoConfig)
	if err != nil {
		log.Fatalf("Error al cargar configuracion, %v", err)
	}
	defer func() { _ = js.Close() }()
	datosByte, _ := ioutil.ReadAll(js)
	err = json.Unmarshal(datosByte, &configuracion)
	if err != nil {
		log.Fatalf("Error al decodificar configuracion, %v", err)
	}

	// obtener el tipo de plantilla
	if genDB {
		tipoPlantilla = "DB"
	}
	if genDD {
		tipoPlantilla = "DD"
	}
	if genSR {
		tipoPlantilla = "SR"
	}
	if genPR {
		tipoPlantilla = "PR"
	}
	if genSE {
		tipoPlantilla = "SE"
	}
	if genRA {
		tipoPlantilla = "RA"
	}
	if genTA {
		tipoPlantilla = "TA"
	}
	if genIN {
		tipoPlantilla = "IN"
	}
	if genST {
		tipoPlantilla = "ST"
	}

	// crear carpetas de salida
	rutaSalida = fmt.Sprintf("%s/%s", configuracion["ruta_exportacion"], configuracion[tipoPlantilla].(map[string]interface{})["alias"])
	err = os.MkdirAll(rutaSalida, os.ModePerm)
	if err != nil {
		log.Fatalf("Error al crear directorios de salida, %v\n", err)
	}

	// prefijo de salida
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Ingrese prefijo para los archivos de salida: ")
	prefijoSalida, _ = reader.ReadString('\n')
	prefijoSalida = strings.TrimSpace(prefijoSalida)

	fmt.Print("Ingrese el segundo prefijo para los archivos de salida: ")
	prefijoExp, _ := reader.ReadString('\n')
	prefijoExp = strings.TrimSpace(prefijoExp)

	// formato de salida
	fechaActual := time.Now()
	formatoSalida = fmt.Sprintf(configuracion["formato"].(string),
		configuracion[tipoPlantilla].(map[string]interface{})["alias"],
		prefijoSalida, prefijoExp, "%s", "%s", fechaActual.Format("020120061504"))

	x := ctrArchivo{}
	x.iniciar(tipo)
}
