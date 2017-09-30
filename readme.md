[![Build Status](https://travis-ci.org/wejick/khash.svg?branch=master)](https://travis-ci.org/wejick/khash) [![Coverage Status](https://coveralls.io/repos/github/wejick/khash/badge.svg?branch=master)](https://coveralls.io/github/wejick/khash?branch=master)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/wejick/khash)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/wejick/khash/blob/master/LICENSE)

# KHash - consistent hashing implementation in go

Khash is based on stathat implementation

# Feature
1. interface for periodic node healtcheck               [todo]
1. automatic node invalidation / removal                [todo]
1. hook for node removal / addition alerting            [todo]

# Difference
1. Some feature like `GetN` are not implemented in khash
1. Slightly faster than stathat implementation (benchmark result at the end of this page)
1. Doesn't have `Set` function, however set on initialization is supported via functional option

# License
This project is licensed under MIT license you can find in LICENSE file

# Get Started
```
        k := New(NumOfReplica(25),
                Node([]string{"cacheA", "cacheB", "cacheC"}))
        users := []string{"user_mcnulty", "user_bunk", "user_omar", "user_bunny", "user_stringer"}
	fmt.Println("initial state [A, B, C]")
	for _, u := range users {
		server, err := k.Get(u)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s => %s\n", u, server)
	}
	k.Remove("cacheC")
	fmt.Println("\ncacheC removed [A, B]")
	for _, u := range users {
		server, err := k.Get(u)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s => %s\n", u, server)
	}
	// Output:
	// initial state [A, B, C]
	// user_mcnulty => cacheA
	// user_bunk => cacheA
	// user_omar => cacheA
	// user_bunny => cacheC
	// user_stringer => cacheC
	//
	// cacheC removed [A, B]
	// user_mcnulty => cacheA
	// user_bunk => cacheA
	// user_omar => cacheA
	// user_bunny => cacheB
	// user_stringer => cacheB

```

# More to read
1. https://www.akamai.com/us/en/multimedia/documents/technical-publication/consistent-hashing-and-random-trees-distributed-caching-protocols-for-relieving-hot-spots-on-the-world-wide-web-technical-publication.pdf (the original paper from akamai)
1. http://www.martinbroadhurst.com/consistent-hash-ring.html (more about replicas in khash)

# Benchmark

```
goos: darwin
goarch: amd64
pkg: github.com/stathat/consistent
BenchmarkAllocations-4            100000             14856 ns/op
--- BENCH: BenchmarkAllocations-4
        consistent_test.go:551: 1: Allocated 4480 bytes (4480.00x)
        consistent_test.go:551: 100: Allocated 283264 bytes (2832.64x)
        consistent_test.go:551: 10000: Allocated 28161888 bytes (2816.19x)
        consistent_test.go:551: 100000: Allocated 281602400 bytes (2816.02x)
BenchmarkMalloc-4                 100000             15403 ns/op
--- BENCH: BenchmarkMalloc-4
        consistent_test.go:564: 1: Mallocd 86 times (86.00x)
        consistent_test.go:564: 100: Mallocd 8202 times (82.02x)
        consistent_test.go:564: 10000: Mallocd 820003 times (82.00x)
        consistent_test.go:564: 100000: Mallocd 8200003 times (82.00x)
BenchmarkCycle-4                  100000             17329 ns/op
BenchmarkCycleLarge-4              20000             70178 ns/op
BenchmarkGet-4                  10000000               203 ns/op
BenchmarkGetLarge-4              5000000               219 ns/op
PASS
ok      github.com/stathat/consistent   24.397s

goos: darwin
goarch: amd64
pkg: github.com/wejick/khash
BenchmarkAllocations-4            100000             16122 ns/op
--- BENCH: BenchmarkAllocations-4
        benchmark_test.go:40: 1: Allocated 3040 bytes (3040.00x)
        benchmark_test.go:40: 100: Allocated 147008 bytes (1470.08x)
        benchmark_test.go:40: 10000: Allocated 14561696 bytes (1456.17x)
        benchmark_test.go:40: 100000: Allocated 145601616 bytes (1456.02x)
BenchmarkMalloc-4                 100000             16522 ns/op
--- BENCH: BenchmarkMalloc-4
        benchmark_test.go:53: 1: Mallocd 94 times (94.00x)
        benchmark_test.go:53: 100: Mallocd 9302 times (93.02x)
        benchmark_test.go:53: 10000: Mallocd 930001 times (93.00x)
        benchmark_test.go:53: 100000: Mallocd 9300001 times (93.00x)
BenchmarkCycle-4                  100000             18659 ns/op
BenchmarkCycleLarge-4              20000             78122 ns/op
BenchmarkGet-4                  10000000               211 ns/op
BenchmarkGetLarge-4              5000000               222 ns/op
PASS
ok      github.com/wejick/khash 13.477s
```