package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./db/nc.db")
	if err != nil {
		panic(err)
	}

	err = InitDb(db)
	if err != nil {
		panic(err)
	}

	go nc_server(db)

	http.HandleFunc("/api/search", func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()

		page := query.Get("page")
		if page == "" {
			page = "1"
		}
		pageNum, err := strconv.Atoi(page)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s", err)
		}

		limit := query.Get("limit")
		if limit == "" {
			limit = "20"
		}
		limitNum, err := strconv.Atoi(limit)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s", err)
		}

		rows, err := GetSessions(db, int64(pageNum), int64(limitNum))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s", err)
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(rows)
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func nc_server(db *sql.DB) {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error in conn: %s", err)
		} else {
			go func(conn net.Conn) {
				start := time.Now().Unix()
				addr := conn.RemoteAddr().String()
				fname := RandomString(32)

				err = InsertSession(db, addr, start, fname)
				if err != nil {
					log.Printf("Error in db: %s", err)
					return
				}

				f, err := os.Create(fmt.Sprintf("./data/%s", fname))
				if err != nil {
					log.Printf("Error in file: %s", err)
					return
				}

				_, err = io.Copy(f, conn)
				if err != nil {
					log.Printf("Error in copy: %s", err)
					return
				}

				err = conn.Close()
				if err != nil {
					log.Printf("Error in conn: %s", err)
					return
				}
			}(conn)
		}
	}
}

func RandomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
