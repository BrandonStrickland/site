package main

import (
	"database/sql" 
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"log"
	"regexp"
	"encoding/json"
	"io/ioutil"
)

var validPath = regexp.MustCompile("/")
var db *sql.DB

// Point is currently a 2D Cartesian point. This can be expanded to xD point via an array or list.
type Point struct {
	ID int
	X int
 	Y int
}

// initTable creates a table for us to use.
func initTable(db *sql.DB) {
	_, err := db.Exec(
		"create table `temp` ( `ID` bigint(20) NOT NULL, `X` bigint(20) NOT NULL, `Y` bigint(20) NOT NULL, PRIMARY KEY (`ID`) ) ENGINE=MEMORY;")
	if err != nil {
		log.Println(err)
	}
}

// fillTable takes the db pointer and a point and adds to to the database
// If it doesn't add it to the table, for this application, it doesn't matter.
// So we aren't going to deal with the error.
func fillTable(db *sql.DB, point Point) {
	db.Exec("insert into temp(id,x,y) values(?,?,?)", point.ID, point.X, point.Y)
}

// cleanTables ensures we aren't playing with any old values. restarting the server
// for each run would be silly.
func cleanTable(db *sql.DB) {
	_, err := db.Exec("drop table if exists `temp`")
	if err != nil {
		log.Println(err)
	}
}

// jsonHandler takes a byte slice from the readAll call and sends it to function that inserts it into the table.
func jsonHandler(body []byte) error {
	var p Point
	err := json.Unmarshal(body, &p)
	if err != nil {
		return err
	}
	go fillTable(db,p)
	return nil
}

// requestHandler will do all error checking. The following is a list of resposibilities:
// 1. check formatting
// 2. send proper error back if need be, otherwise, 200 -- find out how to send this back correctly
// 3. pass to fillTable(db)
func requestHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err) // more error handling I dont understand right now
	}
	if err = jsonHandler(body); err != nil {
		log.Println(err)
	}
}

// main does a few things that could be abstracted out but no need to worry about that right now.
// 
// notes:
// 2nd argument to sql.Open() is standard to open the sql port. I am sure there is a more
// secure way to open a connection without showing it in the code. Probably passing them as
// arguments to the program when it runs. can use os.Args() for this I think. 
// 
// ~~~~~~things to fix~~~~~~~
// THE GLOBAL DB POINTER
// Ideally, there should be one entry and one exit, reconstruct functions to correctly handle errors.
//  
// The program idea follows:
// 1. listenAndServe 
// 2. receive json and unmarshal it
// 3. send the point to the mysql database
// 
// next to implement:
// 1. add auth tokens
// 2. work more with internet response codes
// 3. implement line of best fit (since linear, Pearson or whoever coeff
func main() {
	database, err := sql.Open("mysql", "brandon:password@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Println(err)
	}
	defer database.Close()

	err = database.Ping()
	if err != nil {
		log.Println(err)
	}
	
	db = database
	cleanTable(db)
	initTable(db)
	
	http.HandleFunc("/", requestHandler)
	http.ListenAndServe(":8080", nil)
}
