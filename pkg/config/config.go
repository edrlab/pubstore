package config

import (
	"os"
	"strconv"
)

const OauthSeed = "Edrlab-Rocks"
const PrintRights = 20
const CopyRights = 2000
const NumberOfPublicationsPerPage = 50

var BASE_URL = "http://localhost:8080"

var PORT = 8080

var LCP_SERVER_URL = "https://front-prod.edrlab.org/lcpserver"

var LCP_SERVER_USERNAME = "adm_username"
var LCP_SERVER_PASSWORD = "adm_password"

var DSN = "pub.db"

func Init() {

	var portEnv = os.Getenv("PORT")

	if portEnv != "" {
		portInt, err := strconv.Atoi(portEnv)
		if err != nil {
			if portInt >= 80 && portInt <= 9999 {
				PORT = portInt
			}
		}
	}

	var baseUrlEnv = os.Getenv("BASE_URL")
	if baseUrlEnv != "" {
		BASE_URL = baseUrlEnv
	}

	var lcpServerUrlEnv = os.Getenv("LCP_SERVER_URL")
	if lcpServerUrlEnv != "" {
		LCP_SERVER_URL = lcpServerUrlEnv
	}

	var lcpServerPassword = os.Getenv("LCP_SERVER_PASSWORD")
	if lcpServerPassword != "" {
		LCP_SERVER_PASSWORD = lcpServerPassword
	}

	var lcpServerUsername = os.Getenv("LCP_SERVER_USERNAME")
	if lcpServerUsername != "" {
		LCP_SERVER_USERNAME = lcpServerUsername
	}

	var dsnEnv = os.Getenv("DSN")
	if dsnEnv != "" {
		DSN = dsnEnv
	}
}
