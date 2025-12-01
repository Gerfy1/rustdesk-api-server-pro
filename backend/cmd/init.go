package cmd

import (
	"fmt"
	"rustdesk-api-server-pro/app/model"
	"rustdesk-api-server-pro/config"
	"rustdesk-api-server-pro/db"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize database (create tables)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetServerConfig()
		engine, err := db.NewEngine(cfg.Db)
		if err != nil {
			fmt.Println("Db Engine create error:", err)
			return
		}

		fmt.Println("Syncing database tables...")
		err = engine.Sync2(
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
		)
		if err != nil {
			fmt.Println("Database sync error:", err)
			return
		}

		fmt.Println("✓ Database tables created/updated successfully!")
		fmt.Println("✓ You can now add users with: ./rustdesk-api-server-pro user add <username> <password> --admin")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
