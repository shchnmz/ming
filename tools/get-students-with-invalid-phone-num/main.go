package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"github.com/northbright/ming800"
	"github.com/northbright/pathhelper"
	"github.com/shchnmz/ming"
)

// Config contains the app settings.
type Config struct {
	ServerURL string `json:"server_url"`
	Company   string `json:"company"`
	User      string `json:"user"`
	Password  string `json:"password"`
}

// App represents the application.
type App struct {
	Config
}

var (
	app App
)

// ClassHandler implements ming800.WalkProcessor interface.
func (a *App) ClassHandler(class ming800.Class) {}

// StudentHandler implements ming800.WalkProcessor interface
func (a *App) StudentHandler(class ming800.Class, student ming800.Student) {
	// Check if phone number: 11-digit or 8-digit.
	if !ming.ValidPhoneNum(student.PhoneNum) {
		fmt.Printf("%s,%s,%s,%s\n", class.Category, class.Name, student.Name, student.PhoneNum)
	}
}

func main() {
	var (
		err error
		s   *ming800.Session
	)

	defer func() {
		if err != nil {
			log.Printf("%v", err)
		}
	}()

	if err = loadConfig("config.json", &app.Config); err != nil {
		return
	}

	// New a session
	if s, err = ming800.NewSession(app.ServerURL, app.Company, app.User, app.Password); err != nil {
		err = fmt.Errorf("NewSession() error: %v", err)
		return
	}

	// Login
	if err = s.Login(); err != nil {
		err = fmt.Errorf("Login() error: %v", err)
		return
	}

	// Walk
	// Class and student handler will be called while walking ming800.
	if err = s.Walk(&app); err != nil {
		err = fmt.Errorf("Walk() error: %v", err)
		return
	}

	// Logout
	if err = s.Logout(); err != nil {
		err = fmt.Errorf("Logout() error: %v", err)
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

	return json.Unmarshal(buf, &config)
}
