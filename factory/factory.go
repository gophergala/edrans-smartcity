package factory

import "github.com/gophergala/edrans-smartcity/models"

func CreateCity(numNodes int, name string) *models.City {
	var city = make([]models.Node, numNodes)

	city[1] = models.Node{ID: 1, Outputs: []models.Link{models.Link{Name: "Roca", OriginID: 1, DestinyID: 2, Weight: 30}, models.Link{Name: "Pellegrini", OriginID: 1, DestinyID: 5, Weight: 30}}}
	city[2] = models.Node{ID: 2, Outputs: []models.Link{models.Link{Name: "Roca", OriginID: 2, DestinyID: 3, Weight: 30}}}
	city[3] = models.Node{ID: 3, Outputs: []models.Link{models.Link{Name: "Roca", OriginID: 3, DestinyID: 4, Weight: 30}, models.Link{Name: "Irigoyen", OriginID: 3, DestinyID: 7, Weight: 35}}}
	city[4] = models.Node{ID: 4, Outputs: []models.Link{}}
	city[5] = models.Node{ID: 5, Outputs: []models.Link{models.Link{Name: "Pellegrini", OriginID: 5, DestinyID: 9, Weight: 30}}}
	city[6] = models.Node{ID: 6, Outputs: []models.Link{models.Link{Name: "Rivadavia", OriginID: 6, DestinyID: 5, Weight: 35}, models.Link{Name: "Irigoyen", OriginID: 6, DestinyID: 2, Weight: 35}}}
	city[7] = models.Node{ID: 7, Outputs: []models.Link{models.Link{Name: "Rivadavia", OriginID: 7, DestinyID: 6, Weight: 45}, models.Link{Name: "Palacios", OriginID: 7, DestinyID: 11, Weight: 45}}}
	city[8] = models.Node{ID: 8, Outputs: []models.Link{models.Link{Name: "Rivadavia", OriginID: 8, DestinyID: 7, Weight: 35}, models.Link{Name: "Justo", OriginID: 8, DestinyID: 12, Weight: 30}}}
	city[9] = models.Node{ID: 9, Outputs: []models.Link{models.Link{Name: "Mitre", OriginID: 9, DestinyID: 10, Weight: 35}, models.Link{Name: "Pellegrini", OriginID: 9, DestinyID: 13, Weight: 30}}}
	city[10] = models.Node{ID: 10, Outputs: []models.Link{models.Link{Name: "Irigoyen", OriginID: 10, DestinyID: 6, Weight: 45}, models.Link{Name: "Mitre", OriginID: 10, DestinyID: 11, Weight: 45}}}
	city[11] = models.Node{ID: 11, Outputs: []models.Link{models.Link{Name: "Palacios", OriginID: 11, DestinyID: 15, Weight: 35}, models.Link{Name: "Mitre", OriginID: 11, DestinyID: 12, Weight: 35}}}
	city[12] = models.Node{ID: 12, Outputs: []models.Link{models.Link{Name: "Justo", OriginID: 12, DestinyID: 8, Weight: 30}}}
	city[13] = models.Node{ID: 13, Outputs: []models.Link{}}
	city[14] = models.Node{ID: 14, Outputs: []models.Link{models.Link{Name: "Irigoyen", OriginID: 14, DestinyID: 10, Weight: 35}, models.Link{Name: "Urquiza", OriginID: 14, DestinyID: 13, Weight: 30}}}
	city[15] = models.Node{ID: 15, Outputs: []models.Link{models.Link{Name: "Urquiza", OriginID: 15, DestinyID: 14, Weight: 30}}}
	city[16] = models.Node{ID: 16, Outputs: []models.Link{models.Link{Name: "Justo", OriginID: 16, DestinyID: 12, Weight: 30}, models.Link{Name: "Urquiza", OriginID: 16, DestinyID: 15, Weight: 30}}}
	myCity := models.NewCity(city, name)
	myCity.AddService("hospital", 10, 5, 10)
	myCity.AddService("firehouse", 11, 5, 15)
	myCity.AddService("policeman", 16, 5, 5)
	myCity.LaunchVehicles()
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
