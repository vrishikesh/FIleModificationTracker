package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/osquery/osquery-go"
	"github.com/vrishikesh/FileModificationTracker/internal/model"
)

func main() {
	db, err := sqlx.Connect("sqlite3", "logs.db")
	if err != nil {
		log.Fatal(err)
	}

	dir, err := os.Getwd()
	log.Printf("current working dir: %s\n", dir)
	if err != nil {
		log.Fatal(err)
	}

	client, err := osquery.NewClient("/Users/rishikeshvishwakarma/.osquery/shell.em", 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	go fetchFileModifications(dir, client, db)

	serve := http.Server{}
	serve.Addr = ":8081"
	serve.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("url: %s\n", r.Method)
		log.Printf("%s\n", r.URL.Path)

		if r.Method == http.MethodGet && r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
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

func fetchFileModifications(directory string, client *osquery.ExtensionManagerClient, db *sqlx.DB) {
	for {
		log.Printf("executing query")
		query := fmt.Sprintf("SELECT * FROM file WHERE path LIKE '%s%%%%' AND mtime >= strftime('%%s', 'now') - 60", directory)
		results, err := client.Query(query)
		log.Printf("query result: %v", results.GetResponse())
		if err != nil {
			log.Fatal(err)
		}

		// Process and log results
		for _, row := range results.GetResponse() {
			log.Printf("File Modified: %v", row)

			l := model.Log{}
			l.Path = row["path"]
			l.Directory = row["directory"]
			l.Filename = row["filename"]
			l.Inode = row["inode"]
			l.Uid = row["uid"]
			l.Gid = row["gid"]
			l.Mode = row["mode"]
			l.Device = row["device"]
			l.Size = row["size"]
			l.BlockSize = row["block_size"]
			l.Atime = row["atime"]
			l.Mtime = row["mtime"]
			l.Ctime = row["ctime"]
			l.Btime = row["btime"]
			l.HardLinks = row["hard_links"]
			l.Symlink = row["symlink"]
			l.Type = row["type"]

			_, err := db.NamedExec("INSERT INTO logs(path, directory, filename, inode, uid, gid, mode, device, size, block_size, atime, mtime, ctime, btime, hard_links, symlink, type) VALUES (:path, :directory, :filename, :inode, :uid, :gid, :mode, :device, :size, :block_size, :atime, :mtime, :ctime, :btime, :hard_links, :symlink, :type)", &l)
			if err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
