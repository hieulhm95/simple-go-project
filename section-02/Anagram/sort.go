package main

import (
	"fmt"
	"sort"
	"strings"
)

func isAnagram(s string, t string) bool {
    cloneS := strings.Split(s, "")
	sort.Strings(cloneS)
    cloneT := strings.Split(t, "")
    sort.Strings(cloneT)
    fmt.Println(cloneS, cloneT)
    for i := 0; i < len(cloneS); i++ {
        if cloneS[i] != cloneT[i] {
                return false
        }
    }
    return true
}

func Main() {
	fmt.Println(isAnagram("anagram", "nagaram"))
}
