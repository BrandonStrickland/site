// This is just a test run to see how Go works with http services.
// I might do something more interesting than the standard wikia.
// Since I am going to work with MySQL, I think it would be
// appropriate to generate some data and run over it with some
// algorithm and see if we can get something close to the function
// that generated it.
package main

import (
//	"html/template"
//	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
//	"regexp"
	"bytes"
	"encoding/json"
)

type Point struct {
	ID int
	X int
	Y int
}

func generatePoint(id int) Point {
	x := int(rand.Int31() / 2)
	y := x * 2
	return Point{ id, x, y }
}

func createJsonPackage(id int) ([]byte, error){
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

func postToServer(url string) {
client := &http.Client{}
	i := 0
	for i < 50000 {
		pay, err := createJsonPackage(i)
		if err != nil {
			log.Println(err)
		}

		req, err := http.NewRequest("POST", url, bytes.NewReader(pay))
		if err != nil {
			log.Println(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Println(err) // I dont know how to handle this right now
		}
		defer resp.Body.Close()

		if err != nil {
			log.Println(err) // i dont know how to handle this right now
		}
		defer resp.Body.Close()
		
		if resp.StatusCode >= 400 {
			log.Println(err) // i have no idea what to do here right now
		}
		i = i + 1
	}
}

func main() {
	url := "http://127.0.0.1:8080"
	postToServer(url)
}







/*
// validPath is a solution to keeping people from arbitrarily giving paths
// to be written/read on the server, we are going to use regular expressions to
// make it fit more rigid guidelines.
var validPath = regexp.MustCompile("/")

// i dont dont know what this is for.
//var addr = flag.Bool("addr", false, "find open address and print to final-port.txt")

// Page is the contents of the page the browser will display. A page
// is broken into two pieces, the title and the body of the page.
type Page struct {
	Title string
	Body  []byte // We are using byte slices since the libraries
	// will be expecting byte slices.
}

// loadPage reads from the HDD and returns the resulting page back to the
// user.
func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

// renderTemplate takes a ResponseWriter, the title of a html file, and page
// and populates the html file with the contents of the html file with the body
// and title of the page. If anything were to go wrong,
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// makeHandler passes a http.Handler back to HandleFunc. This is allows us to more
// or less cut down on error handling duplication done for things like checking
// the title. Now, we simply pass the title to the function with this.
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, "home")
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: "home", Body: []byte("the page didn't load correctly")}
	}
	renderTemplate(w, "home", p)
}
*/
// main sends patterns to the handler and then we let the ListenAndServe take
// over and we block on it until the application crashes or is stopped by someone
// who started it.

