package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/northbright/pathhelper"
	"github.com/shchnmz/ming"
)

// Config represents the app settings.
type Config struct {
	RedisServer   string `json:"redis_server"`
	RedisPassword string `json:"redis_password"`
}

const (
	dumpJSON = "class-period.json"
	dumpCSV  = "class-period-utf8-bom.csv"
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

	if err = dumpClassesPeriods(config.RedisServer, config.RedisPassword); err != nil {
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

// dumpClassesPeriods lists class- period relationships in ming800.
func dumpClassesPeriods(redisServer, redisPassword string) error {
	var (
		err     error
		records [][]string
	)

	db := ming.DB{config.RedisServer, config.RedisPassword}
	m, err := db.GetClassesPeriods()
	if err != nil {
		return err
	}

	for k, v := range m {
		// Generate CSV records
		records = append(records, []string{k, v})

		fmt.Printf("%v -> %v\n", k, v)
	}

	// dump JSON
	jsonData, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}

	currentDir, _ := pathhelper.GetCurrentExecDir()
	file := path.Join(currentDir, dumpJSON)

	if err = ioutil.WriteFile(file, jsonData, os.ModePerm); err != nil {
		return err
	}

	fmt.Printf("====================\nAll class-period has been dump to:\n%v\n%v\n", dumpJSON, dumpCSV)

	// dump UTF8BOM CSV
	file = path.Join(currentDir, dumpCSV)

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
	csvWriter.Comma = ';'
	if err = csvWriter.WriteAll(records); err != nil {
		return err
	}

	csvWriter.Flush()

	return nil
}
