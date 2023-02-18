package data

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	mongo_driver "logger-service/mongo-driver"
	"time"
)

var client *mongo.Client

type Models struct {
	LogEntry    LogEntry
	MongoClient mongo_driver.Interactor
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func New(mongo *mongo.Client) Models {
	client = mongo
	return Models{
		LogEntry: LogEntry{},
	}
}

func (m *Models) Insert(ctx context.Context, entry LogEntry) error {
	le := LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := m.MongoClient.Insert(ctx, le); err != nil {
		return fmt.Errorf("failed to replace document in mongo DB %w", err)
	}
	return nil
}

func (m *Models) All(ctx context.Context) ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := m.MongoClient.Find(ctx, bson.D{}, opts)
	if err != nil {
		fmt.Println("Finding all docs error:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry
	for cursor.Next(ctx) {
		var item LogEntry

		err := cursor.Decode(&item)
		if err != nil {
			log.Println("Error decoding into slice:", err)
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}
	return logs, nil
}
