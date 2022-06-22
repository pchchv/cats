package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func testColors(t *testing.T) {
	rows, err := database.Query("select * from cat_colors_info")
	if err != nil {
		t.Fatal(err)
	}
	var colors []catsColors
	for rows.Next() {
		p := catsColors{}
		err := rows.Scan(&p.color, &p.count)
		if err != nil {
			t.Error(err)
			continue
		}
		colors = append(colors, p)
	}
	if len(colors) == 0 {
		t.Fatal()
	}
}

func testStatistics(t *testing.T) {
	rows, err := database.Query("select * from cats_stat")
	if err != nil {
		t.Fatal(err)
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
			t.Error(err)
			continue
		}
		stats = append(stats, p)
	}
	if len(stats) == 0 {
		t.Fatal()
	}
}

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
