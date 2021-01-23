package contact

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// MongoHandler ...
type MongoHandler struct {
	client   *mongo.Client
	database *mongo.Database
}

func InitDataLayer(address string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mongoHandlerCallTimeout)
	defer cancel()

	return mongo.Connect(ctx, options.Client().ApplyURI(address))
}

// MongoHandler ...
func NewHandler(databaseName string, client *mongo.Client) *MongoHandler {
	db := client.Database(databaseName)
	handler := &MongoHandler{
		client:   client,
		database: db,
	}

	return handler
}

func (mh *MongoHandler) DisconnectDataLayer() error {
	return mh.client.Disconnect(context.Background())
}

func (mh *MongoHandler) GetOne(c *Contact, filter interface{}) error {
	collection := mh.database.Collection("contact")
	ctx, cancel := context.WithTimeout(context.Background(), getDocumentCallTimeout)

	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(c)

	return err
}

func (mh *MongoHandler) Get(filter interface{}) []*Contact {
	collection := mh.database.Collection("contact")
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
	collection := mh.database.Collection("contact")
	ctx, cancel := context.WithTimeout(context.Background(), addDocumentCallTimeout)

	defer cancel()

	result, err := collection.InsertOne(ctx, c)

	return result, err
}

func (mh *MongoHandler) Update(update *Contact) (*mongo.UpdateResult, error) {
	collection := mh.database.Collection("contact")
	ctx, cancel := context.WithTimeout(context.Background(), updateDocumentCallTimeout)

	defer cancel()

	result, err := collection.UpdateMany(ctx,
		bson.M{
			"firstName": update.FirstName,
			"lastName":  update.LastName,
		}, bson.D{
			primitive.E{
				Key: "$set",
				Value: bson.M{
					"phoneNumber": update.PhoneNumber,
				},
			},
		},
	)

	return result, err
}

func (mh *MongoHandler) RemoveOne(filter interface{}) (*mongo.DeleteResult, error) {
	collection := mh.database.Collection("contact")
	ctx, cancel := context.WithTimeout(context.Background(), removeDocumentCallTimeout)

	defer cancel()

	result, err := collection.DeleteOne(ctx, filter)

	return result, err
}
