package blueprint_entity

import (
	"encoding/json"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IconInfo struct {
	ItemID   int32
	Name     string
	IconPath string
}

type Icons struct {
	Layout int
	Icon0  *IconInfo
	Icon1  *IconInfo
	Icon2  *IconInfo
	Icon3  *IconInfo
	Icon4  *IconInfo
	Icon5  *IconInfo
}

type Buildings struct {
	ItemId   int    `json:"item_id"`
	Name     string `json:"name"`
	IconPath string `json:"icon_path"`
	Count    int    `json:"count"`
}

type Blueprint struct {
	ID          primitive.ObjectID `bson:"_id"`
	UserID      primitive.ObjectID `bson:"user_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Blueprint   string             `bson:"blueprint"`
	TagsID      []int              `bson:"tags_id"`
	PicList     []string           `bson:"pic_list"`
	GameVersion string             `bson:"game_version"`
	Buildings   string             `bson:"buildings"`
	Icons       string             `bson:"icons"`
	CopyCount   int64              `bson:"copy_count"`
	Original    int8               `bson:"original"`
	Status      int8               `bson:"status"`
	Createtime  time.Time          `bson:"createtime"`
	Updatetime  time.Time          `bson:"updatetime,omitempty"`
}

func (b *Blueprint) CollectionName() string {
	return "blueprint"
}

func (b *Blueprint) SetIcons(icons Icons) error {
	data, err := json.Marshal(icons)
	if err != nil {
		return err
	}
	b.Icons = string(data)
	return nil
}
