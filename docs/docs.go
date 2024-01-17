// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2024-01-17 09:57:40.646016 +0800 CST m=+0.105536126
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/blueprint": {
            "get": {
                "description": "蓝图列表",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "blueprint"
                ],
                "summary": "蓝图列表",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "properties": {
                                "code": {
                                    "type": "integer"
                                },
                                "data": {
                                    "$ref": "#/definitions/blueprint.ListResponse"
                                },
                                "msg": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/BadRequest"
                        }
                    }
                }
            }
        },
        "/blueprint/parse": {
            "post": {
                "description": "蓝图解析",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "blueprint"
                ],
                "summary": "蓝图解析",
                "parameters": [
                    {
                        "name": "body",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/blueprint.ParseRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "properties": {
                                "code": {
                                    "type": "integer"
                                },
                                "data": {
                                    "$ref": "#/definitions/blueprint.ParseResponse"
                                },
                                "msg": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/BadRequest"
                        }
                    }
                }
            }
        },
        "/blueprint/recipe_panel": {
            "get": {
                "description": "获取配方面板",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "blueprint"
                ],
                "summary": "获取配方面板",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "properties": {
                                "code": {
                                    "type": "integer"
                                },
                                "data": {
                                    "$ref": "#/definitions/blueprint.GetRecipePanelResponse"
                                },
                                "msg": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/BadRequest"
                        }
                    }
                }
            }
        },
        "/blueprint/{id}": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "blueprint"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "properties": {
                                "code": {
                                    "type": "integer"
                                },
                                "data": {
                                    "$ref": "#/definitions/blueprint.DetailResponse"
                                },
                                "msg": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/BadRequest"
                        }
                    }
                }
            }
        },
        "/collection/:id": {
            "get": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "collection"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "name": "ID",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "properties": {
                                "code": {
                                    "type": "integer"
                                },
                                "data": {
                                    "$ref": "#/definitions/collection.DetailResponse"
                                },
                                "msg": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/BadRequest"
                        }
                    }
                }
            }
        },
        "/collection/:id/blueprint": {
            "get": {
                "description": "查询蓝图",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "collection"
                ],
                "summary": "查询蓝图",
                "parameters": [
                    {
                        "type": "string",
                        "name": "ID",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "properties": {
                                "code": {
                                    "type": "integer"
                                },
                                "data": {
                                    "$ref": "#/definitions/collection.GetCollectionBlueprintResponse"
                                },
                                "msg": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/BadRequest"
                        }
                    }
                }
            }
        },
        "/collection/:id/download": {
            "get": {
                "description": "下载蓝图zip包",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "collection"
                ],
                "summary": "下载蓝图zip包",
                "parameters": [
                    {
                        "type": "string",
                        "name": "ID",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "properties": {
                                "code": {
                                    "type": "integer"
                                },
                                "data": {
                                    "$ref": "#/definitions/collection.DownloadResponse"
                                },
                                "msg": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/BadRequest"
                        }
                    }
                }
            }
        },
        "/collection/:id/sub": {
            "get": {
                "description": "查询子蓝图集",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "collection"
                ],
                "summary": "查询子蓝图集",
                "parameters": [
                    {
                        "type": "string",
                        "name": "ID",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "properties": {
                                "code": {
                                    "type": "integer"
                                },
                                "data": {
                                    "$ref": "#/definitions/collection.SubCollectionResponse"
                                },
                                "msg": {
                                    "type": "string"
                                }
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/BadRequest"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "BadRequest": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "错误码",
                    "type": "integer",
                    "format": "int32"
                },
                "msg": {
                    "description": "错误信息",
                    "type": "string"
                }
            }
        },
        "blueprint.Blueprint": {
            "type": "object",
            "properties": {
                "blueprint": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "blueprint.Building": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "icon_path": {
                    "type": "string"
                },
                "item_id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "blueprint.DetailResponse": {
            "type": "object"
        },
        "blueprint.GetRecipePanelResponse": {
            "type": "object"
        },
        "blueprint.Item": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "blueprint.ListResponse": {
            "type": "object",
            "properties": {
                "list": {
                    "type": "array",
                    "items": {
                        "type": "object",
                        "$ref": "#/definitions/blueprint.Item"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "blueprint.ParseRequest": {
            "type": "object",
            "properties": {
                "blueprint": {
                    "type": "string"
                }
            }
        },
        "blueprint.ParseResponse": {
            "type": "object",
            "properties": {
                "blueprint": {
                    "$ref": "#/definitions/blueprint.Blueprint"
                },
                "buildings": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/blueprint.Building"
                    }
                },
                "products": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/blueprint.Product"
                    }
                }
            }
        },
        "blueprint.Product": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "number"
                },
                "icon_path": {
                    "type": "string"
                },
                "item_id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "blueprint.RecipePanel": {
            "type": "object",
            "properties": {
                "building_panel": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "$ref": "#/definitions/blueprint.RecipePanelItem"
                        }
                    }
                },
                "thing_panel": {
                    "type": "array",
                    "items": {
                        "type": "array",
                        "items": {
                            "$ref": "#/definitions/blueprint.RecipePanelItem"
                        }
                    }
                }
            }
        },
        "blueprint.RecipePanelItem": {
            "type": "object",
            "properties": {
                "icon_path": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "item_id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "collection.Collection": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "$ref": "#/definitions/primitive.ObjectID"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "collection.DetailResponse": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "id": {
                    "$ref": "#/definitions/primitive.ObjectID"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "collection.DownloadResponse": {
            "type": "object"
        },
        "collection.GetCollectionBlueprintItem": {
            "type": "object",
            "properties": {
                "blueprint": {
                    "type": "string"
                },
                "createtime": {
                    "type": "integer"
                },
                "description": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "updatetime": {
                    "type": "integer"
                }
            }
        },
        "collection.GetCollectionBlueprintResponse": {
            "type": "object",
            "properties": {
                "blueprint": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/collection.GetCollectionBlueprintItem"
                    }
                }
            }
        },
        "collection.SubCollectionResponse": {
            "type": "object",
            "properties": {
                "collection": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/collection.Collection"
                    }
                }
            }
        },
        "httputils.PageRequest": {
            "type": "object",
            "properties": {
                "page": {
                    "description": "Deprecated 请使用方法GetPage",
                    "type": "integer"
                },
                "size": {
                    "description": "Deprecated 请使用方法GetSize",
                    "type": "integer"
                }
            }
        },
        "httputils.PageResponse": {
            "type": "object",
            "properties": {
                "list": {
                    "type": "array",
                    "items": {
                        "type": "any"
                    }
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "primitive.ObjectID": {
            "type": "string"
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "/api/v1",
	Schemes:     []string{},
	Title:       "api文档",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
