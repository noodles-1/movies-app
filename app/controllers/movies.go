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
	Plot   string `json:"plot"`
	Year   int    `json:"year"`
	Rating int    `json:"rating"`
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

func (c Movies) GetMovie(id int) revel.Result {
	c.Response.ContentType = "application/json"

	res, err := app.DynamoClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: strconv.Itoa(id)},
		},
		TableName: aws.String("movies"),
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

	res, err := app.DynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("movies"),
		Item:      item,
	})
	if err != nil {
		c.Response.SetStatus(http.StatusInternalServerError)
		return c.RenderJSON(fmt.Sprintf("Error creating item: %v", err))
	}

	c.Response.SetStatus(http.StatusOK)
	return c.RenderJSON(res)
}
