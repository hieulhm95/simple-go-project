package main

import (
	"fmt"
)

func isAnagram(s string, t string) bool {
    tmp := make(map[int32]int)
	for _, c := range s {
		tmp[c]++
	}
    for _, c := range t {
        tmp[c]--
    }
    fmt.Println(tmp)
    for _, t := range tmp {
        if t != 0 {
            return false
        }
    }
    return true
}

func Main() {
	fmt.Println(isAnagram("anagram", "nagaram"))
}
