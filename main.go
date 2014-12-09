package main

/*
 * go-url-shortener
 *
 * Author: delucks
 * Requirements: github.com/go-sql-driver/mysql
 *
 */

import (
	"database/sql"
	"fmt"
	"log"
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
 * SQL queue:
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
	//fmt.Println(conv_ints)
	sort.Sort(sort.Reverse(sort.IntSlice(conv_ints)))
	//fmt.Println(conv_ints)
	final := ""
	for i := 0; i < len(conv_ints); i++ {
		//final += strconv.Itoa(conv_ints[i])
		final += fmt.Sprintf("%c", dict[i])
	}
	return string(final)
}

func decode(ext string) int {
	return 0
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
		log.Fatal(err.Error())
		return
	}
	fmt.Println(result)
}

func connectDB() *sql.DB {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/go_url_shortener")
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Opening failed")
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Ping failed")
	}
	return db
}

func main() {
	args := os.Args
	if len(args) > 1 {
		if args[1] == "-t" {
			if len(args) < 2 {
				log.Fatal("Invalid number of arguments to -t")
			} else {
				val, err := strconv.Atoi(args[2])
				if err != nil {
					log.Fatal("\x1b[31mConversion failed!\x1b[0m")
				}
				fmt.Printf("Int to convert: %d\n", val)
				fmt.Printf("Converted int:  %s\n", encode(val))
			}
		} else {
			log.Fatalf("Invalid argument %s\n", args[1])
		}
	} else {
		var db *sql.DB
		db = connectDB()
		setupDB(db)
		for true {
			var prev int
			row := db.QueryRow("SELECT max(id) FROM url")
			err := row.Scan(&prev)
			if err != nil {
				log.Fatal("Error getting next id")
				return
			}
			fmt.Printf("\x1b[32m[::]\x1b[0m To insert:         %d\n", prev+1)
			fmt.Printf("\x1b[32m[::]\x1b[0m The string to use: %s\n", encode(prev+1))
			result, err := db.Exec("INSERT INTO url(mapping) values(?)", encode(prev+1))
			if err != nil {
				log.Fatal("Insert failed!")
			} else {
				fmt.Println(result)
			}
		}
		//http.HandleFunc("/", handleURL)
		//http.ListenAndServe("127.0.0.1:8080",nil)
	}
}
