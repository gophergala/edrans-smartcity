package models

import (
	"math/rand"
)

const (
	SERVICE_HOSPITAL    = "Hospital"
	SERVICE_POLICE      = "PoliceDept"
	SERVICE_FIREFIGHTER = "FireDept"

	VEHICLE_AMBULANCE  = "ambulance"
	VEHICLE_POLICE_CAR = "police car"
	VEHICLE_PUMPER     = "pumper"

	CALL_SERVICE_MEDIC   = "Medic"
	CALL_SERVICE_POLICE  = "Police"
	CALL_SERVICE_FIREMAN = "Fireman"
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
		newPatrolman := Vehicle{Service: SERVICE_POLICE, MinWeight: 5, Alert: make(chan Path, 1), Errors: s.Errors, InCity: c}
		s.Vehicles = append(s.Vehicles, newPatrolman)
		go newPatrolman.patrol(rand.Int() % c.GetNumNodes())
	}
}
