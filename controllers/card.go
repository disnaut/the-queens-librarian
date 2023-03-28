package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"runtime"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Could be CardCore - Kirk
//SearchController or CardsHandler
type SearchController struct {
	collection *mongo.Collection
}

//Just need to return the value instead of a pointer
func NewSearchController(collection *mongo.Collection) SearchController {
	return SearchController{collection}
}

//Rename to QueryParams or Params
//Noun/Verb
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

func (cc *SearchController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/cards" {
		//Handle different HTTP methods
		switch r.Method {
        case http.MethodGet:
			cc.searchCards(w, r)
		case http.MethodPost:
			w.WriteHeader(http.StatusForbidden)
		case http.MethodDelete:
			w.WriteHeader(http.StatusForbidden)
		case http.MethodPut:
			w.WriteHeader(http.StatusForbidden)
        }
    }
}

// Try to implement one thing of concurrency with weight groups
// 900 bucks for a course from Arden Labs
// THIS IS WORTH IT.
// Arden Labs -> Have Youtube videos, and could learn this freely.
// Refactor, Decouple

// Private function
// searchCards will grab cards from mongodb based on jQuery arguments
func (cc *SearchController) searchCards(w http.ResponseWriter, r *http.Request) {
	var cardQuery CardQueryParams

	parseUrl(r, &cardQuery)
	query := createQuery(cardQuery)

	opts := options.Find().SetSort(bson.D{{"name", 1}})

	cursor, err := cc.collection.Find(context.Background(), query, opts)
	if err != nil {
		log.Fatal("Error occured while getting cards from collection.")
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	// Create a channel to receive the results
	results := make(chan bson.M)

	// Read all collections into a queue
	var collectionQueue []bson.M
	for cursor.Next(context.Background()) {
		var card bson.M
		err := cursor.Decode(&card)
		if err != nil {
			log.Fatal("Error occured while grabbing card.")
			log.Fatal(err)
		}
		collectionQueue = append(collectionQueue, card)
	}

    var mu sync.Mutex

	// Spawn a pool of worker goroutines to process the results in parallel
	numWorkers := runtime.NumCPU()
	wg := sync.WaitGroup{}
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			for card := range results {
                mu.Lock()
				w.Header().Set("Content-Type", "application/json")
				err := json.NewEncoder(w).Encode(card)
				if err != nil {
					log.Fatal("Error occured while encoding card to JSON.")
					log.Fatal(err)
				}
				w.(http.Flusher).Flush()
                mu.Unlock()
			}
		}()
	}

	// Send out HTTP responses using wait groups
	for _, card := range collectionQueue {
		results <- card
	}
	close(results)
	wg.Wait()

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}

// Private Function
// CreateQuery creates a query based jQuery arguments
// that are parsed within this function
func createQuery(cardQuery CardQueryParams) bson.M {
    var queryCollection []bson.M
 
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

	return bson.M{"$and": queryCollection}
}

// Private function
// ParseUrl gets jQuery arguments from Url
func parseUrl(r *http.Request, card *CardQueryParams) {
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
