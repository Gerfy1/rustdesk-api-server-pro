package cmd

import (
	"fmt"
	"rustdesk-api-server-pro/app/model"
	"rustdesk-api-server-pro/config"
	"rustdesk-api-server-pro/db"
	"rustdesk-api-server-pro/util"

	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User management",
}

var isAdmin bool
var userRole int

var userAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add user [ add username password][--admin or --role=N]",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Usage: user add <username> <password> [--admin] [--role=N]")
			fmt.Println("Roles: 1=User, 2=Support, 3=Support N2, 4=Super Admin")
			return
		}

		username := args[0]
		password, _ := util.Password(args[1])

		// Determine role: if --admin flag is used, set as Super Admin (4)
		// Otherwise use --role flag value (default 1)
		role := userRole
		if isAdmin {
			role = model.ROLE_SUPER_ADMIN
		}
		if role < 1 || role > 4 {
			role = model.ROLE_USER
		}

		// Auto-calculate is_admin based on role
		// Role >= 3 (Support N2, Super Admin) = admin access
		isAdminAccess := role >= model.ROLE_SUPPORT_N2

		user := &model.User{
			Username:        username,
			Password:        password,
			Name:            username,
			LicensedDevices: 0,
			LoginVerify:     model.LOGIN_ACCESS_TOKEN,
			IsAdmin:         isAdminAccess,
			Role:            role,
			Status:          1,
		}
		cfg := config.GetServerConfig()
		engine, err := db.NewEngine(cfg.Db)
		if err != nil {
			fmt.Println("Db Engine create error:", err)
			return
		}

		// Sync database tables (create if not exists)
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
		fmt.Println("Database tables synced successfully!")

		_, err = engine.Insert(user)
		if err != nil {
			fmt.Println("Add error:", err)
			return
		}
		roleNames := map[int]string{
			1: "User",
			2: "Support",
			3: "Support N2",
			4: "Super Admin",
		}
		fmt.Printf("User '%s' added successfully! Role: %s (Level %d)\n", username, roleNames[role], role)
	},
}

func init() {
	userAddCmd.Flags().BoolVarP(&isAdmin, "admin", "a", false, "Set user as Super Admin (role 4)")
	userAddCmd.Flags().IntVarP(&userRole, "role", "r", 1, "Set user role: 1=User, 2=Support, 3=Support N2, 4=Super Admin")
	userCmd.AddCommand(userAddCmd)
	RootCmd.AddCommand(userCmd)
}
