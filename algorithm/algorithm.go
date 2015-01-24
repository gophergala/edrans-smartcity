package algorithm

import "fmt"

func (c City) GetPath(origin, destiny int) (*Path, error) {
	if origin == destiny {
		return nil, fmt.Errorf("Already at destiny")
	}
	org := c.getNode(origin)
	dest := c.getNode(destiny)
	candidates := c.getCandidates(org, dest)
	return &Path{}, c.err
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

func (c *City) getCandidates(org, dst *Nodes, acWeight int) []Path {
	if c.Err != nil {
		return nil
	}
	var paths = make([]Path, 0)
	for i := 0; i < len(org.Outputs); i++ {
		if org.Outputs[i].DestinyID == dst.ID {
			paths = append(paths, Path{Reached: true})
		} else {
			subPaths := getCandidates(c.Nodes[org.Outputs[i].DestinyID], dst, c.Nodes[org.Outputs[i].Weight]+acWeight)
		}
	}
	return paths
}
