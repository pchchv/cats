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
