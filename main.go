package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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

func connectToDB() {
	host := getEnvValue("HOST")
	port := getEnvValue("PORT")
	dbname := getEnvValue("DBNAME")
	user := getEnvValue("USERNAME")
	password := getEnvValue("PASSWORD")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Connected to Postgres")
}

func main() {
	connectToDB()
}
