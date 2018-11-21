### A modular web server in Go

![](https://cdn-images-1.medium.com/max/750/1*NZRlMTB55IffYytMDujUNA.png)
<span class="figcaption_hack">Architecture Representation</span>

The target of this article is to inspect how to build a modular web server. It
will send the requests through our pluggable modules: *middlewares *and
*controller*.

<br> 

how would we like to use it?

> *Our scope is to care only about the mapping of a request to a handlers, adding
> some middlewares from a set.*

![](https://cdn-images-1.medium.com/max/1000/1*AP7etBWgIvojKOBmyekj1Q.png)
<span class="figcaption_hack">Mapping middlewares and controller to an endpoint</span>

This kind of architecture is scalable, customizable and reusable. It enables us
to make:

1.  specialized web services
1.  use the same web services for several projects
1.  **Update separately **the “core” server part and the plugins.

I put particular attention to “**update separately**”. With plugin architecture,
you can distribute the compiled file of only one component. For example: if you
update the core to HTTPS architecture, you can redeploy only the core file. In
the same way, if you update the JWT plugin to use a new ash method, you have
only to redeploy the plugin.

![](https://cdn-images-1.medium.com/max/1000/1*XyZklA8d9HCg5YNrOzuxSg.png)
<span class="figcaption_hack">Distributed deployment example</span>

#### Implementation

We can start building our routes configuration file. An example configuration
can be a single route managed. We attach a plugin to check if the request HTTP
Method is GET or POST and then send it to a controller.

All other routes will return “404 not found”.

`routes.json` will look like this:

    {
       "endpoints":[
          {
             "path":"/myroute",
             "handler":"./plugins/controllers/general.so",
             "middlewares":[
                {
                   "handler":"./plugins/middlewares/method.so",
                   "params":"GET|POST"
                }
             ]
          }
       ]
    }

Creating the file in this way, we can attach several middlewares to a route and
use a middleware in several routes.

### Build the core

#### Make the middleware chain architecture

Our middleware concept will chain a set of functions. This functions will check
the request and if it passes the filter, send it to the next function.

Gate is the type that represents the middleware function with arguments valued:

    type Gate func(http.HandlerFunc) http.HandlerFunc
    func Chain(f http.HandlerFunc, middlewares ...Gate) http.HandlerFunc {
    	for _, m := range middlewares {
    		f = m(f)
    	}
    	return f
    }

#### Read the configurations

Now we can proceed on reading the configuration, mapping it to a struct (with
[this ](https://mholt.github.io/json-to-go/)tool is very simple):

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

and make the function to read from JSON:

    /ReadFromJSON function load a json file into a struct or return error
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

#### Load the Plugins

We can load plugins, using the plugin package. we can import all the exposed
functions and variables
([ELF](https://en.wikipedia.org/wiki/Executable_and_Linkable_Format) symbols).

As we call an exported type method from the plugin, we need to adopt some
conventions, I opted for:

* Controller type with method Fire()
* Middleware type with method Pass()

Walking into the configuration we can **dynamically **link the libraries:

> From “plugin.Open” documentation: If a path has already been opened, then the
> existing *Plugin is returned It is safe for concurrent use by multiple
goroutines.

Load Controller plugin:

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

Load middleware modules to attach on the route:

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

### Plugins Implementation

The package of a plugin needs to be “Main”.

> Unlike that, the package can’t see the entities such as types and functions in
> the “real” main package. So, as a suggestion, maintain plugins *dumber as
possible*.

In our repository create a plugin folder:

    mkdir plugins

Inside we create two folders, one for middlewares, one for controllers

    cd plugins
    mkdir controller
    mkdir middlewares

#### Build the Controllers

Inside the plugins/controllers folder create `general.so`, this will be the HTTP
Request handler:

    package main 
    import ( 
           "fmt" 
           "net/http"
    )

    type controller string

    func (h controller) Fire(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello FROM CONTROLLER PLUGIN!!!") 
    }

    // Controller exported namevar 
    Controller controller

#### Build the Middlewares

We build a method middleware that checks the HTTP Method, else returns a 400 Bad
Request.

To leave middleware “open”, it needs some arguments. In this case, a sequence of
approved HTTP methods that we need to split and check:

    package main

    import (
    	"net/http"
    	"strings"
    )

    type middleware string

    func (m middleware) Pass(args string) func(http.HandlerFunc) http.HandlerFunc {
      return func(f http.HandlerFunc) http.HandlerFunc {
        // Define the http.HandlerFunc
        return func(w http.ResponseWriter, r *http.Request) {
    	//split args and check if the request as this method
    	acceptedMethods := strings.Split(args, "|")
    	for _, v := range acceptedMethods {
    		if r.Method == v {
    		// Call the next middleware in chain
    		f(w, r)
    		return
    		}
    	}

    http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
    	return
        }
      }
    }

    // export as symbol named "Middleware"
    var Middleware middleware

To build the plugin library, we need to use the -buildmode=plugin flag and
specify the result name:

    go build -buildmode=plugin -o plugins/middlewares/method.so plugins/middlewares/method.go

    go build -buildmode=plugin -o plugins/controllers/genearal.so plugins/controllers/genearal.go

Now we can put all together to work starting the web server and test our
service.

    go build -o start -v

> *N.B. it works only on Linux, but with container we can solve this issue*


## Usage
- create middleware plugins under plugins/middlewares
- create handler plugins under plugins/handlers

- configure routes in the configurations/routes.json file
- configure server in the configurations/server.json file

### Tools
the create.sh script provide scaffold to make your middlewares and handlers.
executing command:

```bash
  $ ./create.sh handler mio
```
it will produce a plugins/handlers/mio.go file with the structure needed to use it in the server, as the same of

```bash
  $ ./create.sh middleware mio
```
that will create the plugins/middlewares/mio.go file.


### Building
once you have finish configurations and created the handlers/middlewares plugins, in shell run the command:

```bash
  $ make build
```

if you want to remove all compiled files, run:

```bash
  $ make clean
```
### Test
the project have a test that work for standard configuration and plugin, you need to edit this if you want to test your own implementation

## TODO
### Example
- create a ipfilter middleware
- create basic auth middleware
- create a "only-admin-access" middleware
### Test
- routing configuration test: searh for duplicated or wrong path, search for required plugins
- performance test
### Desired features
- server config to REDIRECT HTTP TO HTTPS
- middleware: CLIENT AUTHENTICATION
- server config to enable HTTPS: use crypto/tls package with ability to rotate TLS session ticket keys by default
- JWT API auth for javascript frontend framework like angular
- csrf token for request validation
