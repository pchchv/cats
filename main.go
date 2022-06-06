package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

type cat struct {
	name           string
	color          string
	tailLength     int
	whiskersLength int
}

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

func connectToDB() *sql.DB {
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
		log.Panic(err)
	}
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	log.Println("Connected to Postgres")
	return db
}

func closeConnect(db *sql.DB) {
	err := db.Close()
	if err != nil {
	}
	log.Println("Connections are closed")
}

func main() {
	db := connectToDB()
	defer closeConnect(db)
	rows, err := db.Query("select * from cats")
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()
	var cats []cat
	for rows.Next() {
		p := cat{}
		err := rows.Scan(&p.name, &p.color, &p.tailLength, &p.whiskersLength)
		if err != nil {
			fmt.Println(err)
			continue
		}
		cats = append(cats, p)
	}
	fmt.Println(cats)
}
