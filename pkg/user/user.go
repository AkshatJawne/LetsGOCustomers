package user 

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-lambda-go/events"
	"github.com/AkshatJawne/go-serverless/pkg/validators"
)

var (
	ErrorFailedToFetchRecord = "Failed to fetch record!"
	ErrorFailedToUnmarshalRecord = "Failed to unmarshal record!"
	ErrorInvalidUser = "Invalid user data!"
	ErrorInvalidEmail = "Invalid email!" 
	ErrorCouldNotMarshalItem = "Could not marshal item!"
	ErrorCouldNotDeleteItem = "Could not delete item!"
	ErrorCouldNotPutItem = "Could not Dynamo put item!"
	ErrorUserAlreadyExists = "User already exists!"
	ErrorDoesNotExist = "User does not exist!"
)

type User struct {
	Email 		string 		`json:"email"`
	FirstName	string 		`json:"firstName"`
	LastName	string 		`json:" lastName"`
} // User defined as struct with Email, FirstName, LastName as strings which contain json

func FetchUser(email, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error) {
	input := &dynamodb.GetItemInput{ 
		Key: map[string]*dynamodb.AttributeValue{
			"email" : {
				S: aws.String(email), // S is used to indicate that the email key is of type string
			}, 
		}, 
		TableName: aws.String(tableName), // TableName is the name of the table in DynamoDB 
	} 

	result, err := dynaClient.GetItem(input)  // get item from DynamoDB table

	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord) 
	}

	item := new(User) // creating new instance of User struct 
	err = dynamodbattribute.UnmarshalMap(result.Item, item) // UnmarshalMap is used to convert the result.Item map to a User struct

	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return item, nil 

}

func FetchUsers(tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*[]User, error) {
	input := &dynamodb.ScanInput {
		TableName : aws.String(tableName), //   is the name of the table in DynamoDB
	}

	result, err := dynaClient.Scan(input) // scan the table in DynamoDB for the input

	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new([]User) // creating new instance of User struct with slice of users 
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, item) // UnmarshalListOfMaps is used to convert the result.Item map to a User struct
	return item, nil 
}

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error,) {
	var u User // create instance of User struct

	if err := json.Unmarshal([]byte(req.Body), &u); err!=nil {
		return nil, errors.New(ErrorInvalidUser)
	} // unmarshal the request body

	if !validators.IsEmailValid (u.Email) {
		return nil, errors.New(ErrorInvalidEmail) // if email is not valid, return error
	}

	currentUser, _ := FetchUser(u.Email, tableName, dynaClient) // attempt to fetch user from table

	if currentUser != nil && len(currentUser.Email) != 0  {
		return nil, errors.New(ErrorUserAlreadyExists) // if user already exists, return error
	} 

	av, err := dynamodbattribute.MarshalMap(u) // MarshalMap is used to convert the User struct to a map

	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem) // if error in marshalling, return error
	}

	input := &dynamodb.PutItemInput {
		Item: av, // Item is the map of attribute name/value pairs
		TableName: aws.String(tableName), // TableName is the name of the table in DynamoDB
	}

	_, err = dynaClient.PutItem(input) // put item into DynamoDB table
	if err != nil {
		return nil, errors.New(ErrorCouldNotPutItem) // if error in putting item, return error
	}

	return &u, nil 
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI)(*User, error, ) {
	var u User // create instance of User struct

	if err := json.Unmarshal([]byte(req.Body), &u); err!=nil {
		return nil, errors.New(ErrorInvalidUser)
	} // unmarshal the request body

	currentUser , _ := FetchUser(u.Email, tableName, dynaClient) // attempt to fetch user from table

	if currentUser == nil || len(currentUser.Email) == 0  {
		return nil, errors.New(ErrorDoesNotExist) // if user does not exist, return error
	}

	av, err := dynamodbattribute.MarshalMap(u) // MarshalMap is used to convert the User struct to a map

	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem) // if error in marshalling, return error 
	}

	input := &dynamodb.PutItemInput {
		Item: av,
		TableName: aws.String(tableName), 
	}

	_, err = dynaClient.PutItem(input) // put item into DynamoDB table
	if err != nil {
		return nil, errors.New(ErrorCouldNotPutItem) // if error in putting item, return error
	}

	return &u, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {
	email := req.QueryStringParameters["email"] // get email from query string parameters

	input := &dynamodb.DeleteItemInput {
		Key: map[string]*dynamodb.AttributeValue{
			"email" : {
				 S: aws.String(email), // S is used to indicate that the email key is of type string
			},
		},
		TableName: aws.String(tableName), // TableName is the name of the table in DynamoDB
	}

	_, err := dynaClient.DeleteItem(input)  // delete item from DynamoDB table

	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem) // if error in deleting item, return error
	}

	return nil
}