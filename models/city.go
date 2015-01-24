package models

import (
	"fmt"
	"math/rand"
)

type City struct {
	nodes    []Node
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

func NewCity(nodeList []Node, name string) (city *City) {
	myCity := City{nodes: nodeList, Name: name}
	myCity.generateSem()
	myCity.enableSem()

	return &myCity
}

func (c *City) GetNumNodes() int {
	return len(c.nodes)
}

func (c *City) AddService(service string, location, vehicles, minWeight int) {
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

func (c *City) LaunchVehicles() {
	for i := 0; i < len(c.Services); i++ {
		if c.Services[i].Service == "hospital" || c.Services[i].Service == "firehouse" {
			for j := 0; j < len(c.Services[i].Vehicles); j++ {
				go c.Services[i].Vehicles[j].wait()
			}
		} else {
			go c.Services[i].readErrors(c)
			for j := 0; j < len(c.Services[i].Vehicles); j++ {
				go c.Services[i].Vehicles[j].patrol(rand.Int() % len(c.nodes))
			}
		}
	}
}

func (c *City) CallService(call string) (*Vehicle, error) {
	switch call {
	case "Medic":
		return c.callService("Hospital", "ambulance")
	case "Fireman":
		return c.callService("FireDept", "pumper")
	case "Police":
		return c.callService("PoliceDept", "patrolman")
	}
	return nil, fmt.Errorf("unknown service")
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

func (c *City) enableSem() {
	for i := 0; i < len(c.nodes); i++ {
		c.nodes[i].Sem.Status <- SemRequest{Status: false}
	}
}

func (c *City) getLinked(node int) []Link {
	var links []Link
	for i := 0; i < len(c.nodes); i++ {
		for j := 0; j < len(c.nodes[i].Outputs); j++ {
			if c.nodes[i].Outputs[j].DestinyID == node {
				links = append(links, c.nodes[i].Outputs[j])
			}
		}
	}
	return links
}

func (c *City) generateSem() {
	for i := 0; i < len(c.nodes); i++ {
		links := c.getLinked(c.nodes[i].ID)
		if len(links) == 0 {
			c.nodes[i].Sem = defaultSemaphore()
			continue
		}
		var sem Semaphore
		sem.Interval = defaultInterval
		sem.Inputs = links
		sem.Status = make(chan SemRequest, 1)
		c.nodes[i].Sem = sem
		go sem.Start()
	}
}

func (c *City) GetNode(ID int) *Node {
	if c.LastError != nil {
		return nil
	}
	if len(c.nodes) < ID || ID <= 0 {
		c.LastError = fmt.Errorf("Node %d does not exist", ID)
		return nil
	}
	return &c.nodes[ID-1]
}
