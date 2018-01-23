package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}

	r := gin.Default()
	r.POST("/temperature", func(c *gin.Context) {
		innerTemperature := c.PostForm("inner")
		outerTemperature := c.PostForm("outer")

		statement, err := db.Prepare("INSERT INTO temperature (session_id, inner, outer) VALUES (?, ?, ?)")
		if err != nil {
			panic(err)
		}
		_, err = statement.Exec(1, innerTemperature, outerTemperature)
		if err != nil {
			panic(err)
		}

		c.Status(200)
	})

	r.GET("/list", func(c *gin.Context) {

		rows, _ := db.Query("SELECT inner FROM temperature")

		var id int
		var ret string

		for rows.Next() {
			_ = rows.Scan(&id)
			ret += fmt.Sprintln(id)
		}

		rows.Close()

		c.String(200, ret)

	})

	r.Run(":8177")
	db.Close()

}
