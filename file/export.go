package file

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func ExportToCsv(inputPath string, xData []float64, tData []float64) (err error) {
	inputFile := filepath.Base(inputPath)
	extension := filepath.Ext(inputFile)
	fileName := inputFile[0 : len(inputFile)-len(extension)]
	fileName += ".csv"

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var csvRow []string

	for i := range xData {
		csvRow = nil
		csvRow = append(csvRow, strconv.FormatFloat(xData[i], 'f', 8, 64))
		csvRow = append(csvRow, strconv.FormatFloat(tData[i], 'f', 8, 64))
		if err = writer.Write(csvRow); err != nil {
			log.Fatal(err)
		}
	}

	if err := writer.Error(); err != nil {
		log.Fatalln("error writing csv:", err)
	}
	return nil
}
