package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CollectionController struct {
    collection *mongo.Collection
}

func NewCollectionController(collection *mongo.Collection) *CollectionController {
    return &CollectionController{collection}
}

func (cc *CollectionController) ServeHTTP(w http.ResponseWriter, 
                                            r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        cc.Collection(w, r)
    case http.MethodPost:
        cc.AddToCollection(w, r)
    case http.MethodDelete:
        cc.RemoveFromCollection(w, r)
    }
}


func (cc *CollectionController) Collection(w http.ResponseWriter, 
                                            r *http.Request) {
    cursor, err := cc.collection.Find(context.Background(), bson.D{})
    if err != nil {
        log.Fatal("There was something wrong getting the Documents from Collection")
        log.Fatal(err)
    }
    defer cursor.Close(context.Background())

    var card bson.M
    for cursor.Next(context.Background()) {
    	err := cursor.Decode(&card)
		if err != nil {
			log.Fatal("Error occured while grabbing card.")
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(card)
		w.(http.Flusher).Flush()
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err) //Construct and query
    }
}

func (cc *CollectionController) AddToCollection(w http.ResponseWriter, 
                                                r *http.Request) {
    var cards []interface{}
    err := json.NewDecoder(r.Body).Decode(&cards)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

    context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, error := cc.collection.InsertMany(context, cards)
    if error != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func (cc *CollectionController) RemoveFromCollection(w http.ResponseWriter, 
                                                        r *http.Request) {
    var cards []interface{}
    err := json.NewDecoder(r.Body).Decode(&cards)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    //Create waitgroup to track the completion of worker goroutines
    var wg sync.WaitGroup

    //Create a a channell to recieve errors from worker goroutines
    errorChannel := make(chan error)

    //Create a channell to send IDs to worker goroutines
    idChannel := make(chan string)

    //Spawn a pool of worker goroutines to delete cards in parallel
    numWorkers := runtime.NumCPU()
    wg.Add(numWorkers)

    for i := 0; i < numWorkers; i++ {
        go func() {
            defer wg.Done()
            for id := range idChannel {
                filter := bson.M{"_id": id}

                _, err := cc.collection.DeleteOne(context.Background(), filter)
                if err != nil {
                    errorChannel <- err
                    return
                }
            }
        }()
    }
}
