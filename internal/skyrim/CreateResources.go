package skyrim

import (
	"bufio"
	"encoding/csv"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"

	"github.com/klaidliadon/next"
)

func potion(i1, i2, i3 *Ingredient) int {
	comb1 := i1.All & i2.All
	comb2 := i1.All & i3.All
	comb3 := i2.All & i3.All

	if comb1 > 0 && comb2 > 0 && comb3 > 0 {
		return int(comb1 | comb2 | comb3)
	}
	return 0
}

func GetNonDuds() [][]string {
	nonDuds, err := os.Open("nonduds.csv")

	if err != nil {
		log.Fatal(err)
	}
	defer nonDuds.Close()

	r := csv.NewReader(nonDuds)

	records, e := r.ReadAll()

	if e != nil {
		log.Fatal(e)
	}

	return records
}

func WriteNonDuds(Im map[string]*Ingredient) {
	ings := make([]interface{}, 0, len(Im))
	for k := range Im {
		ings = append(ings, k)
	}
	var combs [][]interface{}
	for v := range next.Combination(ings, 3, false) {
		combs = append(combs, v)
	}

	var combstr [][]string

	for _, v := range combs {
		var temp []string
		for _, val := range v {
			temp = append(temp, val.(string))
		}
		combstr = append(combstr, temp)
	}
	var pus [][]string
	for _, v := range combstr {
		pot := potion(Im[v[0]], Im[v[1]], Im[v[2]])
		if pot > 0 {
			pus = append(pus, v)
		}
	}

	nonDuds, err := os.Create("nonduds.csv")

	if err != nil {
		log.Fatal(err)
	}
	defer nonDuds.Close()

	w := csv.NewWriter(nonDuds)
	writeErr := w.WriteAll(pus)

	if writeErr != nil {
		log.Fatal(writeErr)
	}

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func DiscoverImEm() (Im map[string]*Ingredient, Em map[string]*Effect) {
	records := slurpFile()

	Im = make(map[string]*Ingredient)
	Em = make(map[string]*Effect)

	effFile, err := os.Open("eff.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer effFile.Close()

	scanner := bufio.NewScanner(effFile)
	i := 0
	for scanner.Scan() {
		calc := Effect(math.Pow(2, float64(i)))
		Em[scanner.Text()] = &calc
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for _, v := range records {
		var num uint64 = 0
		for _, eff := range v[1:5] {
			num += uint64(*Em[eff])
		}
		Im[v[0]] = &Ingredient{num, num}
	}

	// This is some weird Go syntax, why do I need this return statement, you hit the fucking end of the function just
	// return the return params like jesus
	return
}

func slurpFile() [][]string {

	file, err := ioutil.ReadFile("ing.csv")
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(strings.NewReader(string(file)))

	records, err := r.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	return records
}

func fauxMain() {
	Im, _ := DiscoverImEm()

	WriteNonDuds(Im)

}
