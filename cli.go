package main

import (
	"flag"
	"log"
)

type cli struct {
}

var (
	archivoDatos string
	archivoConfig string
	ignorarEncabezado bool
	genDB bool
	genDD bool
	genSR bool
	genPR bool
)

func (cli) inicializar() {
	flag.StringVar(&archivoDatos, "archivo", "", "Archivo donde se extraeran los datos (.csv | .ctrArchivo)")
	flag.StringVar(&archivoConfig, "config", "config.json", "Archivo de configuracion")
	flag.BoolVar(&ignorarEncabezado, "sinEncabezado", false, "Ignora la primera fila del archivo")
	flag.BoolVar(&genDB, "db", false, "Genera plantillas de datos basicos")
	flag.BoolVar(&genDD, "dd", false, "Genera plantillas de datos direccion")
	flag.BoolVar(&genSR, "sr", false, "Genera plantillas de sesiones de registro")
	flag.BoolVar(&genPR, "pr", false, "Genera plantillas de progresiones")
	flag.Parse()
	if archivoDatos == "" || archivoConfig == "" || !(genDB || genDD || genSR || genPR) {
		log.Fatalln("Debe indicar el archivo de origen de datos, el de configuracion y minimo una plantilla a generar. Mas informacion ingrese -help")
	}
}