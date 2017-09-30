// Khash is consistent hashing implementation
// It uses to distribute load or sharding mechanism on multiple nodes
//
// The advantage compared to typical hash table is it doesn't need to remap all key
// after node addition or removal

package khash

import (
	"errors"
	"hash/crc32"
	"log"
	"sort"
	"strconv"
)

const defaultNumberReplica = 20

var (
	errEmptyStringArr  = errors.New("empty array string")
	errEmptyString     = errors.New("empty string")
	errEmptyRing       = errors.New("ring is empty")
	errReplicaNotEmpty = errors.New("replica is not empty. set replica before Node")
)

//New create khash instance
func New(options ...func(*Khash) error) (newInstance *Khash) {
	newInstance = &Khash{}

	//set options
	for _, option := range options {
		err := option(newInstance)
		if err != nil {
			log.Println("couldn't set option", err)
		}
	}

	//initialize when no ring and members set
	if newInstance.ring == nil {
		newInstance.ring = make(map[uint32]string)
		newInstance.members = make(map[string]bool)
	}

	//initialize when no replica set
	if newInstance.numberOfReplica == 0 {
		newInstance.numberOfReplica = defaultNumberReplica
	}

	return
}

//Add node to the ring
func (k *Khash) Add(node string) (err error) {
	k.Lock()
	err = k.add(node)
	k.Unlock()

	return
}

// need c.Lock() before calling
func (k *Khash) add(node string) (err error) {
	if node == "" {
		err = errEmptyString
		return
	}

	for i := 0; i < k.numberOfReplica; i++ {
		k.ring[k.hashKey(k.constructReplicaKey(node, i))] = node
	}
	k.members[node] = true
	k.updateSortedHashes()
	k.numberOfMembers++

	return
}

//Remove node from the ring
func (k *Khash) Remove(node string) (err error) {
	k.Lock()
	err = k.remove(node)
	k.Unlock()

	return
}

// need c.Lock() before calling
func (k *Khash) remove(node string) (err error) {
	for i := 0; i < k.numberOfReplica; i++ {
		delete(k.ring, k.hashKey(k.constructReplicaKey(node, i)))
	}
	delete(k.members, node)
	k.updateSortedHashes()
	k.numberOfMembers--

	return
}

// Get returns an node close to where name hashes to in the ring.
func (k *Khash) Get(name string) (node string, err error) {
	k.RLock()
	defer k.RUnlock()

	if len(k.ring) == 0 {
		err = errEmptyRing
		return
	}
	i := k.search(k.hashKey(name))
	node = k.ring[k.sortedHashes[i]]

	return
}

func (k *Khash) search(key uint32) (i int) {
	i = sort.Search(len(k.sortedHashes), func(x int) bool {
		return k.sortedHashes[x] > key
	})
	if i >= len(k.sortedHashes) {
		i = 0
	}
	return
}

func (k *Khash) updateSortedHashes() {
	// new allocation
	k.sortedHashes = nil
	for r := range k.ring {
		k.sortedHashes = append(k.sortedHashes, r)
	}
	sort.Sort(k.sortedHashes)
}

func (k *Khash) constructReplicaKey(key string, index int) (finalKey string) {
	finalKey = strconv.Itoa(index) + key
	return
}

func (k *Khash) hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}
