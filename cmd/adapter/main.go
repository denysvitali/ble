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
	adapters, err := conn.GetAdapters("aaa")

	for _, adapter := range(adapters) {
		if err != nil {
			log.Fatal(err)
		}
		adapter.Print(os.Stdout)
	}
}
