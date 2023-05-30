package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
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
    var setNames[] string

    filter := bson.D{{}}

    results, err := dc.collection.Distinct(context.Background(), "set_name", filter)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting sets: %v", err), http.StatusInternalServerError)
        return
    }

    for _, value := range results {
        setName, ok := value.(string)
        if !ok {
            http.Error(w, fmt.Sprintf("Error getting set name: %v", err), http.StatusInternalServerError)
            return
        }
        setNames = append(setNames, setName)
    }

    w.Header().Set("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(setNames)
    w.(http.Flusher).Flush()
    if err != nil {
        log.Printf("Error encoding json: %v", err)
    }
}

func (dc *DropdownCore) getRarities(w http.ResponseWriter, r *http.Request) {
    var rarities[] string

    filter := bson.D{{}}

    results, err := dc.collection.Distinct(context.Background(), "rarity", filter)
    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting rarities: %v", err), http.StatusInternalServerError)
        return
    }

    for _, value := range results {
        rarity, ok := value.(string)
        if !ok {
            http.Error(w, fmt.Sprintf("Error getting rarity: %v", err), http.StatusInternalServerError)
            return
        }
        rarities = append(rarities, rarity)
    }

    w.Header().Set("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(rarities)
    w.(http.Flusher).Flush()
    if err != nil {
        log.Printf("Error encoding json: %v", err)
    }
}


