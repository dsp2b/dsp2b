package blueprint_entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Blueprint struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Blueprint   string             `bson:"blueprint"`
	Status      int8               `bson:"status"`
	Createtime  time.Time          `bson:"createtime"`
	Updatetime  time.Time          `bson:"updatetime,omitempty"`
}

func (b *Blueprint) CollectionName() string {
	return "blueprint"
}
