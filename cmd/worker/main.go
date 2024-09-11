package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/vrishikesh/FileModificationTracker/internal/model"
)

type AcceptCommand struct {
	Command string
}

func worker(jobs <-chan string) {
	for j := range jobs {
		app := strings.Fields(j)
		cmd := exec.Command(app[0], app[1:]...)
		stdout, err := cmd.Output()
		if err != nil {
			log.Println(err.Error())
		}
		log.Println(string(stdout))
	}
}

func main() {
	db, err := sqlx.Connect("sqlite3", "logs.db")
	if err != nil {
		log.Fatal(err)
	}

	jobs := make(chan string, 10)

	go worker(jobs)

	serve := http.Server{}
	serve.Addr = ":8080"
	serve.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("url: %s\n", r.Method)
		log.Printf("%s\n", r.URL.Path)

		if r.Method == http.MethodGet && r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodGet && r.URL.Path == "/logs" {
			logs := []model.Log{}
			err := db.Select(&logs, "SELECT * FROM logs LIMIT 100")
			if err != nil {
				if err == sql.ErrNoRows {
					log.Println("no rows in db")
				}
				log.Fatal(err)
			}

			b, err := json.Marshal(logs)
			if err != nil {
				log.Println(err)
			}

			w.Header().Set("Content-Type", "application/json")
			_, err = w.Write(b)
			if err != nil {
				log.Fatal(err)
			}

			return
		}

		if r.Method == http.MethodPost && r.URL.Path == "/command" {
			b, err := io.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
				return
			}
			r.Body.Close()

			ac := AcceptCommand{}
			err = json.Unmarshal(b, &ac)
			if err != nil {
				log.Println(err)
				return
			}

			jobs <- ac.Command

			_, err = w.Write([]byte("command queued"))
			if err != nil {
				log.Fatal(err)
			}

			return
		}

		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(http.StatusText(http.StatusNotFound)))
		if err != nil {
			log.Fatal(err)
		}
	})

	if err := serve.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
