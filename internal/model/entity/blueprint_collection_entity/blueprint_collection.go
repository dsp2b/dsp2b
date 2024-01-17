package blueprint_collection_entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlueprintCollection struct {
	ID           primitive.ObjectID `bson:"_id"`
	BlueprintID  primitive.ObjectID `bson:"blueprint_id"`
	CollectionID primitive.ObjectID `bson:"collection_id"`
	Status       int8               `bson:"status"`
	Createtime   time.Time          `bson:"createtime"`
	Updatetime   time.Time          `bson:"updatetime,omitempty"`
}

func (b *BlueprintCollection) CollectionName() string {
	return "blueprint_collection"
}
