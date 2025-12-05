package admin

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type IndexController struct {
	basicController
}

func (c *IndexController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/userinfo", "HandleUserInfo")
}

func (c *IndexController) HandleUserInfo() mvc.Result {
	user := c.GetUser()
	return c.Success(iris.Map{
		"userId":   user.Id,
		"userName": user.Name,
		"role":     user.Role, // Add role for frontend permission checks
	}, "ok")
}
