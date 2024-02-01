package utils

import (
	"context"

	"github.com/dsp2b/dsp2b-go/internal/repository/collection_repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RootCollection 获取根蓝图集
func RootCollection(ctx context.Context, id primitive.ObjectID) (primitive.ObjectID, error) {
	uniqueMap := make(map[primitive.ObjectID]struct{})
	parentId := id
	for {
		collection, err := collection_repo.Colletcion().Find(ctx, parentId)
		if err != nil {
			return primitive.NilObjectID, err
		}
		if collection.ParentID.IsZero() {
			return parentId, nil
		}
		if _, ok := uniqueMap[collection.ParentID]; ok {
			return primitive.NilObjectID, nil
		}
		uniqueMap[collection.ParentID] = struct{}{}
		parentId = collection.ParentID
	}
}
