package fsm

import (
	"math"
	."globals"
)

func distance(from int, to int) int{
	return (to-from)*int(calculateDir(to, from))
}

func calculateDir(destination int, currentFloor int) Direction_t {
	if destination == currentFloor{
		return DIR_STOP
	}
	return Direction_t(math.Copysign(1, float64(destination-currentFloor)))
}