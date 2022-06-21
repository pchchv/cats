package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"sort"
)

type cat struct {
	Name           string `json:"name"`
	Color          string `json:"color"`
	TailLength     int    `json:"tail_length"`
	WhiskersLength int    `json:"whiskers_length"`
}

type catsColors struct {
	color string
	count int
}

type catsStats struct {
	tailLengthMean       float64
	tailLengthMedian     float64
	tailLengthMode       []uint8
	whiskersLengthMean   float64
	whiskersLengthMedian float64
	whiskersLengthMode   []uint8
}

var database *sql.DB

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
		log.Printf("PostgreSQL not connected!\nError: %v", err)
	}
	return db
}

func pingDB(db *sql.DB) {
	/*
		Checks the connection to the database. If there is no connection, it initializes a new attempt.
		After 15 attempts, generates a Panic.
	*/
	for i := 0; i <= 16; i++ {
		err := db.Ping()
		if err != nil {
			if i == 16 {
				log.Panic(err)
			}
			log.Println("Try reconnect..")
			db = connectToDB()
		} else {
			break
		}
	}
	log.Println("Connected to Postgres")
}

func closeConnect(db *sql.DB) {
	err := db.Close()
	if err != nil {
	}
	log.Println("Connections are closed")
}

func getData(db *sql.DB) {
	/*
		Gets data from the PostgreSQL and runs the functions that process them
	*/
	rows, err := db.Query("select * from cats")
	if err != nil {
		log.Panic(err)
	}
	var catsColorsCounter []catsColors
	var tailsLengths []int
	var whiskersLengths []int
	for rows.Next() {
		p := cat{}
		err := rows.Scan(&p.Name, &p.Color, &p.TailLength, &p.WhiskersLength)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(catsColorsCounter) == 0 {
			catsColorsCounter = append(catsColorsCounter, catsColors{p.Color, 1})
		} else {
			for i, val := range catsColorsCounter {
				if val.color == p.Color {
					catsColorsCounter[i].count = val.count + 1
					break
				}
				if i == len(catsColorsCounter)-1 {
					catsColorsCounter = append(catsColorsCounter, catsColors{p.Color, 1})
				}
			}
		}
		tailsLengths = append(tailsLengths, p.TailLength)
		whiskersLengths = append(whiskersLengths, p.WhiskersLength)
	}
	log.Println(catsColorsCounter)
	// Write the result to the db
	colors(catsColorsCounter, db)
	stats(tailsLengths, whiskersLengths, db)
}

func getCats(db *sql.DB, offset string, limit string) []cat {
	if offset != "" {
		offset = " OFFSET " + offset
	}
	if limit != "" {
		limit = " LIMIT " + limit
	}
	q := fmt.Sprintf("select * from cats%s%s", offset, limit)
	rows, err := db.Query(q)
	if err != nil {
		log.Panic(err)
	}
	var cats []cat
	for rows.Next() {
		p := cat{}
		err := rows.Scan(&p.Name, &p.Color, &p.TailLength, &p.WhiskersLength)
		if err != nil {
			log.Println(err)
			continue
		}
		cats = append(cats, p)
	}
	return cats
}

func colors(catsColorsCounter []catsColors, db *sql.DB) {
	for _, val := range catsColorsCounter {
		var color = val.color
		var count = fmt.Sprint(val.count)
		_, err := db.Exec("Insert into cat_colors_info (color, count) values ($1, $2)",
			color,
			count)
		if err != nil {
			log.Panic(err)
		}
	}
}

func stats(tailsLengths []int, whiskersLengths []int, db *sql.DB) {
	sort.Ints(tailsLengths)
	sort.Ints(whiskersLengths)
	tailLengthMean := means(tailsLengths)
	tailLengthMedian := medians(tailsLengths)
	tailLengthMode := modes(tailsLengths)
	whiskersLengthMean := means(whiskersLengths)
	whiskersLengthMedian := medians(whiskersLengths)
	whiskersLengthMode := modes(whiskersLengths)
	_, err := db.Exec("Insert into cats_stat ("+
		"tail_length_mean,"+
		"tail_length_median,"+
		"tail_length_mode,"+
		"whiskers_length_mean,"+
		"whiskers_length_median,"+
		"whiskers_length_mode)"+
		"values ($1, $2, $3, $4, $5, $6)",
		tailLengthMean,
		tailLengthMedian,
		pq.Array(tailLengthMode), // The sql package does not yet support golang arrays. The pq.Array() function solves this problem
		whiskersLengthMean,
		whiskersLengthMedian,
		pq.Array(whiskersLengthMode))
	if err != nil {
		log.Panic(err)
	}
}

func means(lengths []int) float64 {
	/*
		The mean of a set of numbers, sometimes simply called the average,
		is the sum of the data divided by the total number of data.
	*/
	var mean float64
	for _, val := range lengths {
		mean += float64(val)
	}
	return mean / 2
}

func medians(lengths []int) float64 {
	/*
		The median of a set of numbers is the average number in the sorted set or, if there is an even number of data,
		the median is the average of the two average numbers.
	*/
	var median float64
	if len(lengths)%2 == 0 {
		median = float64((lengths[(len(lengths)%2)-1] + lengths[(len(lengths)%2)]) / 2)
	} else {
		median = float64(lengths[(len(lengths) % 2)])
	}
	return median
}

func modes(lengths []int) []int {
	/*
		The mode of a set of numbers is the number which occurs most often.
	*/
	var mode []int
	modeMap := make(map[int]int)
	for _, l := range lengths {
		modeMap[l]++
	}
	max := 0
	for l, num := range modeMap {
		if num > max {
			mode = []int{l}
		} else if num == max {
			mode = append(mode, l)
		}
	}
	return mode
}

func ping(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Cats Service. Version 0.1\n")
}

func cats(w http.ResponseWriter, req *http.Request) {
	offset := req.URL.Query().Get("offset")
	limit := req.URL.Query().Get("limit")
	catsList := getCats(database, offset, limit)
	for _, cat := range catsList {
		c, err := json.MarshalIndent(cat, " ", "\t")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Panic(err)
		}
		c = bytes.Replace(c, []byte("\\u0026"), []byte("&"), -1)
		w.Header().Set("Content-Type", "application/json")
		w.Write(c)
	}
}

func main() {
	db := connectToDB()
	defer closeConnect(db)
	pingDB(db)
	database = db
	/*getData(db)
	log.Println(testColors(db))
	log.Println(testStatistics(db))*/
	log.Println("Server started")
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/cats", cats)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
