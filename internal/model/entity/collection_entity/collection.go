package collection_entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Collection struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	ParentID    primitive.ObjectID `bson:"parent_id"`
	Status      int8               `bson:"status"`
	Createtime  time.Time          `bson:"createtime"`
	Updatetime  time.Time          `bson:"updatetime,omitempty"`
}

func (c *Collection) CollectionName() string {
	return "collection"
}
