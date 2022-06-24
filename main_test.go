package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	vegeta "github.com/tsenart/vegeta/v12/lib"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func TestColors(t *testing.T) {
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
}

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
		"?offset=5&limit=2",
		"?limit=3&attribute=color",
		"?attribute=color&order=asc&offset=5&limit=2"}
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

func TestLoadPing(t *testing.T) {
	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := 5 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://localhost:8080/ping",
	})
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()
	log.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}

func TestLoadCats(t *testing.T) {
	rate := vegeta.Rate{Freq: 1000, Per: time.Second}
	duration := 5 * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://localhost:8080/cats?offset=7&limit=1",
	})
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()
	log.Printf("99th percentile: %s\n", metrics.Latencies.P99)
}

func TestPostIncorrect(t *testing.T) {
	url := "http://127.0.0.1:8080/cat"
	var jsonStr = []byte(`{"name": "Tihon", "color": "red & white", "tail_length": "15", "whiskers_length": "12"}`)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode == http.StatusOK {
		t.Errorf("status OK but data is incorrect")
	}
	res.Body.Close()
}

func TestPostCorrect(t *testing.T) {
	url := "http://127.0.0.1:8080/cat"
	names := []string{"Barsik", "John", "Murka", "Marfa"}
	colors := []string{"black", "white", "black & white", "red", "red & white", "red & black & white"}
	name := names[rand.Intn(3)]
	color := colors[rand.Intn(5)]
	tl := rand.Intn(21)
	wl := rand.Intn(19)
	v := cat{name, color, tl, wl}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(v)
	res, err := http.Post(url, "application/json", b)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(res.StatusCode)
	if res.StatusCode != http.StatusOK {
		t.Errorf("status not OK")
	}
	res.Body.Close()
}
