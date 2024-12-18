package controllers

import (
	"context"
	//"encoding/json"
	"fmt"
	"movies-app/app"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/oklog/ulid/v2"
	"github.com/revel/revel"
)

type Movies struct {
	*revel.Controller
}

type RequestData struct {
	Title  string `json:"title"`
	Plot   string `json:"plot,omitempty"`
	Year   int    `json:"year,omitempty"`
	Rating int    `json:"rating,omitempty"`
}

func (c Movies) Index() revel.Result {
	c.Response.ContentType = "application/json"

	res, err := app.DynamoClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("movies"),
	})
	if err != nil {
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON(fmt.Sprintf("Error scanning movies table : %v", err))
	}

	c.Response.SetStatus(http.StatusOK)
	return c.RenderJSON(res.Items)
}

func (c Movies) GetMovie(id string) revel.Result {
	c.Response.ContentType = "application/json"

	res, err := app.DynamoClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("movies"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON(fmt.Sprintf("Error retrieving item: %v", err))
	}

	c.Response.SetStatus(http.StatusOK)
	return c.RenderJSON(res.Item)
}

func (c Movies) AddMovie() revel.Result {
	c.Response.ContentType = "application/json"

	data := RequestData{}
	c.Params.BindJSON(&data)

	item := map[string]types.AttributeValue{
		"id":     &types.AttributeValueMemberS{Value: ulid.Make().String()},
		"title":  &types.AttributeValueMemberS{Value: data.Title},
		"plot":   &types.AttributeValueMemberS{Value: data.Plot},
		"year":   &types.AttributeValueMemberN{Value: strconv.Itoa(data.Year)},
		"rating": &types.AttributeValueMemberN{Value: strconv.Itoa(data.Rating)},
	}

	_, err := app.DynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("movies"),
		Item:      item,
	})
	if err != nil {
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON(fmt.Sprintf("Error creating item: %v", err))
	}

	c.Response.SetStatus(http.StatusOK)
	return c.RenderJSON(data)
}

func (c Movies) UpdateMovie(id string) revel.Result {
	c.Response.ContentType = "application/json"

	data := RequestData{}
	c.Params.BindJSON(&data)
	updateExpr := "SET"
	exprValues := map[string]types.AttributeValue{}

	updateExpr += " title = :title,"
	exprValues[":title"] = &types.AttributeValueMemberS{Value: data.Title}

	if data.Plot != "" {
		updateExpr += " plot = :plot,"
		exprValues[":plot"] = &types.AttributeValueMemberS{Value: data.Plot}
	}
	if data.Year != 0 {
		updateExpr += " year = :year,"
		exprValues[":year"] = &types.AttributeValueMemberN{Value: strconv.Itoa(data.Year)}
	}
	if data.Rating != 0 {
		updateExpr += " rating = :rating,"
		exprValues[":rating"] = &types.AttributeValueMemberN{Value: strconv.Itoa(data.Rating)}
	}

	updateExpr = updateExpr[:len(updateExpr)-1]

	_, err := app.DynamoClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String("movies"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeValues: exprValues,
		ReturnValues:              types.ReturnValueUpdatedNew,
	})
	if err != nil {
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON(fmt.Sprintf("Error updating item: %v", err))
	}

	c.Response.SetStatus(http.StatusOK)
	return c.RenderJSON(data)
}

func (c Movies) DeleteMovie(id string) revel.Result {
	c.Response.ContentType = "application/json"

	_, err := app.DynamoClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String("movies"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON(fmt.Sprintf("Error deleting item: %v", err))
	}

	c.Response.SetStatus(http.StatusOK)
	return c.RenderJSON("")
}
