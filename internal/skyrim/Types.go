// author: Gregory Ganley

package skyrim

type Ingredient struct{ Disc, All uint64 } // Disc holds all of the
type Effect uint64                         // 2^i where i is the index of the effect

type PotionStructure struct {
	NumberOfDiscoveredEffects int
	Potion                    uint64
	IngredientNames           []string
}
