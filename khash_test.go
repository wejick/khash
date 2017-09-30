package khash

import (
	"bufio"
	"os"
	"sort"
	"sync"
	"testing"
	"testing/quick"
)

type gtest struct {
	in  string
	out string
}

func checkNum(num, expected int, t *testing.T) {
	if num != expected {
		t.Errorf("got %d, expected %d", num, expected)
	}
}

func TestKhash_constructReplicaKey(t *testing.T) {
	type fields struct {
		RWMutex         sync.RWMutex
		ring            map[uint32]string
		members         map[string]bool
		numberOfReplica int
		numberOfMembers int
	}
	type args struct {
		key   string
		index int
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantFinalKey string
	}{
		{
			name:         "empty",
			wantFinalKey: "0",
		},
		{
			name: "one 1",
			args: args{
				key:   "one",
				index: 1,
			},
			wantFinalKey: "1one",
		},
		{
			name: "two 2",
			args: args{
				key:   "two",
				index: 2,
			},
			wantFinalKey: "2two",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Khash{
				RWMutex:         tt.fields.RWMutex,
				ring:            tt.fields.ring,
				members:         tt.fields.members,
				numberOfReplica: tt.fields.numberOfReplica,
				numberOfMembers: tt.fields.numberOfMembers,
			}
			if gotFinalKey := k.constructReplicaKey(tt.args.key, tt.args.index); gotFinalKey != tt.wantFinalKey {
				t.Errorf("Khash.constructReplicaKey() = %v, want %v", gotFinalKey, tt.wantFinalKey)
			}
		})
	}
}

func TestKhash_hashKey(t *testing.T) {
	type fields struct {
		RWMutex         sync.RWMutex
		ring            map[uint32]string
		members         map[string]bool
		numberOfReplica int
		numberOfMembers int
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint32
	}{
		{
			name: "empty",
			want: 0,
		},
		{
			name: "terima kasih bunda",
			args: args{
				key: "terima kasih bunda",
			},
			want: 2360902884,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Khash{
				RWMutex:         tt.fields.RWMutex,
				ring:            tt.fields.ring,
				members:         tt.fields.members,
				numberOfReplica: tt.fields.numberOfReplica,
				numberOfMembers: tt.fields.numberOfMembers,
			}
			if got := k.hashKey(tt.args.key); got != tt.want {
				t.Errorf("Khash.hashKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKhash_add(t *testing.T) {
	x := New()
	x.Add("abcdefg")
	checkNum(len(x.ring), 20, t)
	checkNum(len(x.sortedHashes), 20, t)
	if sort.IsSorted(x.sortedHashes) == false {
		t.Errorf("expected sorted hashes to be sorted")
	}
	x.Add("qwer")
	checkNum(len(x.ring), 40, t)
	checkNum(len(x.sortedHashes), 40, t)
	if sort.IsSorted(x.sortedHashes) == false {
		t.Errorf("expected sorted hashes to be sorted")
	}
}

func TestKhash_remove(t *testing.T) {
	var rtestsBefore = []gtest{
		{"ggg", "abcdefg"},
		{"hhh", "opqrstu"},
		{"iiiii", "hijklmn"},
	}

	var rtestsAfter = []gtest{
		{"ggg", "abcdefg"},
		{"hhh", "opqrstu"},
		{"iiiii", "opqrstu"},
	}

	x := New()

	//remove
	x.Add("abcdefg")
	x.Remove("abcdefg")
	checkNum(len(x.ring), 0, t)
	checkNum(len(x.sortedHashes), 0, t)

	//remove non existence
	x.Add("abcdefg")
	x.Remove("abcdefghijk")
	checkNum(len(x.ring), 20, t)

	//multiple remove
	x.Add("abcdefg")
	x.Add("hijklmn")
	x.Add("opqrstu")
	for i, v := range rtestsBefore {
		result, err := x.Get(v.in)
		if err != nil {
			t.Fatal(err)
		}
		if result != v.out {
			t.Errorf("%d. got %q, expected %q before rm", i, result, v.out)
		}
	}
	x.Remove("hijklmn")
	for i, v := range rtestsAfter {
		result, err := x.Get(v.in)
		if err != nil {
			t.Fatal(err)
		}
		if result != v.out {
			t.Errorf("%d. got %q, expected %q after rm", i, result, v.out)
		}
	}
}

func TestKhash_Get(t *testing.T) {
	gmtests := []gtest{
		{"ggg", "abcdefg"},
		{"hhh", "opqrstu"},
		{"iiiii", "hijklmn"},
	}

	x := New()

	//empty get
	_, err := x.Get("asdfsadfsadf")
	if err == nil {
		t.Errorf("expected error")
	}

	//single get
	x.Add("abcdefg")
	f := func(s string) bool {
		y, err := x.Get(s)
		if err != nil {
			t.Logf("error: %q", err)
			return false
		}
		t.Logf("s = %q, y = %q", s, y)
		return y == "abcdefg"
	}
	if err := quick.Check(f, nil); err != nil {
		t.Fatal(err)
	}

	// multi get
	x.Add("abcdefg")
	x.Add("hijklmn")
	x.Add("opqrstu")
	for i, v := range gmtests {
		result, err := x.Get(v.in)
		if err != nil {
			t.Fatal(err)
		}
		if result != v.out {
			t.Errorf("%d. got %q, expected %q", i, result, v.out)
		}
	}

}

//test node option
func TestKhash_Node(t *testing.T) {
	gmtests := []gtest{
		{"ggg", "abcdefg"},
		{"hhh", "opqrstu"},
		{"iiiii", "hijklmn"},
	}

	x := New(Node([]string{"abcdefg", "hijklmn", "opqrstu"}))

	// multi get
	x.Add("abcdefg")
	x.Add("hijklmn")
	x.Add("opqrstu")
	for i, v := range gmtests {
		result, err := x.Get(v.in)
		if err != nil {
			t.Fatal(err)
		}
		if result != v.out {
			t.Errorf("%d. got %q, expected %q", i, result, v.out)
		}
	}
}

func TestCollisionsCRC(t *testing.T) {
	// t.SkipNow()
	c := New(NumOfReplica(10))
	f, err := os.Open("/usr/share/dict/words")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	found := make(map[uint32]string)
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		word := scanner.Text()
		for i := 0; i < c.numberOfReplica; i++ {
			ekey := c.constructReplicaKey(word, i)
			// ekey := word + "|" + strconv.Itoa(i)
			k := c.hashKey(ekey)
			exist, ok := found[k]
			if ok {
				t.Logf("found collision: %s, %s", ekey, exist)
				count++
			} else {
				found[k] = ekey
			}
		}
	}
	t.Logf("number of collisions: %d", count)
}

// from @edsrzf on github:
func TestAddCollision(t *testing.T) {
	// These two strings produce several crc32 collisions after "|i" is
	// appended added by Consistent.eltKey.
	const s1 = "abear"
	const s2 = "solidiform"
	x := New()
	x.Add(s1)
	x.Add(s2)
	elt1, err := x.Get("abear")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	y := New()
	// add elements in opposite order
	y.Add(s2)
	y.Add(s1)
	elt2, err := y.Get(s1)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if elt1 != elt2 {
		t.Error(elt1, "and", elt2, "should be equal")
	}
}
