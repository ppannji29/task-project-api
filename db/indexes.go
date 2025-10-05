package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureTaskIndexes(taskCol *mongo.Collection) {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetName("idx_user_id"),
		},
		{
			Keys:    bson.D{{Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_status"),
		},
		{
			Keys:    bson.D{{Key: "priority", Value: 1}},
			Options: options.Index().SetName("idx_priority"),
		},
		{
			Keys:    bson.D{{Key: "due_date", Value: 1}},
			Options: options.Index().SetName("idx_due_date"),
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().SetName("idx_created_at_desc"),
		},
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "status", Value: 1}},
			Options: options.Index().SetName("idx_user_status"),
		},
	}

	for _, idx := range indexes {
		_, err := taskCol.Indexes().CreateOne(context.Background(), idx)
		if err != nil {
			log.Printf("⚠️ Failed to create index: %v", err)
		} else {
			log.Printf("✅ Index ensured: %v", idx.Options.Name)
		}
	}
}

func EnsureUserIndexes(userCol *mongo.Collection) {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_email_unique"),
		},
	}

	for _, idx := range indexes {
		_, err := userCol.Indexes().CreateOne(context.Background(), idx)
		if err != nil {
			log.Printf("⚠️ Failed to create index: %v", err)
		} else {
			log.Printf("✅ Index ensured: %v", idx.Options.Name)
		}
	}
}

// func EnsureTaskIndexes(taskCol *mongo.Collection) {
// 	indexes := []mongo.IndexModel{
// 		{Keys: map[string]interface{}{"user_id": 1}},
// 		{Keys: map[string]interface{}{"status": 1}},
// 		{Keys: map[string]interface{}{"priority": 1}},
// 		{Keys: map[string]interface{}{"due_date": 1}},
// 		{Keys: map[string]interface{}{"created_at": -1}},
// 	}

// 	for _, idx := range indexes {
// 		_, err := taskCol.Indexes().CreateOne(context.Background(), idx, options.CreateIndexes().SetMaxTime(0))
// 		if err != nil {
// 			log.Printf("⚠️ Failed to create index: %v", err)
// 		} else {
// 			log.Printf("✅ Index ensured: %v", idx.Keys)
// 		}
// 	}
// }
