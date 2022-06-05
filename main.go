package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	// Load values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Panic("No .env file found")
	}
}

func getEnvValue(v string) string {
	value, exist := os.LookupEnv(v)
	if !exist {
		log.Panicf("Value %v does not exist", v)
	}
	return value
}

func main() {
	host := getEnvValue("HOST")
	port := getEnvValue("PORT")
	dbname := getEnvValue("DBNAME")
	username := getEnvValue("USERNAME")
	password := getEnvValue("PASSWORD")
	fmt.Println(host)
	fmt.Println(port)
	fmt.Println(dbname)
	fmt.Println(username)
	fmt.Println(password)
}
