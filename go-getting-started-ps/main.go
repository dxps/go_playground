package main

import (
	"net/http"

	"github.com/pluralsight/webservice/controllers"
)

func main() {
	controllers.RegisterControllers()
	_ = http.ListenAndServe(":3000", nil)
}
