package main

import (
	"fmt"
	"os"

	"github.com/gophergala/edrans-smartcity/algorithm"
)

func main() {
	city := algorithm.GetTestCity()
	paths, e := city.GetPath(15, 2)
	if e != nil {
		fmt.Printf("Ohh no... %+v\n", e)
		os.Exit(2)
	}
	for i := 0; i < len(paths); i++ {
		fmt.Printf("\nOption #%d: %+v\n", i+1, paths[i])
	}
}
