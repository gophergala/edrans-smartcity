package models

import (
	"math/rand"
)

type PublicService struct {
	Location int //NodeID
	Service  string
	Vehicles []Vehicle
	Errors   chan error
}

//will be used only for patrolmen
func (s *PublicService) readErrors(c *City) {
	for {
		<-s.Errors
		newPatrolman := Vehicle{Service: "PoliceDept", MinWeight: 5, Alert: make(chan Path, 1), Errors: s.Errors, InCity: c}
		s.Vehicles = append(s.Vehicles, newPatrolman)
		go newPatrolman.patrol(rand.Int() % c.GetNumNodes())
	}
}
