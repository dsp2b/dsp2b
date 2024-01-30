package collection_repo

import (
	"context"

	"github.com/codfrm/cago/database/mongo"
	"github.com/codfrm/cago/pkg/consts"
	"github.com/codfrm/cago/pkg/utils/httputils"
	"github.com/dsp2b/dsp2b-go/internal/model/entity/collection_entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ColletcionRepo interface {
	Find(ctx context.Context, id primitive.ObjectID) (*collection_entity.Collection, error)
	FindPage(ctx context.Context, page httputils.PageRequest) ([]*collection_entity.Collection, int64, error)
	Create(ctx context.Context, colletcion *collection_entity.Collection) error
	Update(ctx context.Context, colletcion *collection_entity.Collection) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByParent(ctx context.Context, id primitive.ObjectID) ([]*collection_entity.Collection, error)
	UpdateDownloadFile(ctx context.Context, id primitive.ObjectID, downloadFile string) error
}

var defaultColletcion ColletcionRepo

func Colletcion() ColletcionRepo {
	return defaultColletcion
}

func RegisterColletcion(i ColletcionRepo) {
	defaultColletcion = i
}

type colletcionRepo struct {
}

func NewColletcion() ColletcionRepo {
	return &colletcionRepo{}
}

func (u *colletcionRepo) Find(ctx context.Context, id primitive.ObjectID) (*collection_entity.Collection, error) {
	colletcion := &collection_entity.Collection{}
	err := mongo.Ctx(ctx).Collection(colletcion.CollectionName()).FindOne(bson.M{
		"_id":    id,
		"status": consts.ACTIVE,
	}).Decode(colletcion)
	if err != nil {
		if mongo.IsNoDocuments(err) {
			return nil, nil
		}
		return nil, err
	}
	return colletcion, nil
}

func (u *colletcionRepo) Create(ctx context.Context, colletcion *collection_entity.Collection) error {
	colletcion.ID = primitive.NewObjectID()
	_, err := mongo.Ctx(ctx).Collection(colletcion.CollectionName()).InsertOne(colletcion)
	return err
}

func (u *colletcionRepo) Update(ctx context.Context, colletcion *collection_entity.Collection) error {
	_, err := mongo.Ctx(ctx).Collection(colletcion.CollectionName()).UpdateOne(bson.M{
		"_id":    colletcion.ID,
		"status": consts.ACTIVE,
	}, bson.M{
		"$set": colletcion,
	})
	return err
}

func (u *colletcionRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	colletcion := &collection_entity.Collection{}
	_, err := mongo.Ctx(ctx).Collection(colletcion.CollectionName()).
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

func (u *colletcionRepo) FindPage(ctx context.Context, page httputils.PageRequest) ([]*collection_entity.Collection, int64, error) {
	colletcion := collection_entity.Collection{}
	filter := bson.M{
		"status": consts.ACTIVE,
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64(page.GetOffset()))
	findOptions.SetLimit(int64(page.GetSize()))
	findOptions.SetSort(bson.M{"createtime": -1})

	total, err := mongo.Ctx(ctx).Collection(colletcion.CollectionName()).CountDocuments(filter)
	if err != nil {
		return nil, 0, err
	}

	curs, err := mongo.Ctx(ctx).Collection(colletcion.CollectionName()).Find(filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	userList := make([]*collection_entity.Collection, 0, page.GetSize())
	if err = curs.All(ctx, &userList); err != nil {
		return nil, 0, err
	}

	return userList, total, nil
}

func (u *colletcionRepo) FindByParent(ctx context.Context, id primitive.ObjectID) ([]*collection_entity.Collection, error) {
	collection := collection_entity.Collection{}
	filter := bson.M{
		"status":    consts.ACTIVE,
		"parent_id": id,
	}
	curs, err := mongo.Ctx(ctx).Collection(collection.CollectionName()).Find(filter)
	if err != nil {
		return nil, err
	}
	list := make([]*collection_entity.Collection, 0)
	if err = curs.All(ctx, &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (u *colletcionRepo) UpdateDownloadFile(ctx context.Context, id primitive.ObjectID, downloadFile string) error {
	collection := collection_entity.Collection{}
	_, err := mongo.Ctx(ctx).Collection(collection.CollectionName()).UpdateOne(bson.M{
		"_id":    id,
		"status": consts.ACTIVE,
	}, bson.M{
		"$set": bson.M{
			"download_file": downloadFile,
		},
	})
	return err
}
