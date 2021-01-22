package contact

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	addDocumentCallTimeout    = 5 * time.Second
	getDocumentCallTimeout    = 5 * time.Second
	updateDocumentCallTimeout = 5 * time.Second
	removeDocumentCallTimeout = 10 * time.Second
	mongoHandlerCallTimeout   = 10 * time.Second
)

// DefaultDatabase ...
const DefaultDatabase = "contactstore"

// MongoHandler ...
type MongoHandler struct {
	client   *mongo.Client
	database string
}

// MongoHandler ...
func NewHandler(address string) *MongoHandler {
	ctx, cancel := context.WithTimeout(context.Background(), mongoHandlerCallTimeout)
	defer cancel()

	mongoClient, _ := mongo.Connect(ctx, options.Client().ApplyURI(address))
	handler := &MongoHandler{
		client:   mongoClient,
		database: DefaultDatabase,
	}

	return handler
}

func (mh *MongoHandler) GetOne(c *Contact, filter interface{}) error {
	// Will automatically create a collection if not available
	collection := mh.client.Database(mh.database).Collection("contact")
	ctx, cancel := context.WithTimeout(context.Background(), getDocumentCallTimeout)

	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(c)

	return err
}

func (mh *MongoHandler) Get(filter interface{}) []*Contact {
	collection := mh.client.Database(mh.database).Collection("contact")
	ctx, cancel := context.WithTimeout(context.Background(), getDocumentCallTimeout)

	defer cancel()

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	var result []*Contact

	for cur.Next(ctx) {
		contact := &Contact{}

		err = cur.Decode(contact)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, contact)
	}

	return result
}

func (mh *MongoHandler) AddOne(c *Contact) (*mongo.InsertOneResult, error) {
	collection := mh.client.Database(mh.database).Collection("contact")
	ctx, cancel := context.WithTimeout(context.Background(), addDocumentCallTimeout)

	defer cancel()

	result, err := collection.InsertOne(ctx, c)

	return result, err
}

func (mh *MongoHandler) Update(filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	collection := mh.client.Database(mh.database).Collection("contact")
	ctx, cancel := context.WithTimeout(context.Background(), updateDocumentCallTimeout)

	defer cancel()

	result, err := collection.UpdateMany(ctx, filter, update)

	return result, err
}

func (mh *MongoHandler) RemoveOne(filter interface{}) (*mongo.DeleteResult, error) {
	collection := mh.client.Database(mh.database).Collection("contact")
	ctx, cancel := context.WithTimeout(context.Background(), removeDocumentCallTimeout)

	defer cancel()

	result, err := collection.DeleteOne(ctx, filter)

	return result, err
}
