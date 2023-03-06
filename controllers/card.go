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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CardsController struct {
	collection *mongo.Collection
}

func NewCardsController(collection *mongo.Collection) *CardsController {
	return &CardsController{collection}
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
func (cc *CardsController) SearchCards(w http.ResponseWriter, r *http.Request) {
	/* region: Grabbing Query Parameters */
	name, colors_arr, types, artist, keywords_arr, set, manacost, rarity, page, pageSize := GetQueryParams(r)

	println(manacost) //To keep manacost available while options are considered for what can be done.
	var and []bson.M
	//Create a regex based on certain params
	name_pattern := bson.M{"$regex": name, "$options": "i"}
	artist_pattern := bson.M{"$regex": artist, "$options": "i"}
	set_pattern := bson.M{"$regex": set, "$options": "i"}
	types_pattern := bson.M{"$regex": types, "$options": "i"}

	//assign regex patterns to respective queries
	name_query := bson.M{"name": name_pattern}
	artist_query := bson.M{"artist": artist_pattern}
	set_query := bson.M{"set_name": set_pattern}
	types_query := bson.M{"type_line": types_pattern}

	and = append(and, name_query, artist_query, set_query, types_query)
	/* endregion */

	/* $and setup for array types */
	if len(colors_arr) != 0 {
		colors_query := bson.M{"colors": colors_arr}
		and = append(and, colors_query)
	}

	if len(keywords_arr) != 0 {
		keywords_query := bson.M{"keywords": keywords_arr}
		and = append(and, keywords_query)
	}

	/*
		Could use an enum to check if it is one of the given rarities
		- common
		- uncommon
		- rare
		- mythic rare
	*/
	if len(rarity) != 0 {
		rarity_query := bson.M{"rarity": rarity}
		and = append(and, rarity_query)
	}

	/* region: setting up query */
	query := bson.M{"$and": and}

	//Calculate the number of documents to skip based on the page number and page
	skip := (page - 1) * pageSize
	/* endregion */

	//Query cards collection with a limit and skip
	cursor, err := cc.collection.Find(context.Background(), query, options.Find().SetLimit(int64(pageSize)).SetSkip(int64(skip)))
	if err != nil {
		log.Fatal("Error occured while getting cards from collection.")
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	//Iterate over the results and write to the response
	var cards []bson.M
	for cursor.Next(context.Background()) {
		var card bson.M
		err := cursor.Decode(&card)
		if err != nil {
			log.Fatal("Error occured while grabbing card.")
			log.Fatal(err)
		}
		cards = append(cards, card)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err) //Construct and query

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}

/*
Parameters that we can get from this:
  - Name => String regex [X]
  - Color => []string [X]
  - Type => []string [X]
  - Artist => string regex [X]
  - Keywords => []string [X]
  - Set => name regex [X]
  - Manacost => using cmc, which is a number.
  - Rarity => string [X]
    artist := r.URL.Query().Get("artist")
    set := r.URL.Query().Get("set")
    rarity := r.URL.Query().Get("rarity")
    types := r.URL.Query().Get("type")
  - page
  - pagesize
*/
func GetQueryParams(r *http.Request) (string, []string, string, string, []string, string, int8, string, int, int) {
	/* Strings */
	name := r.URL.Query().Get("name")
	artist := r.URL.Query().Get("artist")
	set := r.URL.Query().Get("set")
	rarity := r.URL.Query().Get("rarity")
	types := r.URL.Query().Get("types")

	/* Arrays */
	colors := r.URL.Query().Get("colors")
	keywords := r.URL.Query().Get("keywords")

	var colors_arr []string
	var keywords_arr []string

	if len(colors) != 0 {
		colors_arr = strings.Split(colors, ",")
	}

	if len(keywords) != 0 {
		keywords_arr = strings.Split(keywords, ",")
	}

	/* Numbers */
	manacost, err := strconv.Atoi(r.URL.Query().Get("manacost"))
	if err != nil {
		manacost = -1
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1 //default to page 1 if the query parsing breaks.
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pagesize"))
	if err != nil {
		pageSize = 10 //default to size 10 if the query parsing breaks
	}

	return name, colors_arr, types, artist, keywords_arr, set, int8(manacost), rarity, page, pageSize
}
