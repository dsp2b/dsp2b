package blueprint_repo

import (
	"context"

	"github.com/codfrm/cago/database/mongo"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/dsp2b/dsp2b-go/internal/model/entity/blueprint_entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlueprintRepo interface {
	Find(ctx context.Context, id primitive.ObjectID) (*blueprint_entity.Blueprint, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*blueprint_entity.Blueprint, int64, error)
	Create(ctx context.Context, blueprint *blueprint_entity.Blueprint) error
	Update(ctx context.Context, blueprint *blueprint_entity.Blueprint) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

var defaultBlueprint BlueprintRepo

func Blueprint() BlueprintRepo {
	return defaultBlueprint
}

func RegisterBlueprint(i BlueprintRepo) {
	defaultBlueprint = i
}

type blueprintRepo struct {
}

func NewBlueprint() BlueprintRepo {
	return &blueprintRepo{}
}

func (u *blueprintRepo) Find(ctx context.Context, id primitive.ObjectID) (*blueprint_entity.Blueprint, error) {
	blueprint := &blueprint_entity.Blueprint{}
	err := mongo.Ctx(ctx).Collection(blueprint.CollectionName()).FindOne(bson.M{
		"_id":    id,
		"status": consts.ACTIVE,
	}).Decode(blueprint)
	if err != nil {
		if mongo.IsNoDocuments(err) {
			return nil, nil
		}
		return nil, err
	}
	return blueprint, nil
}

func (u *blueprintRepo) Create(ctx context.Context, blueprint *blueprint_entity.Blueprint) error {
	blueprint.ID = primitive.NewObjectID()
	_, err := mongo.Ctx(ctx).Collection(blueprint.CollectionName()).InsertOne(blueprint)
	return err
}

func (u *blueprintRepo) Update(ctx context.Context, blueprint *blueprint_entity.Blueprint) error {
	_, err := mongo.Ctx(ctx).Collection(blueprint.CollectionName()).UpdateOne(bson.M{
		"_id":    blueprint.ID,
		"status": consts.ACTIVE,
	}, bson.M{
		"$set": blueprint,
	})
	return err
}

func (u *blueprintRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	blueprint := &blueprint_entity.Blueprint{}
	_, err := mongo.Ctx(ctx).Collection(blueprint.CollectionName()).
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

func (u *blueprintRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*blueprint_entity.Blueprint, int64, error) {
	blueprint := blueprint_entity.Blueprint{}
	filter := bson.M{
		"status": consts.ACTIVE,
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64(page.GetOffset()))
	findOptions.SetLimit(int64(page.GetSize()))
	findOptions.SetSort(bson.M{"createtime": -1})

	total, err := mongo.Ctx(ctx).Collection(blueprint.CollectionName()).CountDocuments(filter)
	if err != nil {
		return nil, 0, err
	}

	curs, err := mongo.Ctx(ctx).Collection(blueprint.CollectionName()).Find(filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	blueprintList := make([]*blueprint_entity.Blueprint, 0, page.GetSize())
	if err = curs.All(ctx, &blueprintList); err != nil {
		return nil, 0, err
	}

	return blueprintList, total, nil
}
