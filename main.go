package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// BASICOS
const (
	datosExcel         = "../DATA/DB_POS_290720211130.xlsx"
	hojaDatosExcel     = "Sheet1"
	formatoNombre      = "DB_POS_%s_PLA_%s_290720211130"
	nombreplantilla    = "../../../_PLANTILLAS/DATOS_BASICOS.xlsx"
	hojaPlantillaExcel = "PLANTILLA"
	celdaPrograma      = 25
	celdaPlanEstudio   = 26
	filaInicial        = 8
	celdaInicial       = 2
	celdaFinal         = 24
)

// DIRECCIONES
// const (
// 	datosExcel         = "../DATA/DD_POS_290720211130.xlsx"
// 	hojaDatosExcel     = "Sheet1"
// 	formatoNombre      = "DD_POS_%s_PLA_%s_290720211130"
// 	nombreplantilla    = "../../../_PLANTILLAS/DATOS_DIRECCION.xlsx"
// 	hojaPlantillaExcel = "Hoja1"
// 	celdaPrograma      = 14
// 	celdaPlanEstudio   = 15
// 	filaInicial        = 7
// 	celdaInicial       = 2
// 	celdaFinal         = 13
// )

func main() {
	f, err := excelize.OpenFile(datosExcel)
	if err != nil {
		fmt.Println(err)
		return
	}

	//rows := f.Sheet[hojaDatosExcel].Rows
	rows, err := f.GetRows(hojaDatosExcel)
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println(rows)
	programaActual := rows[0][celdaPrograma]
	planActual := rows[0][celdaPlanEstudio]

	// verificar antes de iniciar todo
	fmt.Println(programaActual)
	fmt.Println(planActual)
	//return
	nombreArchivoActual := fmt.Sprintf(fmt.Sprintf("%s.xlsx", formatoNombre), programaActual, planActual)
	nombreArchivoActualTXT := fmt.Sprintf(fmt.Sprintf("%s.txt", formatoNombre), programaActual, planActual)
	fTxt, err := os.Create(nombreArchivoActualTXT)
	archivoActual, err := excelize.OpenFile(nombreplantilla)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err != nil {
		fmt.Println(err)
		return
	}
	c := filaInicial
	archivos := 0
	limitRows := len(rows) - 1

	for i, row := range rows {
		if row[celdaPrograma] != "" {
			if row[celdaPrograma] != programaActual {
				c = filaInicial
				err := archivoActual.SaveAs(nombreArchivoActual)
				if err != nil {
					fmt.Println(err)
					return
				}
				fTxt.Close()
				archivos = archivos + 2
				fmt.Println(nombreArchivoActual)
				fmt.Println(nombreArchivoActualTXT)
				programaActual := row[celdaPrograma]
				planActual := row[celdaPlanEstudio]
				nombreArchivoActual = fmt.Sprintf(fmt.Sprintf("%s.xlsx", formatoNombre), programaActual, planActual)
				nombreArchivoActualTXT = fmt.Sprintf(fmt.Sprintf("%s.txt", formatoNombre), programaActual, planActual)
				archivoActual, err = excelize.OpenFile(nombreplantilla)
				fTxt, err = os.Create(nombreArchivoActualTXT)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			if row[celdaPlanEstudio] != planActual {
				c = filaInicial
				archivoActual.SaveAs(nombreArchivoActual)
				if err != nil {
					fmt.Println(err)
					return
				}
				fTxt.Close()
				archivos = archivos + 2
				fmt.Println(nombreArchivoActual)
				fmt.Println(nombreArchivoActualTXT)
				programaActual := row[celdaPrograma]
				planActual := row[celdaPlanEstudio]
				nombreArchivoActual = fmt.Sprintf(fmt.Sprintf("%s.xlsx", formatoNombre), programaActual, planActual)
				nombreArchivoActualTXT = fmt.Sprintf(fmt.Sprintf("%s.txt", formatoNombre), programaActual, planActual)
				archivoActual, err = excelize.OpenFile(nombreplantilla)
				fTxt, err = os.Create(nombreArchivoActualTXT)
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			if c != filaInicial {
				fTxt.WriteString("\n")
			}

			for ind, colCell := range row {
				if ind <= celdaFinal {
					ref, _ := excelize.CoordinatesToCellName(ind+celdaInicial, c)
					err := archivoActual.SetCellStr(hojaPlantillaExcel, ref, strings.TrimSpace(colCell))
					if err != nil {
						fmt.Println(err)
						return
					}
					fTxt.WriteString(fmt.Sprintf("%s\t", strings.TrimSpace(colCell)))
				}
				ind++
			}
			c++
			if i == limitRows {
				archivoActual.SaveAs(nombreArchivoActual)
				if err != nil {
					fmt.Println(err)
					return
				}
				fTxt.Close()
				archivos = archivos + 2
				fmt.Println(nombreArchivoActual)
				fmt.Println(nombreArchivoActualTXT)
				fmt.Println("Se generaron: ", archivos, "Archivos")
			}
		}
	}

}