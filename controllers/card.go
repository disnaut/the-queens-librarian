package controllers

/*
@TODO: Create function for grabbing page and page size parameters from URL
*/
import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
			w.WriteHeader(http.StatusNotImplemented)
		case http.MethodDelete:
			w.WriteHeader(http.StatusForbidden)
		case http.MethodPut:
			w.WriteHeader(http.StatusNotImplemented)
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
	// Parse Name Parameter
	name := r.URL.Query().Get("name")

	//Parse the page parameters from the request URL
	//This accepts jQuery parameters. Example: http://localhost:8080/users?page=2&pageSize=25
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1 //default to page 1 if the query parsing breaks.
	}

	//Parse the page size parameters from the request URL
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pagesize"))
	if err != nil {
		pageSize = 10 //default to size 10 if the query parsing breaks
	}

	//Create a regex based on the name parameter
	pattern := bson.M{"$regex": name, "$options": "i"}

	//Construct query
	filter := bson.M{"name": pattern}

	//Calculate the number of documents to skip based on the page number and page
	skip := (page - 1) * pageSize

	//Query cards collection with a limit and skip
	cursor, err := cc.collection.Find(context.Background(), filter, options.Find().SetLimit(int64(pageSize)).SetSkip(int64(skip)))
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
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cards)
}
