package module

import (
	"fmt"
	"testing"
)

func TestShowLLDPNeighbors_UnitTest(t *testing.T) {
	show := Show(dummyNode)
	neighbors := show.ShowLLDPNeighbors()

	fmt.Printf("%+v", neighbors)
}
