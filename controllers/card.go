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

func (cc *CardsController) SearchCards(w http.ResponseWriter, r *http.Request) {
	var cardQuery CardQueryParams
	var query_collection []bson.M

	GetQueryParams(r, &cardQuery)

	name_pattern := bson.M{"$regex": cardQuery.name, "$options": "i"}
	artist_pattern := bson.M{"$regex": cardQuery.artist, "$options": "i"}
	set_pattern := bson.M{"$regex": cardQuery.set, "$options": "i"}
	types_pattern := bson.M{"$regex": cardQuery.cardType, "$options": "i"}

	//assign regex patterns to respective parts of the query
	name := bson.M{"name": name_pattern}
	artist := bson.M{"artist": artist_pattern}
	set := bson.M{"set_name": set_pattern}
	card_types := bson.M{"type_line": types_pattern}

	query_collection = append(query_collection, name, artist, set, card_types)

	/* $and setup for non regex types */
	if len(cardQuery.colors) != 0 {
		colors := bson.M{"colors": cardQuery.colors}
		query_collection = append(query_collection, colors)
	}

	if len(cardQuery.keywords) != 0 {
		keywords := bson.M{"keywords": cardQuery.keywords}
		query_collection = append(query_collection, keywords)
	}

	if len(cardQuery.rarity) != 0 {
		rarity := bson.M{"rarity": cardQuery.rarity}
		query_collection = append(query_collection, rarity)
	}

	if len(cardQuery.mana) != 0 {
		mana := strings.Split(cardQuery.mana, "=")
		cmc, _ := strconv.Atoi(mana[1])
		var manaQuery bson.M
		if len(mana) == 0 {
			manaQuery = bson.M{"cmc": cmc}
		} else {
			if mana[1] == ">" {
				gte := bson.M{"$gte": cmc}
				manaQuery = bson.M{"cmc": gte}
			} else {
				lte := bson.M{"lte": cmc}
				manaQuery = bson.M{"cmc": lte}
			}
		}
		query_collection = append(query_collection, manaQuery)
	}

	query := bson.M{"$and": query_collection}

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

func GetQueryParams(r *http.Request, card *CardQueryParams) {
	/* Strings */
	card.name = r.URL.Query().Get("name")
	card.artist = r.URL.Query().Get("artist")
	card.set = r.URL.Query().Get("set")
	card.rarity = r.URL.Query().Get("rarity")
	card.cardType = r.URL.Query().Get("types")
	card.mana = r.URL.Query().Get("manaCost")

	/* Arrays */
	colors := r.URL.Query().Get("colors")
	keywords := r.URL.Query().Get("keywords")

	if len(colors) != 0 {
		card.colors = strings.Split(colors, ",")
	}

	if len(keywords) != 0 {
		card.keywords = strings.Split(keywords, ",")
	}
}
