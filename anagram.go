package main

import (
	"log"
	"sort"
)

func main() {
	words := []string{"kita", "atik", "tika", "aku", "kia", "makan", "kua"}
	refs := map[string][]string{}

	for _, word := range words {
		r := []rune(word)
		sort.Slice(r, func(i, j int) bool { return r[i] < r[j] })
		sortedWord := string(r)

		refs[sortedWord] = append(refs[sortedWord], word)
	}

	// transforming into a slice of slices
	result := [][]string{}
	for _, item := range refs {
		result = append(result, item)
	}

	log.Println(result)
}
