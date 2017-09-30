package khash

import (
	"sync"
)

type (
	uint32Arr []uint32

	//Khash is the main data structure of khash
	Khash struct {
		sync.RWMutex
		ring            map[uint32]string
		members         map[string]bool
		sortedHashes    uint32Arr
		numberOfReplica int
		numberOfMembers int
	}
)

// Len returns the length of the uints array.
func (u uint32Arr) Len() int { return len(u) }

// Less returns true if element i is less than element j.
func (u uint32Arr) Less(i, j int) bool { return u[i] < u[j] }

// Swap exchanges elements i and j.
func (u uint32Arr) Swap(i, j int) { u[i], u[j] = u[j], u[i] }
