package controllers

/*
TODO: Get manacost up and running
TODO: Change out query param variables to pointers
TODO: Update the response to multiple http requests one after another rather than array
*/

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
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
	name     string
	colors   []string
	cardType string
	artist   string
	keywords []string
	set      string
	mana     string
	rarity   string
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

// Because we have a collection of over 20000+ cards, we'll need to paginate the responses
/*
This isn't the most efficient way as it grabs EVERYTHING.
Would be worth looking into getting requests one at a time.
*/
func (cc *CardsController) SearchCards(w http.ResponseWriter, r *http.Request) {
	/* region: Grabbing Query Parameters */
	var cardQuery CardQueryParams

	GetQueryParams(r, &cardQuery)

	var query_collection []bson.M
	//Create a regex based on certain params
	name_pattern := bson.M{"$regex": cardQuery.name, "$options": "i"}
	artist_pattern := bson.M{"$regex": cardQuery.artist, "$options": "i"}
	set_pattern := bson.M{"$regex": cardQuery.set, "$options": "i"}
	types_pattern := bson.M{"$regex": cardQuery.cardType, "$options": "i"} //Change to regular query, something like only Artifacts. No solution for if multiple types are wanted in the query

	//assign regex patterns to respective queries
	name_query := bson.M{"name": name_pattern}
	artist_query := bson.M{"artist": artist_pattern}
	set_query := bson.M{"set_name": set_pattern}
	types_query := bson.M{"type_line": types_pattern}

	query_collection = append(query_collection, name_query, artist_query, set_query, types_query)

	/* $and setup for non regex types */
	if len(cardQuery.colors) != 0 {
		colors_query := bson.M{"colors": cardQuery.colors}
		query_collection = append(query_collection, colors_query)
	}

	if len(cardQuery.keywords) != 0 {
		keywords_query := bson.M{"keywords": cardQuery.keywords}
		query_collection = append(query_collection, keywords_query)
	}

	if len(cardQuery.rarity) != 0 {
		rarity_query := bson.M{"rarity": cardQuery.rarity}
		query_collection = append(query_collection, rarity_query)
	}

	query := bson.M{"$and": query_collection} //Potentially figure out a way to handle an or statement.w
	//Calculate the number of documents to skip based on the page number and page

	//Query cards collection with a limit and skip
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

func GetQueryParams(r *http.Request, cardPtr *CardQueryParams) {
	/* Strings */
	*&cardPtr.name = r.URL.Query().Get("name")
	*&cardPtr.artist = r.URL.Query().Get("artist")
	*&cardPtr.set = r.URL.Query().Get("set")
	*&cardPtr.rarity = r.URL.Query().Get("rarity")
	*&cardPtr.cardType = r.URL.Query().Get("types")
	*&cardPtr.mana = r.URL.Query().Get("manaCost")

	/* Arrays */
	colors := r.URL.Query().Get("colors")
	keywords := r.URL.Query().Get("keywords")

	if len(colors) != 0 {
		*&cardPtr.colors = strings.Split(colors, ",")
	}

	if len(keywords) != 0 {
		*&cardPtr.keywords = strings.Split(keywords, ",")
	}
}
