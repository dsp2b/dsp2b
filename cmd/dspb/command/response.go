package command

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Collection struct {
	ID           primitive.ObjectID `json:"id"`
	UserId       string             `json:"user_id"`
	Title        string             `json:"title"`
	Description  string             `json:"description"`
	ParentId     string             `json:"parent_id"`
	DownloadFile string             `json:"download_file"`
	Public       int                `json:"public"`
	Status       int                `json:"status"`
	Createtime   time.Time          `json:"createtime"`
	Updatetime   time.Time          `json:"updatetime"`
}

type BlueprintItem struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	PicList     []string `json:"pic_list"`
	TagsId      []int    `json:"tags_id"`
	CopyCount   int      `json:"copy_count"`
	User        struct {
		Id       string `json:"id"`
		Username string `json:"username"`
	} `json:"user"`
	Tags []struct {
		Id       int    `json:"id"`
		Name     string `json:"name"`
		IconPath string `json:"icon_path"`
	} `json:"tags"`
	CollectionCount int    `json:"collection_count"`
	LikeCount       int    `json:"like_count"`
	Pic             string `json:"pic,omitempty"`
}

type CollectionDetailResponse struct {
	Collection struct {
		Id           string    `json:"id"`
		UserId       string    `json:"user_id"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		ParentId     string    `json:"parent_id"`
		DownloadFile string    `json:"download_file"`
		Public       int       `json:"public"`
		Status       int       `json:"status"`
		Createtime   time.Time `json:"createtime"`
		Updatetime   time.Time `json:"updatetime"`
		User         struct {
			Id       string `json:"id"`
			Username string `json:"username"`
		} `json:"user"`
	} `json:"collection"`
	Self          bool          `json:"self"`
	SubCollection []*Collection `json:"sub_collection"`
	LikeCount     int           `json:"like_count"`
	IsLike        bool          `json:"is_like"`
	I18N          struct {
		Collection string `json:"collection"`
	} `json:"i18n"`
	Href        string           `json:"href"`
	List        []*BlueprintItem `json:"list"`
	Total       int              `json:"total"`
	Sort        string           `json:"sort"`
	Keyword     string           `json:"keyword"`
	View        string           `json:"view"`
	CurrentPage int              `json:"currentPage"`
}
