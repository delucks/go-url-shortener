package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
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

func encode(newid int) string {
	dict := "abcdefghijklmnopqrstuvwxyz0123456789"
	dict_len := len(dict)
	inc := newid
	conv_ints := make([]int, 0, 1)
	for inc > 0 {
		remainder := inc % 62
		conv_ints = append(conv_ints, remainder)
		inc = inc / dict_len
	}
	fmt.Println(conv_ints)
	sort.Sort(sort.Reverse(sort.IntSlice(conv_ints)))
	fmt.Println(conv_ints)
	final := ""
	for i := 0; i < len(conv_ints); i++ {
		//final += strconv.Itoa(conv_ints[i])
		final += fmt.Sprintf("%c", dict[i])
	}
	return string(final)
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

func setupDB(db *sql.DB) {
	result, err := db.Exec("CREATE TABLE IF NOT EXISTS url(id INT PRIMARY KEY NOT NULL AUTO_INCREMENT, mapping VARCHAR(255) NOT NULL)")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(result)
}

func main() {
	args := os.Args
	if len(args) > 1 {
		if args[1] == "testAlg" {
			if len(args) < 2 {
				fmt.Println("Invalid number of arguments to -t")
			} else {
				val, err := strconv.Atoi(args[2])
				if err != nil {
					fmt.Println("Conversion failed!")
					return
				}
				fmt.Printf("Int to convert: %d\n", val)
				fmt.Printf("Converted int:  %s\n", encode(val))
			}
		} else {
			fmt.Printf("Invalid argument %s\n", args[1])
		}
	} else {
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
		setupDB(db)
		for true {
			var prev int
			row := db.QueryRow("SELECT max(id) FROM url")
			err := row.Scan(&prev)
			if err != nil {
				fmt.Println("Error getting next id")
				return
			}
			fmt.Printf("To insert: %d\n", prev+1)
			fmt.Printf("The string to use %s\n", encode(prev+1))
		}
		//http.HandleFunc("/", handleURL)
		//http.ListenAndServe("127.0.0.1:8080",nil)
	}
}
