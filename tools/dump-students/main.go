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
	RedisServer   string `json:"redis_server"`
	RedisPassword string `json:"redis_password"`
}

const (
	dumpCSV = "all-students.csv"
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

	if err = dumpStudents(config.RedisServer, config.RedisPassword); err != nil {
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

// dumpStudents lists students in ming800.
func dumpStudents(redisServer, redisPassword string) error {
	var (
		err     error
		records [][]string
	)

	db := ming.DB{config.RedisServer, config.RedisPassword}
	students, err := db.GetAllStudents()
	if err != nil {
		return err
	}

	for _, student := range students {
		arr := strings.Split(student, ":")
		if len(arr) < 2 {
			continue
		}
		name := arr[0]
		phoneNum := arr[1]

		// Generate CSV records
		records = append(records, []string{name, phoneNum})

		fmt.Printf("%v\n", student)
	}

	// dump UTF8BOM CSV
	currentDir, _ := pathhelper.GetCurrentExecDir()
	file := path.Join(currentDir, dumpCSV)

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	UTF8BOM := [3]byte{0xEF, 0xBB, 0xBF}
	if _, err = f.Write(UTF8BOM[0:3]); err != nil {
		return err
	}

	csvWriter := csv.NewWriter(f)
	csvWriter.Comma = ','
	if err = csvWriter.WriteAll(records); err != nil {
		return err
	}

	csvWriter.Flush()

	return nil
}
