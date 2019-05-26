// author: Gregory Ganley

package main

import (
	"fmt"

	"github.com/gganley/skyrim/internal/skyrim"
)

func main() {
	Im, nonDuds := skyrim.SetupCycle()
	fmt.Println(skyrim.FullCycle(Im, nonDuds))
}
