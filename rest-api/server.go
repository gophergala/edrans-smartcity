package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	//"github.com/gophergala/edrans-smartcity/models"
	"github.com/gorilla/mux"
)

var sessions map[ID]models.City

type ID string
type handler func(w http.ResponseWriter, r *http.Request, ctx *context)

type context struct {
	Body   string
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
	fmt.Printf("Listening on port %s...\n", conf.ApiPort)
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

	ctx.CityID = r.Header.Get("my-city")
	status, response := h(w, r, ctx)
	if status == 0 {
		status = 200
	}
	if response == nil {
		response = map[string]string{"status": "ok"}
	}
	responseJSON, _ := json.Marhsal(response)
	w.WriteHeader(status)
	w.Write(responseJSON)
}

func getCity(w http.ResponseWriter, r *http.Request, ctx *context) (status int, response interface{}) {
}

type emergencyRequest struct {
	Service string `json:"service"`
	Where   int    `json:"location"`
}

func postEmergency(w http.ResponseWriter, r *http.Request, ctx *context) (status int, response interface{}) {
}
