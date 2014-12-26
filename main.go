package main

/*
 * go-url-shortener
 *
 * Author: delucks
 * Requirements: github.com/go-sql-driver/mysql
 *
 */

/* TODO:
 * Check if url already in database, if so return the already generated ID
 */

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	//"sort"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

/*
 * Base conversion functions
 */

func getord(input byte) int {
	LOWER_OFFSET := 87
	DIGIT_OFFSET := 48
	UPPER_OFFSET := 29
	var result int
	if input <= 57 && input >= 48 {
		result = int(input) - DIGIT_OFFSET
	} else if input >= 97 && input <= 122 {
		result = int(input) - LOWER_OFFSET
	} else if input >= 65 && input <= 90 {
		result = int(input) - UPPER_OFFSET
	} else {
		fmt.Printf("Dafux is this\n")
		result = 0
	}
	//fmt.Printf("%c as base10 is %d\n", input, result)
	return result
}

func getchr(input int) byte {
	LOWER_OFFSET := 87
	DIGIT_OFFSET := 48
	UPPER_OFFSET := 29
	var result byte
	if input < 10 {
		result = byte(input + DIGIT_OFFSET)
	} else if input < 36 {
		result = byte(input + LOWER_OFFSET)
	} else if input < 63 {
		result = byte(input + UPPER_OFFSET)
	} else {
		fmt.Printf("Dafux is this\n")
		result = byte(0)
	}
	//fmt.Printf("%d as base65 is %c\n", input, result)
	return result
}

func encode(newid int) string {
	if newid == 0 {
		return "0"
	}
	dict_len := 62
	inc := newid
	conv_ints := make([]byte, 0, 1)
	for inc > 0 {
		remainder := inc % 62
		conv_ints = append(conv_ints, getchr(remainder))
		inc = inc / dict_len
	}
	final := ""
	for i := len(conv_ints) - 1; i >= 0; i-- {
		final += fmt.Sprintf("%c", conv_ints[i])
	}
	return final
}

func decode(ext string) int {
	runes := []rune(ext)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	bytes := []byte(string(runes))
	//fmt.Println(bytes)
	result := 0
	for i := 0; i < len(bytes); i++ {
		result += getord(bytes[i]) * int(math.Pow(float64(62), float64(i)))
	}
	return result
}

/*
 * Database
 */

func setupDB(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS url(id INT PRIMARY KEY NOT NULL AUTO_INCREMENT, mapping VARCHAR(255) NOT NULL)")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
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

/*
 * HTTP
 */

func geturl(id int) string {
	var final string
	var db *sql.DB
	db = connectDB()
	selstmt := fmt.Sprintf("SELECT mapping from url WHERE id='%d'", id)
	row := db.QueryRow(selstmt)
	err := row.Scan(&final)
	if err != nil {
		fmt.Printf("Error getting mapping for ID %d\n", id)
		fmt.Println(err.Error())
		final = "no such URL"
	}
	db.Close()
	return final
}

func addurl(url string) string {
	var prev int
	var db *sql.DB
	db = connectDB()
	row := db.QueryRow("SELECT max(id) FROM url")
	err := row.Scan(&prev)
	if err != nil {
		log.Fatal("Error getting next id")
		log.Fatal(err.Error())
		return ""
	}
	//fmt.Printf("\x1b[32m[::]\x1b[0m To insert:         %d\n", prev+1)
	//fmt.Printf("\x1b[32m[::]\x1b[0m The string to use: %s\n", encode(prev+1))
	mapping := encode(prev + 1)
	result, err := db.Exec("INSERT INTO url(mapping) values(?)", url)
	if err != nil {
		log.Fatalf("Insert failed for %s, index %d, encoding %s", url, prev+1, mapping)
	} else {
		fmt.Printf("[%d] %s to %s (%s)\n", prev+1, url, mapping, result)
	}
	db.Close()
	return mapping
}

func handleURL(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		if r.Method == "POST" {
			r.ParseForm() // needed so we have the post'd params
			url := r.Form.Get("url")
			if url != "" { // go sets empty strings to ""
				mapping := addurl(url)
				w.Write([]byte("http://thatson.me/" + mapping))
			} else {
				w.Write([]byte("POST the site with a valid url"))
			}
		} else {
			w.Write([]byte("Welcome to my URL shortener. Please POST this path with a valid URL."))
		}
	} else {
		req := strings.Split(r.URL.Path, "/")[1] // grab the first URI after / to use as the url
		if req != "favicon.ico" {
			refresh := fmt.Sprintf("<html><head><meta http-equiv='Refresh' content='0; url=%s' /></head></html>", geturl(decode(req)))
			w.Write([]byte(refresh))
			fmt.Printf("Returned %s on query string %s\n", geturl(decode(req)), req)
		}
	}
}

/*
 * Main
 */

func main() {
	args := os.Args
	if len(args) > 1 {
		if args[1] == "-e" {
			if len(args) < 2 {
				log.Fatal("Invalid number of arguments to -e. Pass in an integer to be converted")
			} else {
				val, err := strconv.Atoi(args[2])
				if err != nil {
					log.Fatal("\x1b[31mConversion failed!\x1b[0m")
				}
				fmt.Printf("Int to convert: %d\n", val)
				fmt.Printf("Converted int:  %s\n", encode(val))
			}
		} else if args[1] == "-d" {
			if len(args) < 2 {
				log.Fatal("Invalid number of arguments to -d. Pass in an string to be converted")
			} else {
				fmt.Printf("String to decode: %s\n", args[2])
				fmt.Printf("Converted int:    %d\n", decode(args[2]))
			}
		} else if args[1] == "-t" {
			if len(args) < 2 {
				log.Fatal("Invalid number of arguments to -e. Pass in an integer to test")
			} else {
				val, err := strconv.Atoi(args[2])
				if err != nil {
					log.Fatal("\x1b[31mConversion failed!\x1b[0m")
				}
				fmt.Printf("int to encode: %d\n", val)
				result := encode(val)
				fmt.Printf("base62 representation: %s\n", result)
				fmt.Printf("and back: %d\n", decode(result))
			}
		} else {
			log.Fatalf("Invalid argument %s\n", args[1])
		}
	} else {
		var DBCONN *sql.DB
		DBCONN = connectDB()
		setupDB(DBCONN)
		DBCONN.Close()
		http.HandleFunc("/", handleURL)
		http.ListenAndServe("127.0.0.1:8080", nil)
	}
}
