package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"

	"github.com/northbright/pathhelper"
	"github.com/shchnmz/ming"
)

type Config struct {
	ServerURL     string `json:"server_url"`
	Company       string `json:"company"`
	User          string `json:"user"`
	Password      string `json:"password"`
	RedisServer   string `json:"redis_server"`
	RedisPassword string `json:"redis_password"`
}

func main() {
	var (
		err                    error
		buf                    []byte
		currentDir, configFile string
		config                 Config
	)

	defer func() {
		if err != nil {
			log.Printf("error: %v", err)
		}
	}()

	currentDir, _ = pathhelper.GetCurrentExecDir()
	configFile = path.Join(currentDir, "config.json")

	// Load Conifg
	if buf, err = ioutil.ReadFile(configFile); err != nil {
		return
	}

	if err = json.Unmarshal(buf, &config); err != nil {
		return
	}

	err = ming.Ming2Redis(config.ServerURL, config.Company, config.User, config.Password, config.RedisServer, config.RedisPassword)
}
