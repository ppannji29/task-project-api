package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureTaskIndexes(taskCol *mongo.Collection) {
	indexes := []mongo.IndexModel{
		{Keys: map[string]interface{}{"user_id": 1}},
		{Keys: map[string]interface{}{"status": 1}},
		{Keys: map[string]interface{}{"priority": 1}},
		{Keys: map[string]interface{}{"due_date": 1}},
		{Keys: map[string]interface{}{"created_at": -1}},
	}

	for _, idx := range indexes {
		_, err := taskCol.Indexes().CreateOne(context.Background(), idx, options.CreateIndexes().SetMaxTime(0))
		if err != nil {
			log.Printf("⚠️ Failed to create index: %v", err)
		} else {
			log.Printf("✅ Index ensured: %v", idx.Keys)
		}
	}
}
