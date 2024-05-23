package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LOAD VARIABLES JUST ONCE
var (
	ApiName      string
	ApiPort      string
	ApiLogLevel  string
	ApiJwtSecret string
	ApiCores     int
)

func Init() []map[string]string {

	var erros []map[string]string

	err := godotenv.Load()
	if err != nil {
		erros = append(erros, map[string]string{"Parâmetro": "Geral", "Mensagem": err.Error()})
		return erros
	}
	ApiName = os.Getenv("API_NAME")
	ApiPort = os.Getenv("API_PORT")
	ApiJwtSecret = os.Getenv("API_JWT_SECRET")
	ApiCoresStr := os.Getenv("API_CORES")

	// VALIDATION OF REQUIRED PARAMETERS
	if ApiName == "" {
		erros = append(erros, map[string]string{"Parâmetro": "API_HUB_DATABASE", "Mensagem": "Não preenchido."})
	}

	// VALIDATION OF OPTIONAL PARAMETERS
	if ApiName == "" {
		ApiName = "API GO"
	}
	if ApiPort == "" {
		ApiPort = "5001"
	}
	if ApiCoresStr == "" {
		ApiCores = 4
	} else {
		ApiCores, err = strconv.Atoi(ApiCoresStr)
		if err != nil {
			ApiCores = 8
		}
	}

	if len(erros) > 0 {
		return erros
	}

	return nil
}
