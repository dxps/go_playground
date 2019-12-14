package route

import (
	"goyave-quickrun/http/controller/hello"
	"goyave-quickrun/http/request/echorequest"

	"github.com/System-Glitch/goyave/v2"
)

// Routing is an essential part of any Goyave application.
// Routes definition is the action of associating a URI, sometimes having
// parameters, with a handler which will process the request and respond to it.

// Routes are defined in routes registrer functions.
// The main route registrer is passed to "goyave.Start()" and is executed
// automatically with a newly created root-level router.

// Register all the application routes. This is the main route registrer.
func Register(router *goyave.Router) {
	// Register your routes here

	// Route without validation
	router.Route("GET", "/hello", hello.SayHi, nil)

	// Route with validation
	router.Route("GET", "/echo", hello.Echo, echorequest.Echo)
}
