// author: Gregory Ganley

package skyrim

import (
	"math/bits"
	// "github.com/klaidliadon/next"
)

// What needs to happen
// store known combinations in a way that can be atomically accessed

// Data structure section
type Ingredient struct{ Disc, All uint64 } // Disc holds all of the
type Effect uint64                         // 2^i where i is the index of the effect

// This is my ever nice term for the return type of evaluatePotion()
type PotionStructure struct {
	NumberOfDiscoveredEffects int
	Potion                    uint64
	IngredientNames           []string
}

// The function that determines all the crucial information
func evaluatePotion(ings []string, Im map[string]*Ingredient) PotionStructure {
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

func findBestPotion(result chan PotionStructure) map[int][][]string {
	// Maybe return a map of map[numberOfDiscoveredEffects]Potion and simple look at the largest numofdisc...
	bestPotionMap := make(map[int][][]string)

	for v := range result {
		bestPotionMap[v.NumberOfDiscoveredEffects] = append(bestPotionMap[v.NumberOfDiscoveredEffects], v.IngredientNames)
	}

	return bestPotionMap
}

func effectsLeft(Im map[string]*Ingredient) bool {
	for _, v := range Im {
		if v.Disc > 0 {
			return true
		}
	}

	return false
}
