package short

import (
	"github.com/ego-component/egorm"
	"github.com/spf13/cast"

	"github.com/clickvisual/clickvisual/api/internal/invoker"
	"github.com/clickvisual/clickvisual/api/pkg/component/core"
	"github.com/clickvisual/clickvisual/api/pkg/model/db"
	"github.com/clickvisual/clickvisual/api/pkg/model/view"
)

func FolderCreate(c *core.Context) {
	var req view.ReqCreateFolder
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "invalid parameter: "+err.Error(), nil)
		return
	}
	obj := &db.Folder{
		Uid:      c.Uid(),
		Name:     req.Name,
		Desc:     req.Desc,
		ParentId: req.ParentId,
		Primary:  req.Primary,
	}
	err := db.FolderCreate(invoker.Db, obj)
	if err != nil {
		c.JSONE(1, "create failed: "+err.Error(), nil)
		return
	}
	c.JSONOK()
}

func FolderUpdate(c *core.Context) {
	id := cast.ToInt(c.Param("id"))
	if id == 0 {
		c.JSONE(1, "invalid parameter", nil)
		return
	}
	var req view.ReqUpdateFolder
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "invalid parameter: "+err.Error(), nil)
		return
	}
	ups := make(map[string]interface{}, 0)
	ups["name"] = req.Name
	ups["desc"] = req.Desc
	ups["parent_id"] = req.ParentId
	if err := db.FolderUpdate(invoker.Db, id, ups); err != nil {
		c.JSONE(1, "update failed: "+err.Error(), nil)
		return
	}
	c.JSONOK()
}

func FolderList(c *core.Context) {
	var req view.ReqListFolder
	if err := c.Bind(&req); err != nil {
		c.JSONE(1, "invalid parameter: "+err.Error(), nil)
		return
	}
	conds := egorm.Conds{}
	conds["iid"] = req.Iid
	conds["primary"] = req.Primary
	fs, err := db.FolderList(conds)
	if err != nil {
		c.JSONE(core.CodeErr, err.Error(), nil)
		return
	}
	// no folder node
	condsNsNoFolder := egorm.Conds{}
	condsNsNoFolder["folder_id"] = 0
	nsnf, _ := db.NodeList(condsNsNoFolder)
	// root
	res := view.RespListFolder{
		Id:       0,
		Name:     "root",
		Desc:     "",
		ParentId: -1,
		Children: make([]view.RespListFolder, 0),
		Nodes:    nsnf,
	}
	// level 1
	level1children := make(map[int][]view.RespListFolder)
	for _, f := range fs {
		// query nodes
		condsNs := egorm.Conds{}
		condsNs["folder_id"] = f.ID
		ns, _ := db.NodeList(condsNs)
		// build item
		item := view.RespListFolder{
			Id:       f.ID,
			Name:     f.Name,
			Desc:     f.Desc,
			ParentId: f.ParentId,
			Children: make([]view.RespListFolder, 0),
			Nodes:    ns,
		}
		if f.ParentId != 0 {
			level1children[f.ParentId] = append(level1children[f.ParentId], item)
		} else {
			res.Children = append(res.Children, item)
		}
	}
	// level 2
	for index, level1 := range res.Children {
		if l1c, ok := level1children[level1.Id]; ok {
			res.Children[index].Children = append(res.Children[index].Children, l1c...)
		}
	}
	c.JSONE(core.CodeOK, "succ", res)
	return
}

func FolderDelete(c *core.Context) {
	id := cast.ToInt(c.Param("id"))
	if id == 0 {
		c.JSONE(1, "invalid parameter", nil)
		return
	}
	conds := egorm.Conds{}
	conds["folder_id"] = id
	ns, err := db.NodeList(conds)
	if err != nil {
		c.JSONE(1, "failed to delete: "+err.Error(), nil)
		return
	}
	if len(ns) != 0 {
		c.JSONE(1, "failed to delete: u should delete nodes first.", nil)
		return
	}
	if err = db.FolderDelete(invoker.Db, id); err != nil {
		c.JSONE(1, "failed to delete: "+err.Error(), nil)
		return
	}
	c.JSONOK()
}

func FolderInfo(c *core.Context) {
	id := cast.ToInt(c.Param("id"))
	if id == 0 {
		c.JSONE(1, "invalid parameter", nil)
		return
	}
	res, err := db.FolderInfo(invoker.Db, id)
	if err != nil {
		c.JSONE(core.CodeErr, err.Error(), nil)
		return
	}
	c.JSONE(core.CodeOK, "succ", res)
	return
}
