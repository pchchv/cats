package main

import (
	"database/sql"

	"fmt"
	"github.com/joho/godotenv"
	"github.com/lib/pq"
	"log"
	"os"
	"sort"
)

type cat struct {
	name           string
	color          string
	tailLength     int
	whiskersLength int
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

func getData(db *sql.DB) {
	rows, err := db.Query("select * from cats")
	if err != nil {
		log.Panic(err)
	}
	var catsColorsCounter []catsColors
	var tailsLengths []int
	var whiskersLengths []int
	for rows.Next() {
		p := cat{}
		err := rows.Scan(&p.name, &p.color, &p.tailLength, &p.whiskersLength)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(catsColorsCounter) == 0 {
			catsColorsCounter = append(catsColorsCounter, catsColors{p.color, 1})
		} else {
			for i, val := range catsColorsCounter {
				if val.color == p.color {
					catsColorsCounter[i].count = val.count + 1
					break
				}
				if i == len(catsColorsCounter)-1 {
					catsColorsCounter = append(catsColorsCounter, catsColors{p.color, 1})
				}
			}
		}
		tailsLengths = append(tailsLengths, p.tailLength)
		whiskersLengths = append(whiskersLengths, p.whiskersLength)
	}
	colors(catsColorsCounter, db)
	stats(tailsLengths, whiskersLengths, db)
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
		pq.Array(tailLengthMode),
		whiskersLengthMean,
		whiskersLengthMedian,
		pq.Array(whiskersLengthMode))
	if err != nil {
		log.Panic(err)
	}
}

func means(lengths []int) float64 {
	var mean float64
	for _, val := range lengths {
		mean += float64(val)
	}
	return mean / 2
}

func medians(lengths []int) float64 {
	var median float64
	if len(lengths)%2 == 0 {
		median = float64((lengths[(len(lengths)%2)-1] + lengths[(len(lengths)%2)]) / 2)
	} else {
		median = float64(lengths[(len(lengths) % 2)])
	}
	return median
}

func modes(lengths []int) []int {
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

func main() {
	db := connectToDB()
	defer closeConnect(db)
	getData(db)
	log.Println(testColors(db))
	log.Println(testStatistics(db))
}
