package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CardsController struct {
	collection *mongo.Collection
}

func NewCardsController(collection *mongo.Collection) *CardsController {
	return &CardsController{collection}
}

func (cc *CardsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Handle different HTTP methods
	switch r.Method {
	case http.MethodGet:
		cc.GetCards(w, r)
	case http.MethodPost:
		w.WriteHeader(http.StatusNotImplemented)
	case http.MethodDelete:
		w.WriteHeader(http.StatusNotImplemented)
	case http.MethodPut:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (cc *CardsController) GetCards(w http.ResponseWriter, r *http.Request) {
	cursor, err := cc.collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal("Error occured grabbing all cards from the database.")
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var cards []bson.M
	for cursor.Next(context.Background()) {
		var card bson.M
		err := cursor.Decode(&card)
		if err != nil {
			log.Fatal("Error occured decoding card from collection.")
			log.Fatal(err)
		}
		cards = append(cards, card)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}
