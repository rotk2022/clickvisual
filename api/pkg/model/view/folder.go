package view

import (
	"github.com/clickvisual/clickvisual/api/pkg/model/db"
)

type ReqCreateFolder struct {
	Iid     int `json:"iid" form:"iid" binding:"required"`
	Primary int `json:"primary" form:"primary" binding:"required"`
	ReqUpdateFolder
}

type ReqUpdateFolder struct {
	Name     string `json:"name" form:"name" binding:"required"`
	Desc     string `json:"desc" form:"desc"`
	ParentId int    `json:"parentId" form:"parentId"`
}

type ReqListFolder struct {
	Iid     int `json:"iid" form:"iid"  binding:"required"`
	Primary int `json:"primary" form:"primary" binding:"required"`
}

type RespListFolder struct {
	Id       int              `json:"id"`
	Name     string           `json:"name"`
	Desc     string           `json:"desc"`
	ParentId int              `json:"parentId"`
	Children []RespListFolder `json:"children"`
	Nodes    []*db.Node       `json:"nodes"`
}
