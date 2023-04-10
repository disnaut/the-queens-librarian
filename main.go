package main

import (
	"context"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/the-queens-librarian/webAPI/controllers"
)

type appContext struct {
	mongoURI string
}

func main() {

	appCtx := appContext{
		mongoURI: "mongodb://localhost:27017",
	}

	//Start Server
	router := http.NewServeMux()
	router.HandleFunc("/cards", appCtx.CardHandler)
    router.HandleFunc("/collection", appCtx.CollectionHandler)
	http.ListenAndServe(":3000", router)
}

func (appCtx *appContext) CollectionHandler(w http.ResponseWriter, r *http.Request) {
    clientOptions := options.Client().ApplyURI(appCtx.mongoURI)
    client, err := mongo.NewClient(clientOptions)
    if err != nil {
        log.Fatal("Something went wrong creating the mongo client.")
        log.Fatal(err)
    }

    err = client.Connect(context.Background())
    if err != nil {
        log.Fatal("Something went wrong connecting to the client.")
        log.Fatal(err)
    }

    defer client.Disconnect(context.Background())

    collection := client.Database("TheQueensLibrary").Collection("Collection")

    collectionController := controllers.NewCollectionController(collection)
    
    collectionController.ServeHTTP(w, r)

}

func (appCtx *appContext) CardHandler(w http.ResponseWriter, r *http.Request) {
	//Create a new client object for each request
	clientOptions := options.Client().ApplyURI(appCtx.mongoURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal("Something went wrong creating the mongo client.")
		log.Fatal(err)
	}
    
    w.Header().Set("Access-Control-Allow-Origin", "*")
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal("Something went wrong connecting to the client.")
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	//Access the cards collection
	collection := client.Database("TheQueensLibrary").Collection("Cards")

	//Use the collection in the cards controller
	searchController := controllers.NewSearchController(collection)

	//Handle the request using the card controller
	searchController.ServeHTTP(w, r)
}
