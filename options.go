package khash

//NumOfReplica set replica number.
//This is option function to pass as argument for New function
func NumOfReplica(n int) func(*Khash) error {
	return func(k *Khash) (err error) {
		if k.numberOfReplica != 0 {
			return errReplicaNotEmpty
		}
		k.numberOfReplica = n
		return
	}
}

//Node set nodes to initialize.
//This is option function to pass as argument for New function
func Node(nodes []string) func(*Khash) error {
	return func(k *Khash) (err error) {
		if len(nodes) == 0 {
			return errEmptyStringArr
		}
		if k.numberOfReplica == 0 {
			k.numberOfMembers = defaultNumberReplica
		}
		k.ring = make(map[uint32]string)
		k.members = make(map[string]bool)
		for i := range nodes {
			k.add(nodes[i])
		}
		return
	}
}
