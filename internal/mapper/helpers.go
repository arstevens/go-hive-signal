package mapper

import (
	"math"

	"github.com/arstevens/go-hive-signal/internal/transmuter"
)

type swarmMapPair struct {
	Manager    transmuter.SwarmManager
	Dataspaces []string
}

func updateMinDataspacesSwarmInfo(sm *SwarmMap, newID string, newSize int) {
	minID := calculateMinDataspaceSwarm(sm)
	if newSize < len(sm.managerMap[minID].Dataspaces) {
		sm.minDataspacesSwarmID = newID
		sm.minDataspaceSwarmSize = newSize
	}
}

func calculateMinDataspaceSwarm(sm *SwarmMap) string {
	if sm.minDataspacesSwarmID != "" {
		return sm.minDataspacesSwarmID
	}
	minSize := math.MaxInt32
	minID := ""
	for id, pair := range sm.managerMap {
		if len(pair.Dataspaces) < minSize {
			minSize = len(pair.Dataspaces)
			minID = id
		}
	}
	return minID
}
