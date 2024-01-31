package command

import (
	"path/filepath"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blueprint struct {
	ID              primitive.ObjectID `json:"id"`               // 蓝图id
	Title           string             `json:"title"`            // 蓝图标题
	Description     string             `json:"description"`      // 蓝图描述
	Hash            string             `json:"hash"`             // 蓝图代码hash
	Updatetime      int64              `json:"updatetime"`       // 更新时间
	LocalUpdatetime int64              `json:"local_updatetime"` // 本地更新时间
}

type Repository struct {
	ID          primitive.ObjectID `json:"id"`                   // 仓库id
	Title       string             `json:"title"`                // 蓝图仓库名
	Description string             `json:"description"`          // 蓝图仓库描述
	Version     string             `json:"version,omitempty"`    // 蓝图仓库版本
	Blueprint   []*Blueprint       `json:"blueprint,omitempty"`  // 蓝图
	Repository  []*Repository      `json:"repository,omitempty"` // 子仓库
}

func (r *Repository) BlueprintMap() map[string]*Blueprint {
	ret := make(map[string]*Blueprint, 0)
	for _, v := range r.Blueprint {
		ret[v.Title] = v
	}
	for _, v := range r.Repository {
		for _, v2 := range v.BlueprintMap() {
			ret[filepath.Join(v.Title, v2.Title)] = v2
		}
	}
	return ret
}

func (r *Repository) RepositoryMap() map[string]*Repository {
	ret := make(map[string]*Repository, 0)
	for _, v := range r.Repository {
		ret[v.Title] = v
	}
	for _, v := range r.Repository {
		for _, v2 := range v.RepositoryMap() {
			ret[filepath.Join(v.Title, v2.Title)] = v2
		}
	}
	return ret
}
