package database

import (
	"context"
	"log"
	"time"
	"github.com/AkshatJawne/LetsGOCustomers/server-based/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const connectionString string = ""
// Add MongoDB Configuration String Here

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURL(connectionString))
	if (err != nil) {
		log.Fatal(err)
	}
	// Create timeout with context package to prevent memory leaks
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// Connect to MongoDB 
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		log.Fatal(err)
	}

	return &DB{
		client: client,
	}
}

func (db *DB) GetCustomer(id string) *model.Customer {
	customerCollec := db.client.Database("graphql-customer-board").Collection("customers")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}
	var customer model.Customer
	err := customerCollec.FindOne(ctx, filter).Decode(&customer)
	if err != nil {
		log.Fatal(err)
	}
	return &customer
}

func (db *DB) GetCustomers() []*model.Customer {
	customerCollec := db.client.Database("graphql-customer-board").Collection("customers")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var customers []*model.Customer
	cursor, err := customerCollec.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(context.TODO(), &customers); err != nil {
		panic(err)
	}

	return customers
}

func (db *DB) CreateCustomer(customerInfo model.CreateCustomerInput) *model.Customer {
	customerCollec := db.client.Database("graphql-customer-board").Collection("customer")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	inserg, err := customerCollec.InsertOne(ctx, bson.M{"name": customerInfo.Name, "description": customerInfo.Description, "company": customerInfo.Company})

	if err != nil {
		log.Fatal(err)
	}

	insertedID := inserg.InsertedID.(primitive.ObjectID).Hex()
	returnCustomer := model.Customer{ID: insertedID, Name: customerInfo.Name, Description: customerInfo.Description, Company: customerInfo.Company}
	return &returnCustomer
}

func (db *DB) UpdateCustomer(customerID string, customerInfo model.DeleteCustomerInput) *model.Customer {
	customerCollec := db.client.Database("graphql-customer-board").Collection("customer")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	updateCustomerInfo := bson.M{}

	if customerInfo.Title != nil {
		updateCustomerInfo["title"] = customerInfo.Title
	}
	if customerInfo.Description != nil {
		updateCustomerInfo["description"] = customerInfo.Description
	}
	if customerInfo.Company != nil {
		updateCustomerInfo["company"] = customerInfo.Company
	}

	_id, _ := primitive.ObjectIDFromHex(customerID)
	filter := bson.M{"_id": _id}
	update := bson.M{"$set": updateCustomerInfo}

	results := customerCollec.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var customer model.Customer

	if err := results.Decode(&customer); err != nil {
		log.Fatal(err)
	}

	return &customer
}

func (db *DB) DeleteCustomer(customerID string) *model.DeleteCustomer {
	customerCollec := db.client.Database("graphql-customer-board").Collection("customer")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_id, _ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id": _id}
	_, err := customerCollec.DeleteOne(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	return &model.DeleteJobResponse{DeletedCustomerId: customerID}
}
