package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"sort"
	"strings"
)

func RemoveDuplicates(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}

func main() {
	file, err := ioutil.ReadFile("ingredients.csv")
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(strings.NewReader(string(file)))
	r.Comma = ';'

	fullMap := make(map[string][]string)
	funkyMask := make(map[string]uint64)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fullMap[record[0]] = record[1:5]
	}

	var effects []string
	var ingredients []string

	for k, v := range fullMap {
		effects = append(effects, v...)
		ingredients = append(ingredients, k)
	}

	RemoveDuplicates(&effects)

	sort.Strings(effects)

	for i, v := range effects {
		funkyMask[v] = uint64(math.Pow(2, float64(i)))
	}

	finalMaskingThing := make(map[string]uint64)

	for k, v := range fullMap {
		var sum uint64 = 0

		for _, entry := range v {
			sum += funkyMask[entry]
		}

		finalMaskingThing[k] = sum
	}
	fmt.Println(fullMap)
	w := csv.NewWriter(os.Stdout)
	for k, v := range fullMap {
		if err := w.Write(append([]string{k}, v...)); err != nil {
			log.Fatalln("error: ", err)
		}
	}

	w.Flush()
}
