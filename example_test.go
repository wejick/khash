package khash

import (
	"fmt"
	"log"
)

func ExampleNew() {
	k := New()
	k.Add("cacheA")
	k.Add("cacheB")
	k.Add("cacheC")
	users := []string{"user_mcnulty", "user_bunk", "user_omar", "user_bunny", "user_stringer"}
	for _, u := range users {
		node, err := k.Get(u)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s => %s\n", u, node)
	}
	// Output:
	// user_mcnulty => cacheA
	// user_bunk => cacheA
	// user_omar => cacheA
	// user_bunny => cacheC
	// user_stringer => cacheC
}

func ExampleKhash_Add() {
	k := New()
	k.Add("cacheA")
	k.Add("cacheB")
	k.Add("cacheC")
	users := []string{"user_mcnulty", "user_bunk", "user_omar", "user_bunny", "user_stringer"}
	fmt.Println("initial state [A, B, C]")
	for _, u := range users {
		server, err := k.Get(u)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s => %s\n", u, server)
	}
	k.Add("cacheD")
	k.Add("cacheE")
	fmt.Println("\nwith cacheD, cacheE [A, B, C, D, E]")
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
	// with cacheD, cacheE [A, B, C, D, E]
	// user_mcnulty => cacheE
	// user_bunk => cacheA
	// user_omar => cacheA
	// user_bunny => cacheE
	// user_stringer => cacheE
}

func ExampleKhash_Remove() {
	k := New()
	k.Add("cacheA")
	k.Add("cacheB")
	k.Add("cacheC")
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
}
