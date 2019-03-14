package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/northbright/pathhelper"
	"github.com/shchnmz/ming"
)

// Config represents app settings.
type Config struct {
	RedisServer   string `json:"redis_server"`
	RedisPassword string `json:"redis_password"`
}

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

	if err = listAllClasses(config.RedisServer, config.RedisPassword); err != nil {
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

// listAllClasses lists all classes in ming800.
func listAllClasses(redisServer, redisPassword string) error {
	var (
		err error
	)

	db = ming.DB{config.RedisServer, config.RedisPassword}
	classes, err := db.GetAllClasses()
	if err != nil {
		return err
	}

	for _, class := range classes {
		fmt.Printf("%s\n", class)
	}

	return nil
}
