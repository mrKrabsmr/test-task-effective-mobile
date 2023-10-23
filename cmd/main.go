package main

import (
	"github.com/joho/godotenv"
	server "github.com/mrKrabsmr/commerce-edu-api/internal"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	s := server.NewAPIServer()
	s.StartMigrations()
	s.Run()
}
