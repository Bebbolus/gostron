//
// Copyright © 2018 Roberto Della Fornace
// Implement a modular web server in Go
// License: MIT included in repository
//

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"plugin"
	"strconv"
	"time"
)

//source routes configuration struct to load from the json configuration file
type routes struct {
	Endpoints []struct {
		Controller  string `json:"controller"`
		Middlewares []struct {
			Handler string `json:"handler"`
			Params  string `json:"params"`
		} `json:"middlewares"`
		Path string `json:"path"`
	} `json:"endpoints"`
}

var RoutesConf routes

//source server configuration struct to load from json configuration file
type server struct {
	Listento     string `json:"listento"`
	Readtimeout  string `json:"readtimeout"`
	Writetimeout string `json:"writetimeout"`
}

var ServerConf server

/* PLUGINS */

//Controller is a local hanlder plugin interface
type Controller interface {
	Fire(w http.ResponseWriter, r *http.Request)
}

/* MIDDLEWARES */

//Middleware is local handler plugin interface, it will return a Gate compatible function
type Middleware interface {
	Pass(args string) func(http.HandlerFunc) http.HandlerFunc
}

//Gate is a type that describe the middleware functions that will be chained to the route
type Gate func(http.HandlerFunc) http.HandlerFunc

//Chain function concatenate the middlewares (typed as Gate function)
func Chain(f http.HandlerFunc, middlewares ...Gate) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}

/* HELPER FUNCTIONS */
//kill function print the message and then exit
func kill(msg interface{}) {
	fmt.Println(msg)
	os.Exit(1)
}

//must function check if there is an error
func must(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//ReadFromJSON function load a json file into a struct or return error
func ReadFromJSON(t interface{}, filename string) error {

	jsonFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(jsonFile), t)
	if err != nil {
		log.Fatalf("error: %v", err)
		return err
	}

	return nil
}

//start point
func main() {

	//load configurations from files
	must(ReadFromJSON(&ServerConf, "configurations/server.json"))
	must(ReadFromJSON(&RoutesConf, "configurations/routes.json"))

	//convert source json strings to integer where needed
	readtimeout, err := strconv.Atoi(ServerConf.Readtimeout)
	if err != nil {
		kill(err)
	}
	writetimeout, err := strconv.Atoi(ServerConf.Writetimeout)
	if err != nil {
		kill(err)
	}

	//set server configurations
	srv := &http.Server{
		ReadTimeout:  time.Duration(readtimeout) * time.Second,
		WriteTimeout: time.Duration(writetimeout) * time.Second,
		Addr:         ServerConf.Listento,
	}

	// based on the source confguration routes, loop on every configuration and load relative plugins
	// plugin.Open: If a path has already been opened, then the existing *Plugin is returned.
	// It is safe for concurrent use by multiple goroutines.
	for _, v := range RoutesConf.Endpoints {
		// load module:
		plug, err := plugin.Open(v.Controller)
		if err != nil {
			kill(err)
		}
		// look up for an exported Controller method
		symController, err := plug.Lookup("Controller")
		if err != nil {
			kill(err)
		}

		// check that loaded symbol is type Controller
		var controller Controller
		controller, ok := symController.(Controller)
		if !ok {
			kill("The Controller module have wrong type")
		}

		//define new middleware chain
		var chain []Gate

		// foreach middleware configured for the actual routepath
		for _, mid := range v.Middlewares {
			// load middleware plugin
			plug, midErr := plugin.Open(mid.Handler)
			if midErr != nil {
				kill(midErr)
			}
			// look up the Pass function
			symMiddleware, midErr := plug.Lookup("Middleware")
			if midErr != nil {
				kill(midErr)
			}

			// check that loaded symbol is type Middleware
			var middleware Middleware
			middleware, ok := symMiddleware.(Middleware)
			if !ok {
				kill("The middleware module have wrong type")
			}

			// build the gate function that contain the middleware instance
			nmid := Gate(middleware.Pass(mid.Params))

			// append to the middlewares chain
			chain = append(chain, nmid)

		}
		// Use all the modules to handle the request
		http.HandleFunc(v.Path, Chain(controller.Fire, chain...))

	}
	
	//SERVER START AND ERROR MANAGEMENT
	//best practise: start a local istance of server mux to avoid imported lib to define malicious handler
	log.Fatal(srv.ListenAndServe(), http.NewServeMux())
}
