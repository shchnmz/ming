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

	if err = getStudentCountOfEachTeacher(config.RedisServer, config.RedisPassword); err != nil {
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

func getStudentCountOfEachTeacher(redisServer, redisPassword string) error {
	db := ming.DB{redisServer, redisPassword}

	teachers, err := db.GetTeachers()
	if err != nil {
		return err
	}

	for _, teacher := range teachers {
		students, err := db.GetStudentsOfTeacher(teacher)
		if err != nil {
			return err
		}

		fmt.Printf("%v: 共%v人\n", teacher, len(students))
	}

	return nil
}
