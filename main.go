package main

import (
        "log"
        "net/http"
        "fmt"
        "encoding/json"
        "os"
        "time"
        "plugin"
        bootstrap "github.com/Bebbolus/gostron/bootstrap"
)

//source configuration struct to map the json configuration file
type sourceConfig struct {
	Endpoints []struct {
		Handler string `json:"handler"`
        Middlewares []struct {
			Handler string `json:"handler"`
			Params  string `json:"params"`
		} `json:"middlewares"`
		Path    string `json:"path"`
	} `json:"endpoints"`
	Server struct {
		Listento     string `json:"listento"`
		Readtimeout  time.Duration `json:"readtimeout"`
		Writetimeout time.Duration `json:"writetimeout"`
	} `json:"server"`
}

//read source configuration file and map into local struct
func ReadConfiguration(confFile string) (sourceConfig, error){
    Conf := sourceConfig{}
    file, err := os.Open(confFile)
    defer file.Close()
    if err != nil {
        return Conf, err
    }
    decoder := json.NewDecoder(file)
    decodingErr := decoder.Decode(&Conf)
    if decodingErr != nil {
        return Conf, decodingErr
    }
    return Conf, nil
}

/* PLUGINS */

//local Http hanlder plugin interface
type Handler interface{
        Fire(w http.ResponseWriter, r *http.Request)
}

/* MIDDLEWARES */

//local Middlewares handler plugin interface
type Middleware interface{
        Pass()
}

//start point
func main() {
        configuration, _ := ReadConfiguration("configuration.json")

        //SET UP SERVER TIMEOUT
        srv := &http.Server{
            ReadTimeout: configuration.Server.Readtimeout * time.Second,
            WriteTimeout: configuration.Server.Writetimeout * time.Second,
            Addr:configuration.Server.Listento,
        }

        for _,v := range configuration.Endpoints{
            // load module
            // 1. open the so file to load the symbols
            plug, err := plugin.Open(v.Handler)
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }

            // 2. look up a symbol (an exported function or variable)
            // in this case, variable Controller
            symController, err := plug.Lookup("Handler")
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }

            // 3. Assert that loaded symbol is of a desired type
            // in this case interface type Handler (defined above)
            var handler Handler
            handler, ok := symController.(Handler)
            if !ok {
                fmt.Println("unexpected type from module symbol")
                os.Exit(1)
            }

            var chain []bootstrap.Gate

            /*
            per ogni middleware configurato da eseguire su questo path:
                carica il plugin di quel middleware e lo carica
                cerca la variabile per i parametri e gli assegna il valore come da configurazione
                lo aggiunge alla catena
            */

            for _,mid := range v.Middlewares {
                // load module
                // 1. open the so file to load the symbols
                plug, midErr := plugin.Open(mid.Handler)
                if midErr != nil {
                    fmt.Println(midErr)
                    os.Exit(1)
                }
                // 2. look up a symbol (an exported function or variable)
                // in this case, function Pass()
                symFunc, midErr := plug.Lookup("Pass")
                if midErr != nil {
                    fmt.Println(midErr)
                    os.Exit(1)
                }

                //chain = append(chain, symFunc.(func() bootstrap.Gate))
                chain = append(chain, symFunc.(func(string) bootstrap.Gate)(mid.Params) )

            }
            /*
                fine
            */

            // 4. use the module to handle the request
            http.HandleFunc(v.Path, bootstrap.Chain(handler.Fire, chain...))



        }
        //best practise: start a local istance of server mux to avoid imported lib to define malicious handler
        mux := http.NewServeMux()

        //SERVER START AND ERROR MANAGEMENT
        log.Fatal(srv.ListenAndServe(),mux)
}
