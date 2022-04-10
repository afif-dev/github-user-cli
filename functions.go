package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/schollz/progressbar/v3"
)

func getData(url string) []byte {
	res, err := http.Get(url)
	res.Header.Add("Accept", "application/json")
	// res.Header.Add("User-Agent", "Qoutes CLI")

	if err != nil {
		log.Println(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	bar := progressbar.DefaultBytes(res.ContentLength, "fetching data")
	for i := 0; i < 100; i++ {
		bar.Add(1)
		time.Sleep(10 * time.Millisecond)
	}
	
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatalln(err)
	}
	return body
}
