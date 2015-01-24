package models

import (
	"fmt"
	"math/rand"
)

type City struct {
	Nodes    []Node
	Services []PublicService
	Name     string

	// Aux vars
	LastError error
}

type Node struct {
	ID       int
	Location []int
	Outputs  []Link
	Sem      Semaphore
	//Latitude int
	//Longitude int
}

type Link struct {
	Name      string
	OriginID  int
	DestinyID int
	Weight    int
}

type Path struct {
	Links            []Link
	Weights          []int
	Estimate         int
	OriginalEstimate int
	Reached          bool
	ForgetMe         bool
}

func (c *City) addService(service string, location, vehicles, minWeight int) {
	var newservice PublicService
	newservice.Service = service
	newservice.Location = location
	newservice.Errors = make(chan error, 5)
	for i := 0; i < vehicles; i++ {
		newservice.Vehicles = append(
			newservice.Vehicles,
			Vehicle{
				Service:      service,
				MinWeight:    minWeight,
				Errors:       newservice.Errors,
				InCity:       c,
				BasePosition: c.GetNode(location),
				Position:     c.GetNode(location),
				Alert:        make(chan Path, 5),
			})
	}
	c.Services = append(c.Services, newservice)
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

func (c *City) EnableSem() {
	for i := 1; i < len(c.Nodes); i++ {
		c.Nodes[i].Sem.Status <- SemRequest{Status: false}
	}
}

func (c *City) getLinked(node int) []Link {
	var links []Link
	for i := 1; i < len(c.Nodes); i++ {
		for j := 0; j < len(c.Nodes[i].Outputs); j++ {
			if c.Nodes[i].Outputs[j].DestinyID == node {
				links = append(links, c.Nodes[i].Outputs[j])
			}
		}
	}
	return links
}

func (c *City) GenerateSem() {
	for i := 1; i < len(c.Nodes); i++ {
		links := c.getLinked(i)
		if len(links) == 0 {
			c.Nodes[i].Sem = defaultSemaphore()
			continue
		}
		var sem Semaphore
		sem.Interval = defaultInterval
		sem.Inputs = links
		sem.Status = make(chan SemRequest, 1)
		c.Nodes[i].Sem = sem
		go sem.Start()
	}
}

func (c *City) GetNode(ID int) *Node {
	if c.LastError != nil {
		return nil
	}
	if len(c.Nodes) < ID {
		c.LastError = fmt.Errorf("Node %d does not exist", ID)
		return nil
	}
	return &c.Nodes[ID]
}
