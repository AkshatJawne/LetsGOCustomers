package handlers 

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
)

func apiResponse (status int, body interface {}) (*events.APIGatewayProxyRequest, error) {
	resp := events.APIGatewayProxyResponse {Headers: map[string]string{"Content-Type": "application/json"}, StatusCode: status}  
	// resp defined as APIGatewayProxyResponse struct from events.go sdk

	stringBody, _ := json.Marshal(body) // GO dodes not recognize jsons natively, thus, need to use json Marshal
	resp.Body = string(stringBody) 
	return &resp, nil
}