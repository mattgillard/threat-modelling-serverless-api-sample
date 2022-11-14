package common

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/smithy-go"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	. "github.com/mattgillard/user-pii-demo/types"
)

var TableName string

var db dynamodb.Client
var km kms.Client

func init() {
	sdkConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	db = *dynamodb.NewFromConfig(sdkConfig)
	km = *kms.NewFromConfig(sdkConfig)
	TableName = os.Getenv("TableName")
}

func GetItem(ctx context.Context, id string, fieldToDecode string) (*User, error) {
	key, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"id": key,
		},
	}

	log.Printf("Calling Dynamodb with input: %v", input)
	result, err := db.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}
	log.Printf("Executed GetItem DynamoDb successfully. Result: %#v", result)

	if result.Item == nil {
		return nil, nil
	}

	user := new(User)
	err = attributevalue.UnmarshalMap(result.Item, user)
	if err != nil {
		return nil, err
	}
	// decode encrypted fields
	err = user.Decode()

	return user, err
}

func ListItems(ctx context.Context) ([]User, error) {
	users := make([]User, 0)
	var token map[string]types.AttributeValue

	for {
		input := &dynamodb.ScanInput{
			TableName:         aws.String(TableName),
			ExclusiveStartKey: token,
		}

		result, err := db.Scan(ctx, input)
		if err != nil {
			return nil, err
		}

		var fetchedUsers []User
		err = attributevalue.UnmarshalListOfMaps(result.Items, &fetchedUsers)
		if err != nil {
			return nil, err
		}

		users = append(users, fetchedUsers...)
		token = result.LastEvaluatedKey
		if token == nil {
			break
		}
	}

	return users, nil
}

func InsertItem(ctx context.Context, createUser CreateUser) (*User, error) {

	user := User{
		Name:     createUser.Name,
		Address:  createUser.Address,
		Status:   false,
		Passport: createUser.Passport,
	}
	log.Printf("user before encode: %v", user)
	err := user.Encode()
	if err != nil {
		return nil, err
	}
	log.Printf("user after encode: %v", user)
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      item,
	}

	res, err := db.PutItem(ctx, input)
	if err != nil {
		return nil, err
	}

	err = attributevalue.UnmarshalMap(res.Attributes, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func DeleteItem(ctx context.Context, id string) (*User, error) {
	key, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"id": key,
		},
		ReturnValues: types.ReturnValue(*aws.String("ALL_OLD")),
	}

	res, err := db.DeleteItem(ctx, input)
	if err != nil {
		return nil, err
	}

	if res.Attributes == nil {
		return nil, nil
	}

	user := new(User)
	err = attributevalue.UnmarshalMap(res.Attributes, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateItem(ctx context.Context, id string, updateUser UpdateUser) (*User, error) {
	key, err := attributevalue.Marshal(id)
	if err != nil {
		return nil, err
	}

	expr, err := expression.NewBuilder().WithUpdate(
		expression.Set(
			expression.Name("name"),
			expression.Value(updateUser.Name),
		).Set(
			expression.Name("address"),
			expression.Value(updateUser.Address),
		).Set(
			expression.Name("status"),
			expression.Value(updateUser.Status),
		),
	).WithCondition(
		expression.Equal(
			expression.Name("id"),
			expression.Value(id),
		),
	).Build()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"id": key,
		},
		TableName:                 aws.String(TableName),
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
		ReturnValues:              types.ReturnValue(*aws.String("ALL_NEW")),
	}

	res, err := db.UpdateItem(ctx, input)
	if err != nil {
		var smErr *smithy.OperationError
		if errors.As(err, &smErr) {
			var condCheckFailed *types.ConditionalCheckFailedException
			if errors.As(err, &condCheckFailed) {
				return nil, nil
			}
		}

		return nil, err
	}

	if res.Attributes == nil {
		return nil, nil
	}

	user := new(User)
	err = attributevalue.UnmarshalMap(res.Attributes, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
