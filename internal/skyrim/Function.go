package skyrim

import (
	"math/bits"
	"os"
	"sync"
)

func SingleCycle() int {
	// `Im` is the Inredient Map, basically because reading binary data gets old after a while
	Im, _ := DiscoverImEm()

	// Get/Make the non-dud potions
	_, err := os.Open("nonduds.csv")
	if os.IsNotExist(err) {
		WriteNonDuds(Im)
	}
	nonDuds := GetNonDuds()

	// The list of potions that determine all the ingredients
	var retVal [][]string

	for effectsLeft(Im) {
		// The channel that all the results of each iteration of f() is sent to
		result := make(chan PotionStructure, 8500)

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
			Im[v].Disc = Im[v].Disc &^ potionToBeSubtracted
		}

		// Append the found best potion to the return value
		retVal = append(retVal, maxCombo)
	}

	return len(retVal)
}

func effectsLeft(Im map[string]*Ingredient) bool {
	for _, v := range Im {
		if v.Disc > 0 {
			return true
		}
	}

	return false
}

// The function that determines all the crucial information
func f(ings []string, Im map[string]*Ingredient) PotionStructure {
	// Determine the complete possible potion
	gamma := Im[ings[0]].All&Im[ings[1]].All |
		Im[ings[0]].All&Im[ings[2]].All |
		Im[ings[1]].All&Im[ings[2]].All

	a := Im[ings[0]].Disc & gamma
	b := Im[ings[1]].Disc & gamma
	c := Im[ings[2]].Disc & gamma

	// All the ingredients that will be discovered
	q := bits.OnesCount64(a) + bits.OnesCount64(b) + bits.OnesCount64(c)
	return PotionStructure{q, gamma, ings}
}
