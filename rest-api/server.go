package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	//"github.com/gophergala/edrans-smartcity/models"
	//"github.com/gophergala/edrans-smartcity/generators"
	"github.com/gorilla/mux"
)

var sessions map[ID]models.City

type ID string
type handler func(w http.ResponseWriter, r *http.Request, ctx *context) (int, interface{})

type context struct {
	Body   []byte
	CityID ID
}

func main() {
	var port int
	flag.IntVar(&port, "port", 2489, "port server will be launched")
	flag.Parse()

	sessions = make(map[ID]models.City)

	muxRouter := mux.NewRouter()
	muxRouter.StrictSlash(false)

	muxRouter.Handle("/city", handler(getCity)).Methods("GET")
	muxRouter.Handle("/emergency", handler(postEmergency)).Methods("POST")

	http.Handle("/", muxRouter)
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		fmt.Println("Cannot launch server:", err)
		os.Exit(2)
	}
	fmt.Printf("Listening on port %s...\n", port)
	http.Serve(listener, nil)

}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var ctx context
	var e error
	var ok bool

	ctx.Body, e = ioutil.ReadAll(r.Body)
	if e != nil {
	}

	ctx.CityID = ID(r.Header.Get("my-city"))
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

func getCity(w http.ResponseWriter, r *http.Request, ctx *context) (status int, response interface{}) {
	if ctx.CityID != "" {
		status = 403
		response = "You already have a city"
		return
	}
	city, e := generators.NewCity()
	response = city
	if e != nil {
		status = 400
		response = e
	}
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
	vehicle, e := sessions[ctx.CityID].CallService(emergency.Service)
	if e != nil {
		status = 400
		response = e
		return
	}
	paths := sessions[ctx.CityID].GetPaths(vehicle.Position.ID, emergency.Where)
	paths = vehicle.CalcPaths(paths)
	return
}
