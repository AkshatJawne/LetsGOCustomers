package handlers 

import (
	"net/http"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/AkshatJawne/go-serverless/pkg/user"
)

const ErrorMethodNotAllowed = "Method not allowed" 

type ErrorBody struct {
	ErrorMsg 	*string 	`json:"error,omitempty"`
} // ErrorBody defined as struct with ErrorMsg as pointer to string which contains json

func GetUser(req events.APIGatewayProxyRequest, tableName string, dynaClient, dynaClient dynamodbiface.DynamoDBAPI) (
	*events.APIGatewayProxyResponse, error) {

		email := req.QueryStringParameters["email"] // get email from query string parameters
		
		if len(email) > 0 {
			result, err := user.FetchUser(email, tableName, dynaClient) // fetch user if there is inputted email
			if (err != nil) {
				return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())}) // return error repsone if error
			}
			return apiResponse(http.StatusOK, result) // return result with OK status if no error
		}

		// if we want to fetch multiple users
		result, err := user.FetchUsers(tableName, dynaClient) // fetch users if no inputted email

		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})  // return error repsone if error
		}

		return apiResponse(http.StatusOK, result)  // return result with OK status if no error
} 

func CreateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*events.APIGatewayProxyResponse, error) {
		
		result, err := user.CreateUser(req, tableName, dynaClient) // create user with inputted request

		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error()),})	// return error repsone if error 
		}

		return apiResponse(http.StatusCreated, result) // return result with Created status if no error
}

func UpdateUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*events.APIGatewayProxyResponse, error) {
		
		result, err := user.UpdateUser(req, tableName, dynaClient) // update user with inputted request

		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())}) // return error repsone if error 
		}

		return apiResponse(http.StatusOK, result)  // return result with OK status if no error

} 

func DeleteUser(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (
	*events.APIGatewayProxyResponse, error) {

		result, err := user.DeleteUser(req, tableName, dynaClient) // delete user with inputted request

		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())}) // return error repsone if error
		}

		return apiResponse(http.StatusOK, nil) // return result with OK status if no error 
}

func UnhandledMethod()(*events.APIGatewayProxyResponse, error) {
	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed) 
	// return apiResponse function with status and error message if http method not one of the above
}