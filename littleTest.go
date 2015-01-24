package main

import (
	"fmt"
	"os"

	"github.com/gophergala/edrans-smartcity/algorithm"
)

func main() {
	city := algorithm.GetTestCity()
	paths, e := city.GetPaths(6, 10)
	if e != nil {
		fmt.Printf("Ohh no... %+v\n", e)
		os.Exit(2)
	}
	vehicle := algorithm.CallService("doctor")
	vehicle.InCity = &city
	elapsedTime := vehicle.Run(algorithm.OrderCandidates(vehicle.CalcPaths(paths))[0])
	fmt.Println("elapsed time:", elapsedTime)
	/*for i := 0; i < 1; i++ {
	  for j := 0; j < len(paths[0].Links); j++ {
	    fmt.Printf("Link #%d: %+v\n", j, paths[0].Links[j])
	  }
	}*/
}
