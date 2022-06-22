package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

/*
func testColors(db *sql.DB) []catsColors {
	rows, err := db.Query("select * from cat_colors_info")
	if err != nil {
		log.Panic(err)
	}
	var colors []catsColors
	for rows.Next() {
		p := catsColors{}
		err := rows.Scan(&p.color, &p.count)
		if err != nil {
			log.Println(err)
			continue
		}
		colors = append(colors, p)
	}
	return colors
}

func testStatistics(db *sql.DB) []catsStats {
	rows, err := db.Query("select * from cats_stat")
	if err != nil {
		log.Panic(err)
	}
	var stats []catsStats
	for rows.Next() {
		p := catsStats{}
		err := rows.Scan(
			&p.tailLengthMean,
			&p.tailLengthMedian,
			&p.tailLengthMode,
			&p.whiskersLengthMean,
			&p.whiskersLengthMedian,
			&p.whiskersLengthMode)
		if err != nil {
			log.Println(err)
			continue
		}
		stats = append(stats, p)
	}
	return stats
}*/

func testServerPing(t *testing.T) {
	res, err := http.Get(fmt.Sprintf("localhost:8080/ping"))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("status not OK")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "Cats Service. Version 0.1\n" && string(body) != "Cats Service. Version 0.1" {
		t.Fatal()
	}
}

func testServer(t *testing.T) {
	res, err := http.Get(fmt.Sprintf("localhost:8080/cats"))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("status not OK")
	}
	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
}
