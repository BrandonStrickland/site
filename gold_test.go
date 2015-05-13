package main

import (
	"testing"
	"net/http"
	"log"
	"encoding/json"
	"math/rand"
	"bytes"
	//"net/http"
	"net/http/httptest"
	"database/sql"
	_ "github.com/erikstmartin/go-testdb"
)

// I don't feel too terrible for making this global in a test.
//var db,_ = sql.Open("testdb", "")
var db,_ = sql.Open("mysql", "brandon:password@tcp(127.0.0.1:3306)/test")

/*
type Point struct {
	ID int
	X  int
	Y  int
}

func TestInitTable(t *testing.T) {
	
}

func TestFillTable(t *testing.T) {

}

func TestCleanTable(t *testing.T) {

}

func TestJsonHandler(t *testing.T) {

}

func TestSetupDatabase(t *tsting.T) {

}
*/

// generatePoint takes in an id so that when we add it to the database,
// we can have the primary key be the id and since it is incremented,
// inserts should be blazingly fast.
func generatePoint(id int) Point {
	x := int(rand.Int31() / 2)
	y := x * 2
	return Point{id, x, y}
}

// createJsonPackage takes in an id and creates a json to send to the server
func createJsonPackage(id int) ([]byte, error) {
	point := generatePoint(id)
	mar, err := json.Marshal(point)
	if err != nil {
		log.Println(err)
	}
	var p Point
	err = json.Unmarshal(mar, &p)
	if err != nil {
		log.Println(err)
	}

	return json.Marshal(p)
}

func TestHandler(t *testing.T) {
	pay, err := createJsonPackage(1)
	if err != nil {
		log.Fatal("The Json could not be Marshaled up")
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:8080", bytes.NewReader(pay))
	if err != nil {
		log.Println("Could not set up the request")
		log.Fatal(err)
	}
	
	w := httptest.NewRecorder()
	Handler(w, req, db)
	if w.Code != 200 {
		t.Fail()
	}
}
