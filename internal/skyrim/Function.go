package skyrim

import (
	"math/bits"
	"os"
	"sync"
)

func FullCycle(Im map[string]*Ingredient, nonDuds [][]string) int {

	// The list of potions that determine all the ingredients
	var retThing [][]int
	var iteration int

	for effectsLeft(Im) {
		// determine best canadates for this iteration
		// in the case of 1 best canadate just remove that potion and continue
		// if there is more than one best canadate then split off each of those canadates onto it's own thread
		// once they've all reterened determine which one is the best
		// return that best
		//
		//

		potions := SingleCycle(Im, nonDuds)
		// Break off on all potions
		if len(potions) == 1 {
			RemovePotion(Im, potions[0])
			continue
		} else {
			for potion := range potions {
				newIm := copyMap(Im)
				RemovePotion(newIm, potion)
				retThing[iteration] = append(retThing[iteration], FullCycle(Im, nonDuds))
			}
		}
	}
}

func copyMap(Im map[string]*Ingredient) map[string]*Ingredient {
	newMap := make(map[string]*Ingredient)
	for k, v := range Im {
		*newMap[k] = *v
	}

	return newMap
}

func RemovePotion(Im map[string]*Ingredient, potion PotionStructure) {

}

func SetupCycle() (map[string]*Ingredient, [][]string) {
	// `Im` is the Inredient Map, basically because reading binary data gets old after a while
	Im, _ := DiscoverImEm()

	// Get/Make the non-dud potions
	_, err := os.Open("nonduds.csv")
	if os.IsNotExist(err) {
		WriteNonDuds(Im)
	}
	nonDuds := GetNonDuds()
	return Im, nonDuds
}

func SingleCycle(Im map[string]*Ingredient, nonDuds [][]string) []PotionStructure {
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

	// This is where we determine whose best, need to care for equal best
	// var max int
	// var maxCombo [][]string
	// var potionToBeSubtracted uint64
	maxMap := make(map[int][]PotionStructure)

	for v := range result {
		maxMap[v.NumberOfDiscoveredEffects] = append(maxMap[v.NumberOfDiscoveredEffects], v)
	}

	var maxKey int
	var finalCombo []PotionStructure
	for key, val := range maxMap {
		if key > maxKey {
			finalCombo = val
		}
	}

	// for _, v := range maxCombo {
	// 	// Turn off the discovered bits in each ingredient
	// 	Im[v].Disc = Im[v].Disc &^ potionToBeSubtracted
	// }
	// Leave the subtraction for later

	// Append the found best potion to the return value
	return finalCombo
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
