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

	if err = findStudents(config.RedisServer, config.RedisPassword); err != nil {
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

// findStudents find the students which are in 2 or more classes then output student's name, phone and classes.
func findStudents(redisServer, redisPassword string) error {
	var (
		err   error
		v     []interface{}
		items []string
	)

	conn, err := redishelper.GetRedisConn(redisServer, redisPassword)
	if err != nil {
		return err
	}
	defer conn.Close()

	k := "ming:students"
	cursor := 0
	for {
		if v, err = redis.Values(conn.Do("ZSCAN", k, cursor, "COUNT", 1000)); err != nil {
			return err
		}

		if _, err = redis.Scan(v, &cursor, &items); err != nil {
			return err
		}

		l := len(items)
		if l <= 0 || l%2 != 0 {
			continue
		}

		for i := 0; i < l; i += 2 {
			key := fmt.Sprintf("ming:%v:classes", items[i])
			count, err := redis.Int64(conn.Do("ZCARD", key))
			if err != nil {
				return err
			}

			if count < 2 {
				continue
			}

			// Output student name / phone num.
			fmt.Printf("%v\n", items[i])
			classes, err := redis.Strings(conn.Do("ZRANGE", key, 0, -1))
			if err != nil {
				return err
			}

			// Output student's  classes.
			for _, class := range classes {
				fmt.Printf("%v\n", class)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}
