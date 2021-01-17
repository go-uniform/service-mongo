package service

import (
	"context"
	"fmt"
	"github.com/go-diary/diary"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	AppClient = "uprate"
	AppProject = "uniform"
)

var mongoUri string = ""

var connect = func(testMode bool, mongoUri, database, collection string) (c *mongo.Collection, close context.CancelFunc) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)

	err = client.Connect(ctx)
	if err != nil {
		cancel()
		panic(err)
	}

	testPrefix := ""
	if testMode {
		testPrefix = "test_"
	}
	c = client.Database(fmt.Sprintf("%s%s", testPrefix, database)).Collection(fmt.Sprintf("%s%s", testPrefix, collection))
	close = func() {
		client.Disconnect(ctx)
	}

	return
}

func Run(p diary.IPage) {
	mongoUri = fmt.Sprint(args["mongo"])
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)

	if err := client.Connect(ctx); err != nil {
		cancel()
		panic(err)
	}
}