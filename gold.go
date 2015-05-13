package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var validPath = regexp.MustCompile("/")

// Point is currently a 2D Cartesian point. This can be expanded to xD point via an array or list.
type Point struct {
	ID int
	X  int
	Y  int
}

// initTable creates a table for us to use.
func InitTable(db *sql.DB) error {
	_, err := db.Exec(
		"create table `temp` ( `ID` bigint(20) NOT NULL, `X` bigint(20) NOT NULL, `Y` bigint(20) NOT NULL, PRIMARY KEY (`ID`) ) ENGINE=MEMORY;")
	if err != nil {
		return err
	}
	return nil
}

// fillTable takes the db pointer and a point and adds to to the database
// If it doesn't add it to the table, for this application, it doesn't matter.
// So we aren't going to deal with the error.
func FillTable(db *sql.DB, point Point) error {
	_, err := db.Exec("insert into temp(id,x,y) values(?,?,?)", point.ID, point.X, point.Y)
	if err != nil {
		return err
	}
	return nil
}

// cleanTables ensures we aren't playing with any old values. restarting the server
// for each run would be silly.
func CleanTable(db *sql.DB) error {
	_, err := db.Exec("drop table if exists `temp`")
	if err != nil {
		return err
	}
	return nil
}

// jsonHandler takes a byte slice from the readAll call and sends it to function that inserts it into the table.
func JsonHandler(db *sql.DB, body []byte) error {
	var p Point
	err := json.Unmarshal(body, &p)
	if err != nil {
		return err
	}

	err = FillTable(db, p)
	if err != nil {
		return err
	}
	return nil
}

func SetupDatabase() (*sql.DB, error) {
	database, err := sql.Open("mysql", "brandon:password@tcp(127.0.0.1:3306)/test")
	if err != nil {
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		return nil, err
	}

	if err = CleanTable(database); err != nil {
		return nil, err
	}
	if err = InitTable(database); err != nil {
		return nil, err
	}
	return database, err
}

func Handler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	r.ParseForm()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "something happened", http.StatusInternalServerError)
	}
	if err = JsonHandler(db, body); err != nil {
		http.Error(w, "something happened", http.StatusInternalServerError)
	}
}

// notes:
// 2nd argument to sql.Open() is standard to open the sql port. I am sure there is a more
// secure way to open a connection without showing it in the code. Probably passing them as
// arguments to the program when it runs. can use os.Args() for this I think.
//
// The program idea follows:
// 1. listenAndServe
// 2. receive json and unmarshal it
// 3. send the point to the mysql database
// 4. send ok if successful or error if not
//
// next to implement:
// 1. add auth tokens
// 2. work more with internet response codes
func main() {
	db, err := SetupDatabase()
	if err != nil {
		log.Fatal("The database could not be correctly configured")
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Handler(w, r, db)
	})
	http.ListenAndServe(":8080", nil)
}
