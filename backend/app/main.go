package app

import (
	"fmt"
	"rustdesk-api-server-pro/app/middleware"
	"rustdesk-api-server-pro/app/model"
	"rustdesk-api-server-pro/config"
	"rustdesk-api-server-pro/db"

	"github.com/kataras/iris/v12"
	"xorm.io/xorm"
)

func newApp(cfg *config.ServerConfig) (*iris.Application, error) {
	app := iris.Default()

	dbEngine, err := db.NewEngine(cfg.Db)
	if err != nil {
		app.Logger().Fatal("Db Engine create error:", err)
		return nil, err
	}

	// Auto-sync database tables on startup
	app.Logger().Info("Syncing database tables...")
	err = dbEngine.Sync2(
		new(model.User),
		new(model.Device),
		new(model.Peer),
		new(model.AddressBook),
		new(model.AddressBookTag),
		new(model.Tags),
		new(model.AuthToken),
		new(model.Audit),
		new(model.FileTransfer),
		new(model.MailLogs),
		new(model.MailTemplate),
		new(model.SystemSettings),
		new(model.VerifyCode),
		// DocHelp tables
		new(model.KnowledgeBaseCategory),
		new(model.KnowledgeBaseArticle),
		new(model.Ticket),
		new(model.TicketComment),
	)
	if err != nil {
		app.Logger().Fatal("Database sync error:", err)
		return nil, err
	}
	app.Logger().Info("Database tables synced successfully!")

	// Auto-fix user roles after sync (migration for existing users)
	app.Logger().Info("Checking user roles...")
	fixUserRoles(dbEngine, app)

	app.RegisterDependency(dbEngine, cfg)

	app.OnErrorCode(iris.StatusNotFound, func(context iris.Context) {
		requestInfo := fmt.Sprintf("(404)▶ %s:%s", context.Method(), context.Request().RequestURI)
		body, _ := context.GetBody()
		context.Application().Logger().Info(requestInfo)
		for header, value := range context.Request().Header {
			fmt.Println(header+":", value)
		}
		fmt.Println(string(body))
	})

	app.Use(iris.Compression)
	if cfg.HttpConfig.PrintRequestLog {
		app.Use(middleware.RequestLogger())
	}

	SetRoute(app)

	// Serve uploaded files
	app.HandleDir("/uploads", iris.Dir("./uploads"))

	app.HandleDir("/", iris.Dir(cfg.HttpConfig.StaticDir))

	return app, nil
}

// fixUserRoles - Auto-migrate existing users to role-based system
func fixUserRoles(dbEngine *xorm.Engine, app *iris.Application) {
	var users []model.User
	err := dbEngine.Where("role = 0 OR role IS NULL").Find(&users)
	if err != nil {
		app.Logger().Warn("Failed to check user roles:", err)
		return
	}

	if len(users) == 0 {
		app.Logger().Info("All users have valid roles ✓")
		return
	}

	app.Logger().Infof("Found %d users without role, fixing...", len(users))
	
	for _, user := range users {
		// If user has is_admin=true, set as Super Admin (4)
		// Otherwise set as regular User (1)
		newRole := model.ROLE_USER
		if user.IsAdmin {
			newRole = model.ROLE_SUPER_ADMIN
		}

		_, err := dbEngine.ID(user.Id).Cols("role").Update(&model.User{Role: newRole})
		if err != nil {
			app.Logger().Warnf("Failed to update role for user %s: %v", user.Username, err)
		} else {
			roleName := "USER"
			if newRole == model.ROLE_SUPER_ADMIN {
				roleName = "SUPER_ADMIN"
			}
			app.Logger().Infof("✓ User '%s' (ID: %d) → role: %s", user.Username, user.Id, roleName)
		}
	}

	app.Logger().Info("User roles migration completed!")
}

func StartServer() (bool, error) {
	cfg := config.GetServerConfig()

	StartJobs(cfg)

	app, err := newApp(cfg)
	if err != nil {
		return false, err
	}

	err = app.Listen(cfg.HttpConfig.Port, iris.WithoutBodyConsumptionOnUnmarshal)
	if err != nil {
		return false, err
	}

	return true, nil
}
