package admin

import (
	"rustdesk-api-server-pro/app/model"
	"rustdesk-api-server-pro/config"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"xorm.io/xorm"
)

type basicController struct {
	Ctx iris.Context
	Db  *xorm.Engine
}

func (c *basicController) GetUser() *model.User {
	return c.Ctx.Values().Get(config.AdminUserKey).(*model.User)
}

func (c *basicController) GetToken() string {
	return c.Ctx.Values().Get(config.AdminAuthTokenString).(string)
}

func (c *basicController) GetAuthToken() *model.AuthToken {
	return c.Ctx.Values().Get(config.AdminAuthToken).(*model.AuthToken)
}

// CheckPermission - Check if user has required role level
func (c *basicController) CheckPermission(requiredRole int) bool {
	user := c.GetUser()
	if user == nil {
		return false
	}
	return user.Role >= requiredRole
}

// RequirePermission - Return error if user doesn't have required role level
func (c *basicController) RequirePermission(requiredRole int, action string) mvc.Result {
	if !c.CheckPermission(requiredRole) {
		return c.Error(nil, "Permission denied: "+action+" requires higher role level")
	}
	return nil
}

func (c *basicController) Success(data interface{}, message string) mvc.Result {
	return c.response(200, data, message)
}

func (c *basicController) Error(data interface{}, message string) mvc.Result {
	return c.response(500, data, message)
}

func (c *basicController) response(code int, data interface{}, message string) mvc.Result {
	return mvc.Response{
		Object: iris.Map{
			"code":    code,
			"message": message,
			"data":    data,
		},
	}
}
