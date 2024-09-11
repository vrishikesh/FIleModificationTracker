package logging

import (
	"log"
	"os"
)

func InitLogger() {
	file, err := os.OpenFile("file_tracker.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
}
