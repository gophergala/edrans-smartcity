package algorithm

import (
	"fmt"
	"math/rand"
	"time"
)

type City struct {
	Nodes    []Node
	Services []PublicService
	Name     string
	err      error
}

type PublicService struct {
	Location int //NodeID
	Service  string
	Vehicles []Vehicle
	Errors   chan error
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
	Service      string
	MinWeight    int
	Busy         bool
	Alert        chan Path
	Errors       chan error
	InCity       *City
	Position     *Node
	BasePosition *Node
}

func (c *City) addService(service string, location, vehicles, minWeight int) {
	var newservice PublicService
	newservice.Service = service
	newservice.Location = location
	newservice.Errors = make(chan error, 5)
	for i := 0; i < vehicles; i++ {
		newservice.Vehicles = append(newservice.Vehicles, Vehicle{Service: service, MinWeight: minWeight, Errors: newservice.Errors, InCity: c, BasePosition: c.getNode(location), Position: c.getNode(location), Alert: make(chan Path, 5)})
	}
	c.Services = append(c.Services, newservice)
}

func (v *Vehicle) patrol(start int) {
	patrol := time.After(1 * time.Second)
	for {
		select {
		case <-patrol:
			node := v.InCity.getNode(start)
			if len(node.Outputs) == 0 {
				v.Errors <- fmt.Errorf("can not go on patrol")
				return
			}
			v.Position = node
			patrol = time.After(time.Duration(node.Outputs[0].Weight) * time.Second)
			start = node.Outputs[0].DestinyID
		case path := <-v.Alert:
			v.run(path)
			v.Position = v.InCity.getNode(path.Links[len(path.Links)-1].DestinyID)
			start = v.Position.ID
		}
	}
}

func (v *Vehicle) wait() {
	for {
		path := <-v.Alert
		fmt.Println("models:106")
		v.run(path)
		v.Position = v.BasePosition
	}
}

func (c *City) launchVehicles() {
	for i := 0; i < len(c.Services); i++ {
		if c.Services[i].Service == "hospital" || c.Services[i].Service == "firehouse" {
			for j := 0; j < len(c.Services[i].Vehicles); j++ {
				go c.Services[i].Vehicles[j].wait()
			}
		} else {
			go c.Services[i].readErrors(c)
			for j := 0; j < len(c.Services[i].Vehicles); j++ {
				go c.Services[i].Vehicles[j].patrol(rand.Int() % len(c.Nodes))
			}
		}
	}
}

//will be used only for patrolmen
func (s *PublicService) readErrors(c *City) {
	for {
		<-s.Errors
		newPatrolman := Vehicle{Service: "policeman", MinWeight: 5, Alert: make(chan Path, 1), Errors: s.Errors, InCity: c}
		s.Vehicles = append(s.Vehicles, newPatrolman)
		go newPatrolman.patrol(rand.Int() % len(c.Nodes))
	}
}

func (c *City) CallService(service string) (*Vehicle, error) {
	switch service {
	case "doctor":
		return c.callService("hospital", "ambulance")
	case "fireman":
		return c.callService("firehouse", "pumper")
	}
	return c.callService("policeman", "patrolman")
}

func (c *City) callService(service, name string) (*Vehicle, error) {
	var base PublicService
	for i := 0; i < len(c.Services); i++ {
		if c.Services[i].Service == service {
			base = c.Services[i]
		}
	}
	for i := 0; i < len(base.Vehicles); i++ {
		if !base.Vehicles[i].Busy {
			return &base.Vehicles[i], nil
		}
	}
	return nil, fmt.Errorf("There is no %s available", name)
}

func (v *Vehicle) run(path Path) time.Duration {
	fmt.Println("running")
	v.Busy = true
	now := time.Now()
	var i int
	for i = 0; i < len(path.Links); i++ {
		v.InCity.getNode(path.Links[i].DestinyID).Sem.Status <- SemRequest{Status: true, Allow: path.Links[i].Name}
		time.Sleep(time.Duration(path.Weights[i]) * time.Second)
		v.InCity.getNode(path.Links[i].DestinyID).Sem.Status <- SemRequest{Status: false, Allow: path.Links[i].Name}
		v.Position = v.InCity.getNode(path.Links[i].DestinyID)
	}
	v.Busy = false
	return time.Since(now)
}
