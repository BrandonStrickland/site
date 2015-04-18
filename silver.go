// This is just a test run to see how Go works with http services.
// I might do something more interesting than the standard wikia.
// Since I am going to work with MySQL, I think it would be
// appropriate to generate some data and run over it with some
// algorithm and see if we can get something close to the function
// that generated it.
package main

import (
	"errors"
	"flag"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

// validPath is a solution to keeping people from arbitrarily giving paths
// to be written/read on the server, we are going to use regular expressions to
// make it fit more rigid guidelines.
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

// i dont dont know what this is for.
var addr = flag.Bool("addr", false, "find open address and print to final-port.txt")

// Page is the contents of the page the browser will display. A page
// is broken into two pieces, the title and the body of the page.
type Page struct {
	Title string
	Body  []byte // We are using byte slices since the libraries
	// will be expecting byte slices.
}

// Save takes the page and converts it into a .txt for storage in the HDD.
func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
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

// viewHandler loads the page requested from the user via URL and returns html
// for the user.
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

// editHandler loads the html for editting a page and returns the user to
// to the edit page.
func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

// saveHandler takes form data and sends the Page to the save function to write
// it to the HDD.
func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err = p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

// makeHandler passes a http.Handler back to HandleFunc. This is allows us to more
// or less cut down on error handling duplication done for things like checking
// the title. Now, we simply pass the title to the function with this.
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.ULR.Path)
		if m != nill {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

// main sends patterns to the handler and then we let the ListenAndServe take
// over and we block on it until the application crashes or is stopped by someone
// who started it.
func main() {
	flag.Parse()
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))

}
