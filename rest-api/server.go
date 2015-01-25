package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gophergala/edrans-smartcity/algorithm"
	"github.com/gophergala/edrans-smartcity/factory"
	"github.com/gophergala/edrans-smartcity/models"
	"github.com/gorilla/mux"
)

var sessions map[int]*models.City

type handler func(w http.ResponseWriter, r *http.Request, ctx *context) (int, interface{})

type context struct {
	Body   []byte
	CityID int
}

func main() {
	var port int
	flag.IntVar(&port, "port", 2489, "port server will be launched")
	flag.Parse()

	sessions = make(map[int]*models.City)
	sessions[0], _ = factory.CreateRectangularCity(3, 3, "default")

	muxRouter := mux.NewRouter()
	muxRouter.StrictSlash(false)

	muxRouter.Handle("/city", handler(getCity)).Methods("GET")
	//muxRouter.Handle("/sample-city", handler(postSampleCity)).Methods("POST")
	muxRouter.Handle("/sample-city", handler(postSampleCity)).Methods("POST")
	muxRouter.Handle("/emergency", handler(postEmergency)).Methods("POST")
	muxRouter.Handle("/city/{cityID}", handler(getIndex)).Methods("GET")
	muxRouter.HandleFunc("/city/img/0.jpg", handleFile("img/0.jpg"))
	muxRouter.HandleFunc("/city/img/1.jpg", handleFile("img/1.jpg"))
	muxRouter.HandleFunc("/city/img/2.jpg", handleFile("img/2.jpg"))
	muxRouter.HandleFunc("/city/img/3.jpg", handleFile("img/3.jpg"))
	muxRouter.HandleFunc("/city/img/-1.jpg", handleFile("img/-1.jpg"))

	http.Handle("/", muxRouter)
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		fmt.Println("Cannot launch server:", err)
		os.Exit(2)
	}
	fmt.Printf("Listening on port %d...\n", port)
	http.Serve(listener, nil)
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var ctx context
	var e error

	ctx.Body, e = ioutil.ReadAll(r.Body)
	if e != nil {
		fmt.Println("Error when reading body!")
	}

	//ctx.CityID = r.Header.Get("city-name")
	vars := mux.Vars(r)
	value, _ := vars["cityID"]
	ctx.CityID, _ = strconv.Atoi(value)

	status, response := h(w, r, &ctx)
	if status == -1 {
		return
	}
	if status == 0 {
		status = 200
	}
	if response == nil {
		response = map[string]string{"status": "ok"}
	}
	if status < 200 || status >= 300 {
		response = map[string]interface{}{"error": response}
	}
	responseJSON, _ := json.Marshal(response)
	w.WriteHeader(status)
	w.Write(responseJSON)
}

func postSampleCity(w http.ResponseWriter, r *http.Request, ctx *context) (status int, response interface{}) {
	type cityParams struct {
		SizeHorizontal int    `json:"size-horizontal"`
		SizeVertical   int    `json:"size-vertical"`
		Name           string `json:"name"`
	}
	type cityOut struct {
		CityName string `json:"city-name"`
	}
	var in cityParams
	status = 302
	var url string
	err := json.Unmarshal(ctx.Body, &in)
	if err != nil {
		status = 400
	}
	/*if err != nil {
		status = http.StatusBadRequest
		fmt.Printf("error in body %+v\n", string(ctx.Body))
		response = "invalid json body"
		return
	}*/

	cityID := len(sessions) + 1
	sessions[cityID], err = factory.CreateRectangularCity(in.SizeHorizontal, in.SizeVertical, in.Name)
	if err != nil {
		status = 400
		//response = fmt.Sprintf("Error: %s", err)
		//return
	}

	/*response = cityOut{
		CityName: in.Name,
	}*/
	if status != 302 {
		url = "/error"
	} else {
		url = fmt.Sprintf("/%d", cityID)
	}
	http.Redirect(w, r, url, status)
	return -1, nil
}

func getCity(w http.ResponseWriter, r *http.Request, ctx *context) (status int, response interface{}) {
	if ctx.CityID != 0 {
		status = 403
		response = "You already have a city"
		return
	}

	city := factory.SampleCity() // TODO MUST REPLACE THIS
	response = city
	/*if e != nil {
		status = 400
		response = e
	}*/
	return
}

type emergencyRequest struct {
	Service string `json:"service"`
	Where   int    `json:"location"`
}

func postEmergency(w http.ResponseWriter, r *http.Request, ctx *context) (status int, response interface{}) {
	var emergency emergencyRequest
	e := json.Unmarshal(ctx.Body, &emergency)
	if e != nil {
		status = 400
		response = e
		return
	}
	city := sessions[ctx.CityID]
	vehicle, e := city.CallService(emergency.Service)
	if e != nil {
		status = 400
		response = e
		return
	}
	paths, e := algorithm.GetPaths(city, vehicle.Position.ID, emergency.Where)
	if e != nil {
		status = 400
		response = e
		return
	}
	paths = algorithm.CalcEstimatesForVehicle(vehicle, paths)
	vehicle.Alert <- algorithm.SortCandidates(paths)[0]
	response = fmt.Sprintf("%s on the way to %d", emergency.Service, emergency.Where)
	return
}

func getIndex(w http.ResponseWriter, r *http.Request, ctx *context) (status int, response interface{}) {
	var index = make([]string, 0)
	file, e := ioutil.ReadFile("index.html")
	if e != nil {
		return 503, "index not found"
	}
	fileLines := strings.Split(string(file), "\n")
	for i := 0; i < len(fileLines); i++ {
		index = append(index, fileLines[i])
		if strings.Contains(fileLines[i], "<table") {
			table := createTable(ctx.CityID)
			for j := 0; j < len(table); j++ {
				index = append(index, table[j])
			}
		}
	}
	http.ServeContent(w, r, "city", time.Now(), bytes.NewReader([]byte(strings.Join(index, "\n"))))
	status = -1
	response = strings.Join(index, "\n")
	return
}

func createTable(cityID int) []string {
	var table = make([]string, 0)
	city := sessions[cityID]
	locations := city.GetLocations()
	nodesRoot := int(math.Sqrt(float64(len(locations))))
	for i := 0; i < nodesRoot; i++ {
		table = append(table, "<tr>")
		myNodes := getNodes(locations, i)
		for j := 0; j < len(myNodes); j++ {
			var color string
			switch myNodes[j].Vehicle {
			case 0:
				color = fmt.Sprintf(`bgcolor="#0000FF"`)
			case 1:
				color = fmt.Sprintf(`bgcolor="#FF0000"`)
			case 2:
				color = fmt.Sprintf(`bgcolor="#228B22"`)
			}
			insert := fmt.Sprintf(`<img src="img/%d.jpg" height="20" width="20" />`, myNodes[j].Input)
			table = append(table, fmt.Sprintf(`<td style="width:100px" %s> %s </td>`, color, insert))
		}
		table = append(table, "</tr>")
	}
	return table
}

func getNodes(locations []models.Location, nodes int) []models.Location {
	var local = make([]models.Location, 0)
	for i := 0; i < len(locations); i++ {
		if locations[i].Lat == nodes {
			local = append(local, locations[i])
		}
	}
	/*var done bool
	for i := 0; i < len(local)-1 && !done; i++ {
		done = true
		if local[i].Long > local[i+1].Long {
			aux := local[i]
			local[i] = local[i+1]
			local[i+1] = aux
			done = false
		}
	}*/
	return local
}

func handleFile(path string) http.HandlerFunc {
	path = filepath.Join("", path)
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}
