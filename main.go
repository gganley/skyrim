// author: Gregory Ganley

package main

import (
	"fmt"
	"math/bits"
	"os"
	"sync"
)

// What needs to happen
// store known combinations in a way that can be atomically accessed

type Ingredient struct{ disc, all uint64 } // Disc holds all of the
type Effect uint64                         // 2^i where i is the index of the effect

// This is my ever nice term for the return type of f()
type FunkyType struct {
	NumberOfDiscoveredEffects int
	Potion                    uint64
	IngredientNames           []string
}

// The function that determines all the crucial information
func f(ings []string, Im map[string]*Ingredient) FunkyType {
	// Determine the complete possible potion
	gamma := Im[ings[0]].all&Im[ings[1]].all |
		Im[ings[0]].all&Im[ings[2]].all |
		Im[ings[1]].all&Im[ings[2]].all

	a := Im[ings[0]].disc & gamma
	b := Im[ings[1]].disc & gamma
	c := Im[ings[2]].disc & gamma

	// All the ingredients that will be discovered
	q := bits.OnesCount64(a) + bits.OnesCount64(b) + bits.OnesCount64(c)
	return FunkyType{q, gamma, ings}
}

func main() {
	// `Im` is the Inredient Map, basically because reading binary data gets old after a while
	Im, _ := discoverImEm()

	// Get/Make the non-dud potions
	_, err := os.Open("nonduds.csv")
	if os.IsNotExist(err) {
		WriteNonDuds(Im)
	}
	nonDuds := getNonDuds()

	// The list of potions that determine all the ingredients
	var retVal [][]string

	for effectsLeft(Im) {
		// The channel that all the results of each iteration of f() is sent to
		result := make(chan FunkyType, 8500)

		// This is a sync utility that lets me manage the lifetime of the channel by incrementing it when a new process
		// is added and decrementing it when its done calculating it
		var wg sync.WaitGroup
		for _, v := range nonDuds {
			wg.Add(1)
			val := v
			go func(x []string, i map[string]*Ingredient) {
				defer wg.Done()
				result <- f(x, i)
			}(val, Im)
		}

		// This runs when all the goroutines that run f() have finished, then contines to the next for-loop
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
			// Turn off the discovered bits in each ingredient
			Im[v].disc = Im[v].disc &^ potionToBeSubtracted
		}

		// Append the found best potion to the return value
		retVal = append(retVal, maxCombo)
	}

	// Printing the result
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
