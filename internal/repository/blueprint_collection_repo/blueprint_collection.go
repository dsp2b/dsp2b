package blueprint_collection_repo

import (
	"context"

	"github.com/codfrm/cago/database/mongo"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/dsp2b/dsp2b-go/internal/model/entity/blueprint_collection_entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlueprintCollectionRepo interface {
	Find(ctx context.Context, id primitive.ObjectID) (*blueprint_collection_entity.BlueprintCollection, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*blueprint_collection_entity.BlueprintCollection, int64, error)
	Create(ctx context.Context, blueprintCollection *blueprint_collection_entity.BlueprintCollection) error
	Update(ctx context.Context, blueprintCollection *blueprint_collection_entity.BlueprintCollection) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByCollection(ctx context.Context, id primitive.ObjectID) ([]*blueprint_collection_entity.BlueprintCollection, error)
}

var defaultBlueprintCollection BlueprintCollectionRepo

func BlueprintCollection() BlueprintCollectionRepo {
	return defaultBlueprintCollection
}

func RegisterBlueprintCollection(i BlueprintCollectionRepo) {
	defaultBlueprintCollection = i
}

type blueprintCollectionRepo struct {
}

func NewBlueprintCollection() BlueprintCollectionRepo {
	return &blueprintCollectionRepo{}
}

func (u *blueprintCollectionRepo) Find(ctx context.Context, id primitive.ObjectID) (*blueprint_collection_entity.BlueprintCollection, error) {
	blueprintCollection := &blueprint_collection_entity.BlueprintCollection{}
	err := mongo.Ctx(ctx).Collection(blueprintCollection.CollectionName()).FindOne(bson.M{
		"_id":    id,
		"status": consts.ACTIVE,
	}).Decode(blueprintCollection)
	if err != nil {
		if mongo.IsNoDocuments(err) {
			return nil, nil
		}
		return nil, err
	}
	return blueprintCollection, nil
}

func (u *blueprintCollectionRepo) Create(ctx context.Context, blueprintCollection *blueprint_collection_entity.BlueprintCollection) error {
	blueprintCollection.ID = primitive.NewObjectID()
	_, err := mongo.Ctx(ctx).Collection(blueprintCollection.CollectionName()).InsertOne(blueprintCollection)
	return err
}

func (u *blueprintCollectionRepo) Update(ctx context.Context, blueprintCollection *blueprint_collection_entity.BlueprintCollection) error {
	_, err := mongo.Ctx(ctx).Collection(blueprintCollection.CollectionName()).UpdateOne(bson.M{
		"_id":    blueprintCollection.ID,
		"status": consts.ACTIVE,
	}, bson.M{
		"$set": blueprintCollection,
	})
	return err
}

func (u *blueprintCollectionRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	blueprintCollection := &blueprint_collection_entity.BlueprintCollection{}
	_, err := mongo.Ctx(ctx).Collection(blueprintCollection.CollectionName()).
		UpdateOne(bson.M{
			"_id":    id,
			"status": consts.ACTIVE,
		}, bson.M{
			"$set": bson.M{
				"status": consts.DELETE,
			},
		})
	return err
}

func (u *blueprintCollectionRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*blueprint_collection_entity.BlueprintCollection, int64, error) {
	blueprintCollection := blueprint_collection_entity.BlueprintCollection{}
	filter := bson.M{
		"status": consts.ACTIVE,
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64(page.GetOffset()))
	findOptions.SetLimit(int64(page.GetSize()))
	findOptions.SetSort(bson.M{"createtime": -1})

	total, err := mongo.Ctx(ctx).Collection(blueprintCollection.CollectionName()).CountDocuments(filter)
	if err != nil {
		return nil, 0, err
	}

	curs, err := mongo.Ctx(ctx).Collection(blueprintCollection.CollectionName()).Find(filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	userList := make([]*blueprint_collection_entity.BlueprintCollection, 0, page.GetSize())
	if err = curs.All(ctx, &userList); err != nil {
		return nil, 0, err
	}

	return userList, total, nil
}

func (u *blueprintCollectionRepo) FindByCollection(ctx context.Context, id primitive.ObjectID) ([]*blueprint_collection_entity.BlueprintCollection, error) {
	blueprintCollection := blueprint_collection_entity.BlueprintCollection{}
	filter := bson.M{
		"collection_id": id,
		"status":        consts.ACTIVE,
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"createtime": -1})

	curs, err := mongo.Ctx(ctx).Collection(blueprintCollection.CollectionName()).Find(filter, findOptions)
	if err != nil {
		return nil, err
	}
	list := make([]*blueprint_collection_entity.BlueprintCollection, 0)
	if err = curs.All(ctx, &list); err != nil {
		return nil, err
	}

	return list, nil
}
