package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/northbright/pathhelper"
	"github.com/shchnmz/ming"
)

// Config represents the app settings.
type Config struct {
	RedisServer                 string `json:"redis_server"`
	RedisPassword               string `json:"redis_password"`
	InputCSVNameColumnIndex     int    `json:"input_csv_name_column_index"`
	InputCSVPhoneNumColumnIndex int    `json:"input_csv_phone_num_column_index"`
}

const (
	inputCSV = "input.csv"
	dumpCSV  = "left-students.csv"
)

var (
	config Config
	db     ming.DB
)

func main() {
	var err error

	defer func() {
		if err != nil {
			log.Printf("%v", err)
		}
	}()

	if err = loadConfig("config.json", &config); err != nil {
		return
	}

	if err = findLeftStudents(config.RedisServer, config.RedisPassword); err != nil {
		return
	}
}

func loadConfig(file string, config *Config) error {
	var (
		err        error
		buf        []byte
		currentDir string
	)

	currentDir, _ = pathhelper.GetCurrentExecDir()
	file = path.Join(currentDir, file)

	// Load Conifg
	if buf, err = ioutil.ReadFile(file); err != nil {
		return err
	}

	return json.Unmarshal(buf, config)
}

// findLeftStudents lists students in ming800.
func findLeftStudents(redisServer, redisPassword string) error {
	var (
		err     error
		records [][]string
	)

	// Read input CSV
	currentDir, _ := pathhelper.GetCurrentExecDir()
	file := path.Join(currentDir, inputCSV)
	inputFile, err := os.Open(file)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	csvReader := csv.NewReader(inputFile)
	csvReader.Comma = ','
	records, err = csvReader.ReadAll()
	if err != nil {
		return err
	}

	l := len(records)
	m := map[string]bool{}
	for i := 1; i < l; i++ {
		record := records[i]
		name := ""
		phoneNum := ""
		n := len(record)

		if config.InputCSVNameColumnIndex >= n {
			continue
		}
		if config.InputCSVPhoneNumColumnIndex >= n {
			continue
		}
		name = records[i][config.InputCSVNameColumnIndex]
		phoneNum = records[i][config.InputCSVPhoneNumColumnIndex]

		// Trim TAB
		phoneNum = strings.Trim(phoneNum, "\x09")
		k := fmt.Sprintf("%v:%v", name, phoneNum)

		m[k] = true
	}

	db := ming.DB{config.RedisServer, config.RedisPassword}
	students, err := db.GetAllStudents()
	if err != nil {
		return err
	}

	// dump UTF8BOM CSV
	file = path.Join(currentDir, dumpCSV)

	outputFile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	UTF8BOM := [3]byte{0xEF, 0xBB, 0xBF}
	if _, err = outputFile.Write(UTF8BOM[0:3]); err != nil {
		return err
	}

	csvWriter := csv.NewWriter(outputFile)
	csvWriter.Comma = ','

	// Find left students
	leftNum := 0
	for _, student := range students {
		if _, ok := m[student]; ok {
			continue
		}
		leftNum++
		arr := strings.Split(student, ":")
		if len(arr) < 2 {
			continue
		}
		name := arr[0]
		phoneNum := arr[1]

		classes, err := db.GetClassesByNameAndPhoneNum(name, phoneNum)
		if err != nil {
			return err
		}

		for _, str := range classes {
			arr := strings.Split(str, ":")
			if len(arr) < 3 {
				continue
			}
			campus := arr[0]
			category := arr[1]
			class := arr[2]

			teachers, err := db.GetTeachersOfClass(campus, category, class)
			if err != nil {
				return err
			}
			for _, teacher := range teachers {
				if err = csvWriter.Write([]string{name, phoneNum, class, teacher}); err != nil {
					return err
				}
			}
		}
	}

	csvWriter.Flush()
	fmt.Printf("Left students number: %v\n", leftNum)
	fmt.Printf("All left student - class - teacher record are dumped to left-students.csv\n")

	return nil
}
