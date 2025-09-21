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

// ReadCsvFile - читает csv файл и при завершении вызывает push
//
//	seeds.ReadCsvFile("users.csv", func(row map[string]string) *ent.UserCreate {
//		model := db.Client.User.Create().SetMap(row)
//		ID, _ := strconv.Atoi(row["id"])
//		model.SetID(ID)
//		return model
//	}, func(bulk []*ent.UserCreate) {
//
//		db.Client.User.CreateBulk(bulk...).OnConflictColumns(user.FieldID).UpdateNewValues().ExecX(ctx)
//	})
//
//	seeds.ReadCsvFile("courses.csv", func(row map[string]string) *ent.CourseCreate {
//		model := db.Client.Course.Create().SetMap(row)
//		ID, _ := strconv.Atoi(row["id"])
//		model.SetID(ID)
//		return model
//	}, func(bulk []*ent.CourseCreate) {
//
//		db.Client.Course.CreateBulk(bulk...).OnConflictColumns(course.FieldID).UpdateNewValues().ExecX(ctx)
//	})
//
//	seeds.ReadCsvFile("promocodes.csv", func(row map[string]string) *ent.PromocodeCreate {
//		model := db.Client.Promocode.Create().SetMap(row)
//		ID, _ := strconv.Atoi(row["id"])
//		model.SetID(ID)
//		CourseId, _ := strconv.Atoi(row["promocodeCourse"])
//		if CourseId != 0 {
//			model.SetCourseID(CourseId)
//		}
//		return model
//	}, func(bulk []*ent.PromocodeCreate) {
//
//		db.Client.Promocode.CreateBulk(bulk...).OnConflictColumns(promocode.FieldID).UpdateNewValues().ExecX(ctx)
//	})
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
