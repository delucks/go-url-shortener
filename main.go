package main

import (
	"net/http"
	"fmt"
)

func main() {
	http.HandleFunc("/", handleMain)
	http.ListenAndServe("127.0.0.1:8080",nil)
}

func handleMain(w http.ResponseWriter, r *http.Request) {
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
}
