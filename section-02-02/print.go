package main

import (
	"fmt"
	"runtime"
)

func init() {
	fmt.Println("This is init a")
}

func printAllocation() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%d KB\n", m.Alloc/1024)

}
