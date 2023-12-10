package main

import (
	"os"
	"github.com/aws/aws-sdk-go/aws" 
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/AkshatJawne/go-serverless/pkg/handlers"
)

const tableName = "LambdaInGoUser"  

var (
	dynaClient dynamodbiface.DynamoDBAPI // dynaClient defined as DynamoDBAPI interface from dynamodbiface.go sdk
)

func main() {
	region := os.getEnv("AWS_REGION");

	awsSession, err := session.NewSession(&aws.Config{
		region: aws.String(region),
	}) // awsSession defined with NewSession function from aws.go sdk

	if err != nil { // if error, print error
		return
	}

	dynaClient := dynamodb.New(awsSession) // dynaClient defined with New from dynamodb.go sdk
	lambda.Start(handler) 
}

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod { // switch statement with different cases for each method with corresponding handlers
		case "GET": 
			return handlers.GetUser(req, tableName, dynaClient)
		case "POST": 
			return handlers.CreateUser(req, tableName, dynaClient)
		case "PUT": 
			return handlers.UpdateUser(req, tableName, dynaClient)
		case "DELETE": 
			return handlers.DeleteUser(req, tableName, dynaClient)
		default:
			return handlers.UnhandledMethod() // if method not defined, return UnhandledMethod function from handlers
	}
}