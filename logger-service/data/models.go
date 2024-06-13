package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Models struct {
	LogEntry LogEntry
}

func New(mongo *mongo.Client) Models {
	client = mongo
	return Models{
		LogEntry: LogEntry{},
	}
}

type LogEntry struct {
	ID        string `bson:"_id" json:"id"`
	Name      string `bson:"name" json:"name"`
	Data      string `bson:"data" json:"data"`
	CreateAt  string `bson:"create_at" json:"created_at"`
	UpdatedAt string `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("log")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreateAt:  entry.CreateAt,
		UpdatedAt: entry.CreateAt,
	})

	if err != nil {
		return err
	}

	return err

}

func (l *LogEntry) GetAll() ([]*LogEntry, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{})

	cur, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var logs []*LogEntry

	for cur.Next(ctx) {
		var item LogEntry

		err := cur.Decode(&item)
		if err != nil {
			return nil, err
		}

		logs = append(logs, &item)

	}

	return logs, nil
}


// Get by 30
func(l *LogEntry)  GetOne() (*LogEntry, error) {
	return nil, nil 
}

func(l LogEntry) DropCollection() error {
	return nil
}


