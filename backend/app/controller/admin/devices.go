package admin

import (
	"rustdesk-api-server-pro/app/model"
	"rustdesk-api-server-pro/config"
	"rustdesk-api-server-pro/db"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"xorm.io/xorm"
)

type DevicesController struct {
	basicController
}

func (c *DevicesController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/devices/list", "HandleList")
	b.Handle("GET", "/devices/online", "HandleOnlineList")
}

func (c *DevicesController) HandleList() mvc.Result {
	currentPage := c.Ctx.URLParamIntDefault("current", 1)
	pageSize := c.Ctx.URLParamIntDefault("size", 10)
	hostname := c.Ctx.URLParamDefault("hostname", "")
	username := c.Ctx.URLParamDefault("username", "")
	rustdesk_id := c.Ctx.URLParamDefault("rustdesk_id", "")
	status := c.Ctx.URLParamDefault("status", "") // "online", "offline", ""
	query := func() *xorm.Session {
		q := c.Db.Table(&model.Device{})

		if hostname != "" {
			q.Where("hostname LIKE ?", "%"+hostname+"%")
		}
		if username != "" {
			q.Where("username LIKE ?", "%"+username+"%")
		}
		if rustdesk_id != "" {
			q.Where("rustdesk_id LIKE ?", "%"+rustdesk_id+"%")
		}
		if status == "online" {
			q.Where("is_online = ?", true)
		} else if status == "offline" {
			q.Where("is_online = ?", false)
		}
		q.Desc("is_online").Asc("username")
		return q
	}

	pagination := db.NewPagination(currentPage, pageSize)
	deviceList := make([]model.Device, 0)

	err := pagination.Paginate(query, &model.Audit{}, &deviceList)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	list := make([]iris.Map, 0)
	for _, a := range deviceList {
		list = append(list, iris.Map{
			"id":           a.Id,
			"rustdesk_id":  a.RustdeskId,
			"hostname":     a.Hostname,
			"username":     a.Username,
			"uuid":         a.Uuid,
			"version":      a.Version,
			"os":           a.Os,
			"memory":       a.Memory,
			"is_online":    a.IsOnline,
			"last_seen_at": a.LastSeenAt.Format(config.TimeFormat),
			"ip_address":   a.IpAddress,
			"conns":        a.Conns,
			"created_at":   a.CreatedAt.Format(config.TimeFormat),
		})
	}
	return c.Success(iris.Map{
		"total":   pagination.TotalCount,
		"records": list,
		"current": currentPage,
		"size":    pageSize,
	}, "ok")
}

func (c *DevicesController) HandleOnlineList() mvc.Result {
	deviceList := make([]model.Device, 0)
	err := c.Db.Where("is_online = ?", true).OrderBy("last_seen_at DESC").Find(&deviceList)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	list := make([]iris.Map, 0)
	for _, a := range deviceList {
		list = append(list, iris.Map{
			"id":           a.Id,
			"rustdesk_id":  a.RustdeskId,
			"hostname":     a.Hostname,
			"username":     a.Username,
			"uuid":         a.Uuid,
			"version":      a.Version,
			"os":           a.Os,
			"memory":       a.Memory,
			"is_online":    a.IsOnline,
			"last_seen_at": a.LastSeenAt.Format(config.TimeFormat),
			"ip_address":   a.IpAddress,
			"conns":        a.Conns,
			"created_at":   a.CreatedAt.Format(config.TimeFormat),
		})
	}

	return c.Success(iris.Map{
		"total":   len(list),
		"records": list,
	}, "ok")
}
