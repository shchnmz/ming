package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"

	"github.com/garyburd/redigo/redis"
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
	var (
		err                    error
		buf                    []byte
		currentDir, configFile string
	)

	defer func() {
		if err != nil {
			fmt.Printf("%v", err)
		}
	}()

	currentDir, _ = pathhelper.GetCurrentExecDir()
	configFile = path.Join(currentDir, "config.json")

	// Load Conifg
	if buf, err = ioutil.ReadFile(configFile); err != nil {
		err = fmt.Errorf("load config file error: %v", err)
		return
	}

	if err = json.Unmarshal(buf, &config); err != nil {
		err = fmt.Errorf("parse config err: %v", err)
		return
	}

	err = FindClasses(config.RedisServer, config.RedisPassword)
}

func FindClasses(redisServer, redisPassword string) error {
	var (
		err error
	)

	conn, err := redishelper.GetRedisConn(redisServer, redisPassword)
	if err != nil {
		return err
	}
	defer conn.Close()

	k := ":classes"
	classes, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
	if err != nil {
		return err
	}

	for _, class := range classes {
		fmt.Printf("%v\n", class)
	}

	return nil
}
