package main

import (
	"log"
	"os"

	"github.com/denysvitali/ble"
)

func main() {
	conn, err := ble.Open()
	if err != nil {
		log.Fatal(err)
	}
	conn.Print(os.Stdout)
}
