//MIAUSEproject-founderKJ/multi-platform-AI/simulation/replay/twist_injector.go

package replay

import (
	"multi-platform-AI/simulation/digital_twin"
	"math/rand"
)

type Twist struct {
	Type  string  // "FrictionLoss", "SensorNoise", "GhostObstacle"
	Value float64
}

func InjectTwist(world *digital_twin.VoxelWorld) Twist {
	// Randomly degrade the simulation to test the Student Policy's robustness
	twists := []string{"FrictionLoss", "SensorNoise", "GhostObstacle"}
	selected := twists[rand.Intn(len(twists))]
	
	t := Twist{Type: selected, Value: rand.Float64()}
	
	if t.Type == "GhostObstacle" {
		// Insert a fake wall into the voxel grid to see if AI avoids it
		world.Grid["ghost_obj"] = digital_twin.VoxelCell{Occupied: true, Hardness: 1.0}
	}
	
	return t
}