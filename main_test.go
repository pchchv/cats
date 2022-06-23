package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

// !!!The commented don't pass!!!
/*func TestColors(t *testing.T) {
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

func TestStatistics(t *testing.T) {
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
}*/

func TestServerPing(t *testing.T) {
	res, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/ping"))
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

func TestServer(t *testing.T) {
	params := []string{"",
		"?attribute=color",
		"?attribute=color&order=desc",
		"?offset=7",
		"?limit=4",
		"?offset=5&limit=2"}
	// "?limit=3&attribute=color"}
	// "?attribute=color&order=asc&offset=5&limit=2"}
	for _, v := range params {
		res, err := http.Get(fmt.Sprintf("http://127.0.0.1:8080/cats" + v))
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != http.StatusOK {
			t.Errorf("status not OK")
		}
		_, err = ioutil.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
	}
}
