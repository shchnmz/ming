package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/gomodule/redigo/redis"
	"github.com/northbright/pathhelper"
	"github.com/northbright/redishelper"
)

// Config represents the app settings.
type Config struct {
	RedisServer   string `json:"redis_server"`
	RedisPassword string `json:"redis_password"`
}

var (
	config Config
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

	if err = findClasses(config.RedisServer, config.RedisPassword); err != nil {
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

func findClasses(redisServer, redisPassword string) error {
	var (
		err error
	)

	conn, err := redishelper.Dial(redisServer, redisPassword)
	if err != nil {
		return err
	}
	defer conn.Close()

	k := "ming::classes"
	classes, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
	if err != nil {
		return err
	}

	for _, class := range classes {
		fmt.Printf("%v\n", class)
	}

	return nil
}
