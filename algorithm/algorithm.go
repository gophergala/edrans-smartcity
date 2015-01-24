package algorithm

import "fmt"

/*func (c City) GetPath(origin, destiny int, service string) (*Path, error) {
	if origin == destiny {
		return nil, fmt.Errorf("Already at destiny")
	}
	org := c.getNode(origin)
	dest := c.getNode(destiny)
	candidates := c.getCandidates(org, dest, nil)
	vehicle := CallService(service)
	if c.err != nil || len(candidates) == 0 {
		return nil, c.err
	}
	candidates = calcEstimates(candidates)
	vehicle.InCity = &c
	candidates = vehicle.CalcPaths(candidates)
	candidates = OrderCandidates(candidates)
	return &candidates[0], c.err
}*/

func (c City) GetPaths(origin, destiny int) ([]Path, error) {
	if origin == destiny {
		return nil, fmt.Errorf("Already at destiny")
	}
	org := c.getNode(origin)
	dest := c.getNode(destiny)
	candidates := c.getCandidates(org, dest, nil)
	return orderLinks(candidates), c.err
}

func orderLinks(paths []Path) []Path {
	for i := 0; i < len(paths); i++ {
		var lnk = make([]Link, len(paths[i].Links))
		for j := 0; j < len(lnk); j++ {
			lnk[j] = paths[i].Links[len(lnk)-1-j]
		}
		paths[i].Links = lnk
	}
	return paths
}

func (c City) getNode(ID int) *Node {
	if c.err != nil {
		return nil
	}
	if len(c.Nodes) < ID {
		c.err = fmt.Errorf("Node %d does not exist", ID)
		return nil
	}
	return &c.Nodes[ID]
}

func (c *City) getCandidates(org, dst *Node, visited []int) []Path {
	if c.err != nil {
		return nil
	}
	vlen := len(visited)
	var paths = make([]Path, 0)
	for i := 0; i < len(org.Outputs); i++ {
		if org.Outputs[i].DestinyID == dst.ID {
			paths = append(paths, Path{Reached: true, Links: []Link{org.Outputs[i]}})
			continue
		}
		if alreadyVisited(org.Outputs[i].DestinyID, visited) {
			continue
		}
		visited = append(visited, org.Outputs[i].DestinyID)
		subPaths := c.getCandidates(c.getNode(org.Outputs[i].DestinyID), dst, visited)
		for j := 0; j < len(subPaths); j++ {
			lnks := subPaths[j].Links
			lnks = append(lnks, org.Outputs[i])
			paths = append(paths, Path{Links: lnks, Reached: subPaths[j].Reached})
		}
	}
	if vlen == 0 && len(paths) == 0 {
		c.err = fmt.Errorf("There's no way to the requested address")
	}
	return paths
}

func alreadyVisited(ID int, visited []int) bool {
	for i := 0; i < len(visited); i++ {
		if visited[i] == ID {
			return true
		}
	}
	return false
}

func calcEstimates(paths []Path) []Path {
	for i := 0; i < len(paths); i++ {
		for j := 0; j < len(paths[i].Links); j++ {
			paths[i].OriginalEstimate += paths[i].Links[j].Weight
		}
	}
	return paths
}

func OrderCandidates(paths []Path) []Path {
	var done bool
	var x int
	for i := 0; i < len(paths) && !done; i++ {
		done = true
		for j := 0; j < len(paths)-x-1; j++ {
			if paths[i].Estimate > paths[i+1].Estimate {
				aux := paths[i]
				paths[i] = paths[i+1]
				paths[i+1] = aux
				done = false
				x++
			}
		}
	}
	return paths
}

func (v *Vehicle) CalcPaths(paths []Path) []Path {
	for i := 0; i < len(paths); i++ {
		if len(paths[i].Links) == 0 {
			continue
		}
		paths[i].Weights = make([]int, 0)
		lastLink := paths[i].Links[0]
		weight := paths[i].Links[0].Weight
		paths[i].Weights = append(paths[i].Weights, paths[i].Links[0].Weight)
	LinksLoop:
		for j := 1; j < len(paths[i].Links); j++ {
			if paths[i].Links[j].Name == lastLink.Name {
				paths[i].Weights = append(paths[i].Weights, paths[i].Links[j].Weight)
				weight += paths[i].Links[j].Weight
				lastLink = paths[i].Links[j]
				continue LinksLoop
			}
			newLinkWeight := paths[i].Links[j].Weight - lastLink.Weight
			if newLinkWeight < v.MinWeight {
				newLinkWeight = v.MinWeight
			}
			lastLink = paths[i].Links[j]
			weight += newLinkWeight
			paths[i].Weights = append(paths[i].Weights, newLinkWeight)
		}
		paths[i].Estimate = weight
	}
	return paths
}