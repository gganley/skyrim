// author: Gregory Ganley

package main

import (
	"fmt"
	"math/bits"
	"sync"
)

// What needs to happen
// store known combinations in a way that can be atomically accessed

type Ingredient struct{ disc, all uint64 } // Disc holds all of the
type Effect uint64                         // 2^i where i is the index of the effect

type FunkyType struct {
	NumberOfDiscoveredEffects int
	Potion                    uint64
	IngredientNames           []string
}

func f(ings []string, Im map[string]*Ingredient) FunkyType {
	gamma := Im[ings[0]].all&Im[ings[1]].all |
		Im[ings[0]].all&Im[ings[2]].all |
		Im[ings[1]].all&Im[ings[2]].all

	a := Im[ings[0]].disc & gamma
	b := Im[ings[1]].disc & gamma
	c := Im[ings[2]].disc & gamma

	q := bits.OnesCount64(a) + bits.OnesCount64(b) + bits.OnesCount64(c)
	return FunkyType{q, gamma, ings}
}

func main() {

	Im, _ := discoverImEm()
	WriteNonDuds(Im)
	nonDuds := getNonDuds()

	var retVal [][]string

	for effectsLeft(Im) {
		result := make(chan FunkyType, 8500)
		var wg sync.WaitGroup
		for _, v := range nonDuds {
			wg.Add(1)
			val := v
			go func(x []string, i map[string]*Ingredient) {
				defer wg.Done()
				result <- f(x, i)
			}(val, Im)
		}

		go func() {
			wg.Wait()
			close(result)
		}()

		max := 0
		var maxCombo []string
		potionToBeSubtracted := uint64(0)
		for v := range result {
			if v.NumberOfDiscoveredEffects > max {
				max = v.NumberOfDiscoveredEffects
				maxCombo = v.IngredientNames
				potionToBeSubtracted = v.Potion
			}
		}
		for _, v := range maxCombo {
			Im[v].disc = Im[v].disc &^ potionToBeSubtracted
		}

		retVal = append(retVal, maxCombo)
	}
	// I really need to figure out if I want to index by the string name of the ingredient or weirdly by the value it
	// holds

	// ingsmap := make(map[string]int)
	//
	// for _, v := range retVal {
	// 	for _, val := range v {
	// 		ingsmap[val]++
	// 	}
	// 	fmt.Println(v)
	// }
	// total := 0
	// for k, v := range ingsmap {
	// 	total += v
	// 	fmt.Println(k, v)
	// }
	// fmt.Println(total)

	for _, v := range retVal {
		fmt.Println(v)
	}

	fmt.Println(len(retVal))
}

func effectsLeft(Im map[string]*Ingredient) bool {
	for _, v := range Im {
		if v.disc > 0 {
			return true
		}
	}

	return false
}
