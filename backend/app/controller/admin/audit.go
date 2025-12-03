package admin

import (
	"rustdesk-api-server-pro/app/model"
	"rustdesk-api-server-pro/config"
	"rustdesk-api-server-pro/db"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"xorm.io/xorm"
)

type AuditController struct {
	basicController
}

func (c *AuditController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/audit/list", "HandleList")
	b.Handle("GET", "/audit/file-transfer-list", "HandleFileTransferList")
	b.Handle("GET", "/audit/stats", "HandleStats")
}

func (c *AuditController) HandleList() mvc.Result {
	currentPage := c.Ctx.URLParamIntDefault("current", 1)
	pageSize := c.Ctx.URLParamIntDefault("size", 10)
	conn_id := c.Ctx.URLParamDefault("conn_id", "")
	_type := c.Ctx.URLParamDefault("type", "")
	rustdesk_id := c.Ctx.URLParamDefault("rustdesk_id", "")
	ip := c.Ctx.URLParamDefault("ip", "")
	session_id := c.Ctx.URLParamDefault("session_id", "")
	uuid := c.Ctx.URLParamDefault("uuid", "")
	created_at_0 := c.Ctx.URLParamDefault("created_at[0]", "")
	created_at_1 := c.Ctx.URLParamDefault("created_at[1]", "")
	closed_at_0 := c.Ctx.URLParamDefault("closed_at[0]", "")
	closed_at_1 := c.Ctx.URLParamDefault("closed_at[1]", "")

	query := func() *xorm.Session {
		q := c.Db.Table(&model.Audit{})
		if conn_id != "" {
			q.Where("audit.conn_id = ?", conn_id)
		}
		if _type != "" {
			q.Where("audit.type = ?", _type)
		}
		if rustdesk_id != "" {
			q.Where("audit.rustdesk_id = ?", rustdesk_id)
		}
		if ip != "" {
			q.Where("audit.ip = ?", ip)
		}
		if session_id != "" {
			q.Where("audit.session_id = ?", session_id)
		}
		if uuid != "" {
			q.Where("audit.uuid = ?", uuid)
		}
		if created_at_0 != "" && created_at_1 != "" {
			q.Where("audit.created_at BETWEEN ? AND ?", created_at_0, created_at_1)
		}
		if closed_at_0 != "" && closed_at_1 != "" {
			q.Where("audit.closed_at BETWEEN ? AND ?", closed_at_0, closed_at_1)
		}
		q.Desc("id")
		return q
	}

	pagination := db.NewPagination(currentPage, pageSize)
	auditList := make([]model.Audit, 0)
	err := pagination.Paginate(query, &model.Audit{}, &auditList)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	list := make([]iris.Map, 0)
	for _, a := range auditList {
		// Get username from user_id
		var user model.User
		username := "-"
		if a.UserId > 0 {
			has, err := c.Db.ID(a.UserId).Get(&user)
			if err == nil && has {
				username = user.Username
			}
		}
		
		list = append(list, iris.Map{
			"id":          a.Id,
			"user_id":     a.UserId,
			"username":    username,
			"conn_id":     a.ConnId,
			"rustdesk_id": a.RustdeskId,
			"ip":          a.IP,
			"session_id":  a.SessionId,
			"uuid":        a.Uuid,
			"type":        a.Type,
			"created_at":  a.CreatedAt.Format(config.TimeFormat),
			"closed_at":   a.ClosedAt.Format(config.TimeFormat),
		})
	}
	return c.Success(iris.Map{
		"total":   pagination.TotalCount,
		"records": list,
		"current": currentPage,
		"size":    pageSize,
	}, "ok")
}

func (c *AuditController) HandleFileTransferList() mvc.Result {
	currentPage := c.Ctx.URLParamIntDefault("current", 1)
	pageSize := c.Ctx.URLParamIntDefault("size", 10)
	_type := c.Ctx.URLParamDefault("type", "")
	rustdesk_id := c.Ctx.URLParamDefault("rustdesk_id", "")
	peer_id := c.Ctx.URLParamDefault("peer_id", "")
	uuid := c.Ctx.URLParamDefault("uuid", "")
	created_at_0 := c.Ctx.URLParamDefault("created_at[0]", "")
	created_at_1 := c.Ctx.URLParamDefault("created_at[1]", "")

	query := func() *xorm.Session {
		q := c.Db.Table(&model.FileTransfer{})
		if _type != "" {
			q.Where("type = ?", _type)
		}
		if rustdesk_id != "" {
			q.Where("rustdesk_id = ?", rustdesk_id)
		}
		if peer_id != "" {
			q.Where("peer_id = ?", peer_id)
		}
		if uuid != "" {
			q.Where("audit.uuid = ?", uuid)
		}
		if created_at_0 != "" && created_at_1 != "" {
			q.Where("audit.created_at BETWEEN ? AND ?", created_at_0, created_at_1)
		}
		q.Desc("id")
		return q
	}

	pagination := db.NewPagination(currentPage, pageSize)
	fileTransferList := make([]model.FileTransfer, 0)
	err := pagination.Paginate(query, &model.FileTransfer{}, &fileTransferList)
	if err != nil {
		return c.Error(nil, err.Error())
	}

	list := make([]iris.Map, 0)
	for _, a := range fileTransferList {
		list = append(list, iris.Map{
			"id":          a.Id,
			"rustdesk_id": a.RustdeskId,
			"peer_id":     a.PeerId,
			"path":        a.Path,
			"uuid":        a.Uuid,
			"type":        a.Type,
			"created_at":  a.CreatedAt.Format(config.TimeFormat),
		})
	}
	return c.Success(iris.Map{
		"total":   pagination.TotalCount,
		"records": list,
		"current": currentPage,
		"size":    pageSize,
	}, "ok")
}

// HandleStats returns audit statistics
func (c *AuditController) HandleStats() mvc.Result {
	// Top 10 most accessed devices
	type DeviceStats struct {
		RustdeskId string `xorm:"rustdesk_id"`
		Count      int    `xorm:"count"`
	}
	topDevices := make([]DeviceStats, 0)
	err := c.Db.SQL(`
		SELECT rustdesk_id, COUNT(*) as count 
		FROM audit 
		WHERE rustdesk_id != '' 
		GROUP BY rustdesk_id 
		ORDER BY count DESC 
		LIMIT 10
	`).Find(&topDevices)
	if err != nil {
		c.Ctx.Application().Logger().Errorf("Error getting top devices: %v", err)
	}

	// Format top devices for frontend
	topDevicesList := make([]iris.Map, 0)
	for _, d := range topDevices {
		topDevicesList = append(topDevicesList, iris.Map{
			"rustdesk_id": d.RustdeskId,
			"count":       d.Count,
		})
	}

	// Average session duration (in seconds)
	type AvgDuration struct {
		AvgSeconds float64 `xorm:"avg_seconds"`
	}
	var avgDuration AvgDuration
	_, err = c.Db.SQL(`
		SELECT AVG(CAST((julianday(closed_at) - julianday(created_at)) * 86400 AS INTEGER)) as avg_seconds
		FROM audit 
		WHERE closed_at IS NOT NULL AND closed_at > created_at
	`).Get(&avgDuration)
	if err != nil {
		c.Ctx.Application().Logger().Errorf("Error getting avg duration: %v", err)
		avgDuration.AvgSeconds = 0
	}

	// Connections per day (last 7 days)
	type DailyStats struct {
		Date  string `xorm:"date"`
		Count int    `xorm:"count"`
	}
	dailyStats := make([]DailyStats, 0)
	err = c.Db.SQL(`
		SELECT DATE(created_at) as date, COUNT(*) as count
		FROM audit
		WHERE created_at >= DATE('now', '-7 days')
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`).Find(&dailyStats)
	if err != nil {
		c.Ctx.Application().Logger().Errorf("Error getting daily stats: %v", err)
	}

	// Format daily stats for chart
	dailyStatsList := make([]iris.Map, 0)
	for _, d := range dailyStats {
		dailyStatsList = append(dailyStatsList, iris.Map{
			"date":  d.Date,
			"count": d.Count,
		})
	}

	// Total connections
	totalConnections, _ := c.Db.Count(&model.Audit{})

	// Top 10 most active users
	type UserStats struct {
		UserId   int    `xorm:"user_id"`
		Count    int    `xorm:"count"`
		Username string `xorm:"username"`
	}
	topUsers := make([]UserStats, 0)
	err = c.Db.SQL(`
		SELECT a.user_id, COUNT(*) as count, u.username
		FROM audit a
		LEFT JOIN user u ON a.user_id = u.id
		WHERE a.user_id > 0
		GROUP BY a.user_id, u.username
		ORDER BY count DESC
		LIMIT 10
	`).Find(&topUsers)
	if err != nil {
		c.Ctx.Application().Logger().Errorf("Error getting top users: %v", err)
	}

	// Format top users for frontend
	topUsersList := make([]iris.Map, 0)
	for _, u := range topUsers {
		topUsersList = append(topUsersList, iris.Map{
			"user_id":  u.UserId,
			"username": u.Username,
			"count":    u.Count,
		})
	}

	return c.Success(iris.Map{
		"top_devices":       topDevicesList,
		"top_users":         topUsersList,
		"avg_duration":      int(avgDuration.AvgSeconds),
		"daily_stats":       dailyStatsList,
		"total_connections": totalConnections,
	}, "ok")
}
