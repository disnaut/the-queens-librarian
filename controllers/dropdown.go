package controllers

import (
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

type DropdownCore struct {
    collection *mongo.Collection
}

func NewDropdownCore(collection *mongo.Collection) DropdownCore {
    return DropdownCore{collection: collection}
}

func (dc *DropdownCore) ServeHttp(w http.ResponseWriter , r *http.Request) {

    if r.URL.Path == "/sets" {
        switch r.Method {
        case "GET":
            dc.getSets(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    } else if r.URL.Path == "/rarities" {
        switch r.Method {
        case "GET":
            dc.getRarities(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    }
}

func (dc *DropdownCore) getSets(w http.ResponseWriter, r *http.Request) {
}

func (dc *DropdownCore) getRarities(w http.ResponseWriter, r *http.Request) {
}


