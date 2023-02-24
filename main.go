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
	router.HandleFunc("/cards", appCtx.cardHandler)
	http.ListenAndServe(":3000", router)
}

func (appCtx *appContext) cardHandler(w http.ResponseWriter, r *http.Request) {
	//Create a new client object for each request
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

	//Access the cards collection
	collection := client.Database("MtgAppDatabase").Collection("Card")

	//Use the collection in the cards controller
	cardController := controllers.NewCardsController(collection)

	//Handle the request using the card controller
	cardController.ServeHTTP(w, r)
}
