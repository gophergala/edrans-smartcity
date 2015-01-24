package algorithm

import (
	"container/ring"
	"time"
)

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
	Inputs      *ring.Ring
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
	Estimate         int
	OriginalEstimate int
	Reached          bool
	ForgetMe         bool
}
