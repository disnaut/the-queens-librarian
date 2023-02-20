package main

import (
	"net/http"

	"github.com/the-queens-library/webservice/controllers"
)

func main() {
	controllers.RegisterControllers()
	http.ListenAndServe(":3000", nil)
}
