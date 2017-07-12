package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"strconv"

	_ "github.com/GO-SQL-Driver/MySQL"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(ip:3306)/scoreboard?charset=utf8")
	if err != nil {
		panic("Database Connection Failed!")
	}
}

//	/api?op=new&name=瑞典&data='flag url'
//  /api?op=update&name=瑞典&data=10000000 //append data
//  /api?op=get
func handleIndex(w http.ResponseWriter, r *http.Request) {
	operator, name, data := r.URL.Query().Get("op"), r.URL.Query().Get("name"), r.URL.Query().Get("data")

	if operator == "new" {

		//flag := data

	} else if operator == "update" {

		score, _ := strconv.Atoi(data)

		stmt, err := db.Prepare("UPDATE countries SET `score` = `score` + ? WHERE country = ?")
		if err != nil {
			panic(err)
		}
		res, err := stmt.Exec(score, name)
		if err != nil {
			panic(err)
		}
		stmt.Close()

		affect, err := res.RowsAffected()
		if err != nil {
			panic(err)
		}
		fmt.Println(affect)
		fmt.Fprintf(w, `{ affect: `+strconv.FormatInt(affect, 10)+`}`)

	} else if operator == "get" {

		rows, err := db.Query("SELECT * FROM countries")
		if err != nil {
			panic(err)
		}

		var json string = "["
		for rows.Next() {
			var country string
			var score int64
			var flag string
			rows.Scan(&country, &score, &flag)

			fmt.Println("coutnry: ", country, ", score: ", score)

			json += "{\"country\":\"" + country + "\",\"score\":" + strconv.FormatInt(score, 10) + ", \"flag\":\"" + flag + "\"}"
		}
		json += "]"

		fmt.Fprintf(w, json)

	} else {
		fmt.Fprintf(w, "{ error: \"No Operator Set In Parameter!\"}")
	}
	//score, _ := strconv.Atoi(data)
	//fmt.Fprintf(w, "data: %s", string(operator))
}

func main() {

	http.Handle("/", http.FileServer(http.Dir("./src")))
	http.HandleFunc("/api", handleIndex)
	http.ListenAndServe(":9090", nil)
}
