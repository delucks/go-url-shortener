package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strings"
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
 *
 * SQL used:
 * create database if not exists go_url_shortener;
 */

func shorten(original string) string {
	dict := "abcdefghijklmnopqrstuvwxyz0123456789"
	//dict_len := len(dict)
	fmt.Printf("%q", dict[0])
	final := ""
	for i := 0; i < len(original); i++ {
		final += fmt.Sprintf("%c", original[i])
	}
	return final
}

func handleURL(w http.ResponseWriter, r *http.Request) {
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
			w.Write([]byte("Welcome to my URL shortener. Please POST this path with a valid URL."))
		}
	} else {
		req := strings.Split(r.URL.Path, "/")[1] // grab the first URI after / to use as the url
		fmt.Println(req)
	}
}

func main() {
	fmt.Println(shorten("bar"))
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/go_url_shortener")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Opening failed")
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Ping failed")
		return
	}
	//http.HandleFunc("/", handleURL)
	//http.ListenAndServe("127.0.0.1:8080",nil)
}
