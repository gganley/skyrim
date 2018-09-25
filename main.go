package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	file, err := ioutil.ReadFile("ing.csv")
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(strings.NewReader(string(file)))

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(records)
}
