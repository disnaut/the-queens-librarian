//Set up when a network request is recieved
//and go to the right controller.

package controllers

import (
	"encoding/json"
	"io"
	"net/http"
)

func RegisterControllers() {
	uc := newUserController()
	mc := newMongoConnection()

	http.Handle("/users", *uc)
	http.Handle("/users/", *uc)
	http.Handle("/card/", *mc)
}

func encodeResponseAsJson(data interface{}, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.Encode(data)
}
