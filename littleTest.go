package main

import (
	"fmt"
	"os"

	"github.com/gophergala/edrans-smartcity/algorithm"
)

func main() {
	var i int
	city := algorithm.GetTestCity()
	vehicle, e := city.CallService("doctor")
	if e != nil {
		fmt.Println(e)
		os.Exit(2)
	}
	paths, e := algorithm.GetPaths(&city, vehicle.Position.ID, 3)
	if e != nil {
		fmt.Printf("Ohh no... %+v\n", e)
		os.Exit(2)
	}
	path := algorithm.SortCandidates(algorithm.CalcEstimatesForVehicle(vehicle, paths))[0]
	vehicle.Alert <- path
	fmt.Scanf("%d", &i)
	/*for i := 0; i < 1; i++ {
	  for j := 0; j < len(paths[0].Links); j++ {
	    fmt.Printf("Link #%d: %+v\n", j, paths[0].Links[j])
	  }
	}*/
}
