package service

import (
	"context"
	"fmt"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"github.com/nats-io/go-nats"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type ConstraintsCheck struct {
	Collection string
	Id         string
	Document   M
}

type ReadRequest struct {
	Database       string
	Collection     string
	Id             string
	IncludeDeleted bool
}

type MongoInsertRequest struct {
	Database   string
	Collection string
	Document   M
}

func init() {
	subscribe("mongo.insert", mongoInsert)
}

func mongoInsert(r uniform.IRequest, p diary.IPage) {
	var request MongoInsertRequest
	r.Read(&request)
	x := r.Context()

	p.Notice("", diary.M{
		"request": request,
		"context": x,
	})

	c, cancel := connect(testMode, mongoUri, request.Database, request.Collection)
	defer cancel()

	conn := r.Conn()

	if v, exists := x["field-tags"]; exists && v != nil {
		if _, exists := x["no-constraints"]; !exists {
			if err := conn.Request(p, "constraints.check", r.Remainder(), uniform.Request{
				Model: ConstraintsCheck{
					Collection: request.Collection,
					Document:   request.Document,
				},
				Context: r.Context(),
				Parameters: r.Parameters(),
			}, func(r uniform.IRequest, p diary.IPage) {
				r.Read(&request.Document)
			}); err != nil {
				if err == nats.ErrTimeout {
					p.Warning("constraints", "failed to check constraints", diary.M{
						"database": request.Database,
						"collection": request.Collection,
					})
				} else {
					panic(err)
				}
			}
		}
	}

	request.Document["created-at"] = time.Now().UTC()
	request.Document["modified-at"] = time.Now().UTC()
	delete(request.Document, "id")
	delete(request.Document, "deleted-at")

	for key, value := range request.Document {
		datetime, err := time.Parse(time.RFC3339, fmt.Sprintf("%s", value))
		if err == nil {
			request.Document[key] = datetime
		}
	}

	ctx := context.Background()
	res, err := c.InsertOne(ctx, request.Document)
	if err != nil {
		panic(err)
	}

	id := res.InsertedID.(primitive.ObjectID)

	if err := conn.Request(p, "mongo.read", r.Remainder(), uniform.Request{
		Model: ReadRequest{
			Database: request.Database,
			Collection: request.Collection,
			Id: id.Hex(),
		},
		Context: r.Context(),
		Parameters: r.Parameters(),
	}, func(sr uniform.IRequest, p diary.IPage) {
		sr.Read(&request.Document)
		raw := sr.Raw()

		if err := r.Reply(raw); err != nil {
			panic(err)
		}

		// todo: calculate record hash and use it to determine if the record has actually changed to avoid potential infinite loops
		if err := conn.Publish(p, fmt.Sprintf("mongo.%s.%s.inserted", request.Database, request.Collection), raw); err != nil {
			panic(err)
		}
	}); err != nil {
		panic(err)
	}
}