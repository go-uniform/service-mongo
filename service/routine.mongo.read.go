package service

import (
	"context"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoReadRequest struct {
	Database       string
	Collection     string
	Id             string
	IncludeDeleted bool
}

func init() {
	subscribe("mongo.read", mongoRead)
}

func mongoRead(r uniform.IRequest, p diary.IPage) {
	var request MongoReadRequest
	r.Read(&request)
	x := r.Context()

	p.Notice("", diary.M{
		"request": request,
		"context": x,
	})

	id, err := primitive.ObjectIDFromHex(request.Id)
	if err != nil {
		panic(err)
	}

	c, cancel := connect(testMode, mongoUri, request.Database, request.Collection)
	defer cancel()

	query := bson.D{{"_id", id}}
	if !request.IncludeDeleted {
		query = append(query, bson.E{
			Key:   "deleted-at",
			Value: nil,
		})
	}

	ctx := context.Background()
	res := c.FindOne(ctx, query)
	if err := res.Err(); err != nil {
		panic(err)
	}

	var responseObject map[string]interface{}
	if err := res.Decode(&responseObject); err != nil {
		panic(err)
	}
	responseObject["id"] = responseObject["_id"]
	delete(responseObject, "_id")

	if err := r.Reply(uniform.Request{
		Model: responseObject,
		Context: r.Context(),
	}); err != nil {
		panic(err)
	}
}
