package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
	"github.com/mattgillard/user-pii-demo/common"
)

var validate *validator.Validate = validator.New()

func clientError(status int) (events.APIGatewayProxyResponse, error) {

	return events.APIGatewayProxyResponse{
		Body:       http.StatusText(status),
		StatusCode: status,
	}, nil
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	log.Println(err.Error())

	return events.APIGatewayProxyResponse{
		Body:       http.StatusText(http.StatusInternalServerError),
		StatusCode: http.StatusInternalServerError,
	}, nil
}

func processGetPassport(ctx context.Context, id string) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received GET todo request with id = %s", id)

	todo, err := common.GetItem(ctx, id, "Passport")
	if err != nil {
		return serverError(err)
	}

	if todo == nil {
		return clientError(http.StatusNotFound)
	}

	json, err := json.Marshal(todo)
	if err != nil {
		return serverError(err)
	}
	log.Printf("Successfully fetched todo item %s", json)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(json),
	}, nil
}

func processGetTodos(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	log.Print("Received GET todos request")

	todos, err := common.ListItems(ctx)
	if err != nil {
		return serverError(err)
	}

	json, err := json.Marshal(todos)
	if err != nil {
		return serverError(err)
	}
	log.Printf("Successfully fetched todos: %s", json)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(json),
	}, nil
}

func processGet(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, ok := req.PathParameters["id"]

	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string("Invalid Request, pls specify an id"),
		}, nil
	} else {
		return processGetPassport(ctx, id)
	}
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Received req %#v", req)

	return processGet(ctx, req)
}

func main() {
	lambda.Start(handler)
}
