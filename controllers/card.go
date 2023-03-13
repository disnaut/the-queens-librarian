package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CardsController struct {
	collection *mongo.Collection
}

func NewCardsController(collection *mongo.Collection) *CardsController {
	return &CardsController{collection}
}

type CardQueryParams struct {
	name        string
	colors      []string
	cardType    string
	artist      string
	keywords    []string
	set         string
	mana        int
	manaCompare string
	rarity      string
}

func (cc *CardsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/cards" {
		//Handle different HTTP methods
		switch r.Method {
		case http.MethodGet:
			cc.SearchCards(w, r)
		case http.MethodPost:
			w.WriteHeader(http.StatusForbidden)
		case http.MethodDelete:
			w.WriteHeader(http.StatusForbidden)
		case http.MethodPut:
			w.WriteHeader(http.StatusForbidden)
		}
	} else {
		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusNotImplemented)
		case http.MethodPost:
			w.WriteHeader(http.StatusNotImplemented)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNotImplemented)
		case http.MethodPut:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func (cc *CardsController) SearchCards(w http.ResponseWriter, r *http.Request) {
	var cardQuery CardQueryParams
	var queryCollection []bson.M

	GetQueryParams(r, &cardQuery)

	nameRegex := bson.M{"$regex": cardQuery.name, "$options": "i"}
	artistRegex := bson.M{"$regex": cardQuery.artist, "$options": "i"}
	setRegex := bson.M{"$regex": cardQuery.set, "$options": "i"}
	cardTypesRegex := bson.M{"$regex": cardQuery.cardType, "$options": "i"}

	//assign regex patterns to respective parts of the query
	name := bson.M{"name": nameRegex}
	artist := bson.M{"artist": artistRegex}
	set := bson.M{"set_name": setRegex}
	cardTypes := bson.M{"type_line": cardTypesRegex}

	queryCollection = append(queryCollection, name, artist, set, cardTypes)

	/* $and setup for non regex types */
	if len(cardQuery.colors) != 0 {
		colors := bson.M{"colors": cardQuery.colors}
		queryCollection = append(queryCollection, colors)
	}

	if len(cardQuery.keywords) != 0 {
		keywords := bson.M{"keywords": cardQuery.keywords}
		queryCollection = append(queryCollection, keywords)
	}

	if len(cardQuery.rarity) != 0 {
		rarity := bson.M{"rarity": cardQuery.rarity}
		queryCollection = append(queryCollection, rarity)
	}

	if len(cardQuery.manaCompare) != 0 {
		manaCompare := bson.M{cardQuery.manaCompare: cardQuery.mana}
		mana := bson.M{"cmc": manaCompare}
		queryCollection = append(queryCollection, mana)
	}

	query := bson.M{"$and": queryCollection}

	cursor, err := cc.collection.Find(context.Background(), query)
	if err != nil {
		log.Fatal("Error occured while getting cards from collection.")
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	//Iterate over the results and write to the response
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

func (cc *CardsController) UpdateCollection(w http.ResponseWriter, r *http.Request) {
	// IMPORTANT: We need to make sure that collection and amount exist as fields.
	// First, we grab parameters which would be the following
	//		- Card ID
	//		- Amount of card
	// These parameters would be gotten via a json payload that is an array of objects those members would be the parameters stated above
	// Second, Decode the JSON into an array we can access.
	// Third, get a connection to TheQueensLibrary, and the Collection of Cards
	// Forth, Loop through the JSON array and do the following operations:
	//		- Get the card data from the database.
	//		- Check if the 'collection field is false'
	//			- if false, then set to true
	//		- Update amount to be the amount of card
	// 		- Loop until every card is updated
}

func GetQueryParams(r *http.Request, card *CardQueryParams) {
	/* Strings */
	card.name = r.URL.Query().Get("name")
	card.artist = r.URL.Query().Get("artist")
	card.set = r.URL.Query().Get("set")
	card.rarity = r.URL.Query().Get("rarity")
	card.cardType = r.URL.Query().Get("types")

	/* Arrays */
	colors := r.URL.Query().Get("colors")
	keywords := r.URL.Query().Get("keywords")

	if len(colors) != 0 {
		card.colors = strings.Split(colors, ",")
	}

	if len(keywords) != 0 {
		card.keywords = strings.Split(keywords, ",")
	}

	manaParam := r.URL.Query().Get("mana")
	if len(manaParam) != 0 {
		if strings.HasPrefix(manaParam, "lte") {
			manaParam, _ = strings.CutPrefix(manaParam, "lte")
			card.mana, _ = strconv.Atoi(manaParam)
			card.manaCompare = "$lte"
		} else if strings.HasPrefix(manaParam, "gte") {
			manaParam, _ = strings.CutPrefix(manaParam, "gte")
			card.mana, _ = strconv.Atoi(manaParam)
			card.manaCompare = "$gte"
		} else {
			var err error
			card.mana, err = strconv.Atoi(manaParam)
			if err != nil {
				log.Panicln(err)
				card.mana = -1
				card.manaCompare = "$gte"
			} else {
				card.manaCompare = "$eq"
			}
		}
	}
}
