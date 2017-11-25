package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

	if err = json.Unmarshal(buf, &config); err != nil {
		return err
	}

	return nil
}

// listAllClasses lists all classes in ming800.
func listAllClasses(redisServer, redisPassword string) error {
	var (
		err error
	)

	conn, err := redishelper.GetRedisConn(redisServer, redisPassword)
	if err != nil {
		return err
	}
	defer conn.Close()

	k := "ming:campuses"
	campuses, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
	if err != nil {
		return err
	}

	for _, campus := range campuses {
		k = fmt.Sprintf("ming:%v:categories", campus)
		categories, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
		if err != nil {
			return err
		}

		for _, category := range categories {
			k = fmt.Sprintf("ming:%v:%v:periods", campus, category)
			periods, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
			if err != nil {
				return err
			}

			for _, period := range periods {
				fmt.Printf("%v:%v:%v\n", campus, category, period)
			}

		}
	}

	return nil
}
