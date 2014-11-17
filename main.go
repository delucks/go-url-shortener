package main

import (
	"net/http"
	"fmt"
)

/*
 * Options for generating unique url paths: 
 * Increment something globally and base64 encode it for the mapping
 * Take a hash of the url and keep that for the mapping
 *
 * Either way: when a request comes in for something other than /
 * Look it up in the database. If the key exists, return a HTTP redirect to its url
 * If not, return an error page with a link back to /
 * / should be some kind of a welcome page with a form to submit the URL to the POST
 */

func main() {
	http.HandleFunc("/", handleMain)
	http.ListenAndServe("127.0.0.1:8080",nil)
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		if r.Method == "POST" {
			r.ParseForm() // needed so we have the post'd params
			url := r.Form.Get("url")
			if url != "" { // go sets empty strings to ""
				w.Write([]byte(url))
				fmt.Println(url)
			} else {
				w.Write([]byte("POST the site with a valid url"))
			}
		} else {
			w.Write([]byte("Welcome to my URL shortener. Please POST this path with a valid url."))
		}
	} else {
		fmt.Println(r.URL.Path)
	}
}
