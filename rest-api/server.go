package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gophergala/edrans-smartcity/algorithm"
	"github.com/gophergala/edrans-smartcity/factory"
	"github.com/gophergala/edrans-smartcity/models"
	"github.com/gorilla/mux"
)

var sessions map[string]*models.City

type handler func(w http.ResponseWriter, r *http.Request, ctx *context) (int, interface{})

type context struct {
	Body   []byte
	CityID string
}

func main() {
	var port int
	flag.IntVar(&port, "port", 2489, "port server will be launched")
	flag.Parse()

	sessions = make(map[string]*models.City)

	muxRouter := mux.NewRouter()
	muxRouter.StrictSlash(false)

	muxRouter.Handle("/city", handler(getCity)).Methods("GET")
	muxRouter.Handle("/emergency", handler(postEmergency)).Methods("POST")
	muxRouter.Handle("/sample-city", handler(postSampleCity)).Methods("POST")
	muxRouter.HandleFunc("/index", handleFile("index.html"))

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

	ctx.CityID = r.Header.Get("city-name")
	status, response := h(w, r, &ctx)
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
	err := json.Unmarshal(ctx.Body, &in)
	if err != nil {
		status = http.StatusBadRequest
		fmt.Printf("error in body %+v\n", string(ctx.Body))
		response = "invalid json body"
		return
	}

	sessions[in.Name], err = factory.CreateRectangularCity(in.SizeHorizontal, in.SizeVertical, in.Name)
	if err != nil {
		status = 400
		response = fmt.Sprintf("Error: %s", err)
		return
	}

	response = cityOut{
		CityName: in.Name,
	}

	return
}

func getCity(w http.ResponseWriter, r *http.Request, ctx *context) (status int, response interface{}) {
	if ctx.CityID != "" {
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

func handleFile(path string) http.HandlerFunc {
	path = filepath.Join("", path)
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}
