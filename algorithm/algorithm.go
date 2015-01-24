package algorithm

import "fmt"

func (c City) GetPath(origin, destiny int) ([]Path, error) {
	if origin == destiny {
		return nil, fmt.Errorf("Already at destiny")
	}
	org := c.getNode(origin)
	dest := c.getNode(destiny)
	candidates := c.getCandidates(org, dest, nil)
	candidates = calcEstimates(candidates)
	candidates = orderCandidates(candidates)
	return candidates, c.err
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

func orderCandidates(paths []Path) []Path {
	var done bool
	var x int
	for i := 0; i < len(paths) && !done; i++ {
		done = true
		for j := 0; j < len(paths)-x-1; j++ {
			if paths[i].OriginalEstimate > paths[i+1].OriginalEstimate {
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
