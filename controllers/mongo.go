package controllers

import (
	"context"
	"log"
	"net/http"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoConnection struct {
	client    *mongo.Client
	cardRegex *regexp.Regexp
}

// The following regular expression only works for simple queries
// This will be replaced as the queries get more complex.
// This is to test functionality first.
func (mc mongoConnection) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Fatal("This is not implemented yet!")
}

// Constructor
func newMongoConnection() *mongoConnection {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	println("Connected to MongoDB!")

	return &mongoConnection{
		client:    client,
		cardRegex: regexp.MustCompile(`^/card/(name|color|manacost|type)=(.+)?$`),
	}
}

// Get a card
func (connection *mongoConnection) GetCardByName(name string) primitive.M {
	collection := connection.client.Database("MtgAppDatabase").Collection("Card")

	var result bson.M
	err := collection.FindOne(context.TODO(), bson.M{"name": name}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return result
}
