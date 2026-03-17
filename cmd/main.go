package main

import (
	"job_vacancies/config"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("config/.env")
	if err != nil {
		log.Fatalf("err loading env: %v", err)
	}
	httpcfg := config.LoadHttpConfig() //only has port
	apikeys := config.LoadApiKeys()

}
