package services

import "os"

type EnvironmentVariable struct {
	ConnStr          string
	RefreshToken     string
	DropboxAppKey    string
	DropboxAppSecret string
	Port             string
	Config           string
	Password         string
}

func ReadEnvironmentVariables() EnvironmentVariable {
	var envVars EnvironmentVariable
	envVars.ConnStr = os.Getenv("ConnStr")
	envVars.RefreshToken = os.Getenv("RefreshToken")
	envVars.DropboxAppKey = os.Getenv("DropboxAppKey")
	envVars.DropboxAppSecret = os.Getenv("DropboxAppSecret")
	envVars.Port = os.Getenv("Port")
	envVars.Config = os.Getenv("Config")
	envVars.Password = os.Getenv("Password")
	return envVars
}
