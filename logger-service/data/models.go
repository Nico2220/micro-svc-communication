package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type LogEntry struct {
	ID        string `bson:"_id" json:"id"`
	Name      string `bson:"name" json:"name"`
	Data      string `bson:"data" json:"data"`
	CreateAt  string `bson:"create_at" json:"created_at"`
	UpdatedAt string `bson:"updated_at" json:"updated_at"`
}


type Models struct {
	LogEntry LogEntryModel
}

func New(mongo *mongo.Client) Models {
	
	return Models{
		LogEntry: LogEntryModel{mongo},
	}
}


type LogEntryModel struct{
	DB *mongo.Client
}

func (m *LogEntryModel) Insert(entry LogEntry) error {
	collection := m.DB.Database("logs").Collection("log")
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

func (m *LogEntryModel) GetAll() ([]*LogEntry, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := m.DB.Database("logs").Collection("logs")

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
func(m *LogEntryModel)  GetOne() (*LogEntry, error) {
	return nil, nil 
}

func(m LogEntryModel) DropCollection() error {
	return nil
}


