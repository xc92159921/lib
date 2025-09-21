package seeds

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Folder   string
	Overflow int
}

var defaultConfig = Config{
	Folder:   "seeds",
	Overflow: 100,
}

func Init(folder string, overflow int) {
	defaultConfig.Folder = folder
	defaultConfig.Overflow = overflow
}

func ReadCsvFile[T any](fileName string, handleRow func(row map[string]string) T, pushData func(bulk []T)) {
	absPath, _ := filepath.Abs(defaultConfig.Folder + "/" + fileName)
	fmt.Println(absPath)

	file, err := os.Open(absPath)
	if err != nil {
		log.Println("Не удалось открыть файл:", absPath)
		return
	}
	defer file.Close()

	data, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Println("Не удалось прочитать файл:", absPath)
		return
	}

	var header = data[0]
	var bulk []T
	for _, values := range data[1:] {
		var row = make(map[string]string)

		for index, collumn := range header {
			if len(values) > index {
				value := values[index]

				if strings.TrimSpace(value) != "" {
					row[collumn] = value
				}
			}
		}

		bulk = append(bulk, handleRow(row))

		if len(bulk) > defaultConfig.Overflow {
			pushData(bulk)
			bulk = []T{}
		}

	}

	pushData(bulk)
}
