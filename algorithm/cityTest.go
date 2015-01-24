package algorithm

import (
	"container/ring"
	"fmt"
	"time"
)

var (
	defaultSem      = Semaphore{Inputs: ring.New(0), ActiveInput: nil, Interval: defaultInterval, Status: make(chan SemRequest, 1), Paused: false}
	defaultInterval = 1 * time.Minute
)

func GetTestCity() City {
	var city = make([]Node, 17)
	city[1] = Node{ID: 1, Outputs: []Link{Link{Name: "Roca", OriginID: 1, DestinyID: 2, Weight: 30}, Link{Name: "Pellegrini", OriginID: 1, DestinyID: 5, Weight: 30}}}
	city[2] = Node{ID: 2, Outputs: []Link{Link{Name: "Roca", OriginID: 2, DestinyID: 3, Weight: 30}}}
	city[3] = Node{ID: 3, Outputs: []Link{Link{Name: "Roca", OriginID: 3, DestinyID: 4, Weight: 30}, Link{Name: "Irigoyen", OriginID: 3, DestinyID: 7, Weight: 35}}}
	city[4] = Node{ID: 4, Outputs: []Link{}}
	city[5] = Node{ID: 5, Outputs: []Link{Link{Name: "Pellegrini", OriginID: 5, DestinyID: 9, Weight: 30}}}
	city[6] = Node{ID: 6, Outputs: []Link{Link{Name: "Rivadavia", OriginID: 6, DestinyID: 5, Weight: 35}, Link{Name: "Irigoyen", OriginID: 6, DestinyID: 2, Weight: 35}}}
	city[7] = Node{ID: 7, Outputs: []Link{Link{Name: "Rivadavia", OriginID: 7, DestinyID: 6, Weight: 45}, Link{Name: "Palacios", OriginID: 7, DestinyID: 11, Weight: 45}}}
	city[8] = Node{ID: 8, Outputs: []Link{Link{Name: "Rivadavia", OriginID: 8, DestinyID: 7, Weight: 35}, Link{Name: "Justo", OriginID: 8, DestinyID: 12, Weight: 30}}}
	city[9] = Node{ID: 9, Outputs: []Link{Link{Name: "Mitre", OriginID: 9, DestinyID: 10, Weight: 35}, Link{Name: "Pellegrini", OriginID: 9, DestinyID: 13, Weight: 30}}}
	city[10] = Node{ID: 10, Outputs: []Link{Link{Name: "Irigoyen", OriginID: 10, DestinyID: 6, Weight: 45}, Link{Name: "Mitre", OriginID: 10, DestinyID: 11, Weight: 45}}}
	city[11] = Node{ID: 11, Outputs: []Link{Link{Name: "Palacios", OriginID: 11, DestinyID: 15, Weight: 35}, Link{Name: "Mitre", OriginID: 11, DestinyID: 12, Weight: 35}}}
	city[12] = Node{ID: 12, Outputs: []Link{Link{Name: "Justo", OriginID: 12, DestinyID: 8, Weight: 30}}}
	city[13] = Node{ID: 13, Outputs: []Link{}}
	city[14] = Node{ID: 14, Outputs: []Link{Link{Name: "Irigoyen", OriginID: 14, DestinyID: 10, Weight: 35}, Link{Name: "Urquiza", OriginID: 14, DestinyID: 13, Weight: 30}}}
	city[15] = Node{ID: 15, Outputs: []Link{Link{Name: "Urquiza", OriginID: 15, DestinyID: 14, Weight: 30}}}
	city[16] = Node{ID: 16, Outputs: []Link{Link{Name: "Justo", OriginID: 16, DestinyID: 12, Weight: 30}, Link{Name: "Urquiza", OriginID: 16, DestinyID: 15, Weight: 30}}}
	myCity := City{Nodes: city, Name: "Fake Buenos Aires"}
	myCity.GenerateSem()
	//myCity.EnableSem()
	return myCity
}

/*

d      u      d      u

1 ---- 2 ---- 3 ---- 4    r Roca
|      |      |      |
5 ---- 6 ---- 7 ---- 8    l Rivadavia
|      |      |      |
9 ---- a ---- b ---- c    r Mitre
|      |      |      |
d ---- e ---- f ---- g    l Urquiza

- Pellegrini
- Irigoyen
- Palacios
- Justo

*/

func (c City) GenerateSem() {
	for i := 1; i < len(c.Nodes); i++ {
		links := c.getLinked(i)
		if len(links) == 0 {
			c.Nodes[i].Sem = defaultSemaphore()
			continue
		}
		var sem Semaphore
		sem.Interval = defaultInterval
		sem.Inputs = ring.New(len(links))
		for j := 0; j < len(links); j++ {
			sem.Inputs = sem.Inputs.Next()
			sem.ActiveInput, _ = sem.Inputs.Next().Value.(*Link)
			sem.Inputs.Value = &links[j]
			sem.ActiveInput, _ = sem.Inputs.Next().Value.(*Link)
		}
		sem.Status = make(chan SemRequest, 500)
		go sem.Start()
	}
}

func defaultSemaphore() Semaphore {
	return Semaphore{Inputs: ring.New(0), ActiveInput: nil, Interval: defaultInterval, Status: make(chan SemRequest, 1), Paused: false}
}

func (c City) EnableSem() {
	fmt.Println("len:", len(c.Nodes))
	for i := 1; i < len(c.Nodes); i++ {
		fmt.Println("#85,", i)
		c.Nodes[i].Sem.Status <- SemRequest{Status: false}
	}
}

func (c City) getLinked(node int) []Link {
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

func (sem *Semaphore) Start() {
	change := time.After(sem.Interval)
	for {
		select {
		case <-change:
			if !sem.Paused {
				sem.ActiveInput = sem.Inputs.Next().Value.(*Link)
				sem.Inputs = sem.Inputs.Next()
			}
			change = time.After(sem.Interval)
		case req := <-sem.Status:
			sem.Paused = req.Status
			if req.Status {
				for i := 0; i < sem.Inputs.Len(); i++ {
					sem.Inputs = sem.Inputs.Next()
					if sem.Inputs.Next().Value.(*Link).Name == req.Allow {
						sem.ActiveInput = sem.Inputs.Next().Value.(*Link)
						sem.Inputs = sem.Inputs.Next()
					}
				}
			} else {
				change = time.After(1 * time.Second)
			}
		}
	}
}
