package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		app                    App
		err                    error
		buf                    []byte
		currentDir, configFile string
		s                      *ming800.Session
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

	if err = json.Unmarshal(buf, &app.Config); err != nil {
		err = fmt.Errorf("parse config err: %v", err)
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
