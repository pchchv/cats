package main

import (
	"database/sql"
	"log"
)

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
}
