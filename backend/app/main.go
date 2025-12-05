package app

import (
	"fmt"
	"rustdesk-api-server-pro/app/middleware"
	"rustdesk-api-server-pro/app/model"
	"rustdesk-api-server-pro/config"
	"rustdesk-api-server-pro/db"

	"github.com/kataras/iris/v12"
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
