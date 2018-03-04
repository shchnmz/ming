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

// Config represents app settings.
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

	if err = findPhones(config.RedisServer, config.RedisPassword); err != nil {
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

// findPhones find the phones which have 2 or more students, then output phone number, student's names.
func findPhones(redisServer, redisPassword string) error {
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

	k := "ming:phone_nums"
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
			key := fmt.Sprintf("ming:%v:students", items[i])
			count, err := redis.Int64(conn.Do("ZCARD", key))
			if err != nil {
				return err
			}

			if count < 2 {
				continue
			}

			// Output phone num.
			fmt.Printf("%v\n", items[i])

			names, err := redis.Strings(conn.Do("ZRANGE", key, 0, -1))
			if err != nil {
				return err
			}

			// Output student's names.
			for _, name := range names {
				fmt.Printf("%v\n", name)
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}
