package algorithm

import "time"

type City struct {
	Nodes []Node
	Name  string
	err   error
}

type Node struct {
	ID      int
	Outputs []Link
	Sem     Semaphore
	//Latitude int
	//Longitude int
}

type Link struct {
	Name      string
	OriginID  int
	DestinyID int
	Weight    int
}

type Semaphore struct {
	Inputs      []Link
	ActiveInput *Link
	Interval    time.Duration
	Status      chan SemRequest
	Paused      bool
}

type SemRequest struct {
	Status bool
	Allow  string //link active's street
}

type Path struct {
	Links            []Link
	Weights          []int
	Estimate         int
	OriginalEstimate int
	Reached          bool
	ForgetMe         bool
}

type Vehicle struct {
	Service   string
	MinWeight int
	InCity    *City
	Position  *Node
}

func CallService(service string) Vehicle {
	switch service {
	case "doctor":
		return callDoctors()
	case "fireman":
		return callFiremen()
	}
	return callPolicemen()
}

func callDoctors() Vehicle {
	return Vehicle{Service: "Ambulance", MinWeight: 10}
}

func callFiremen() Vehicle {
	return Vehicle{Service: "Pumper", MinWeight: 15}
}

func callPolicemen() Vehicle {
	return Vehicle{Service: "Patrolman", MinWeight: 5}
}
