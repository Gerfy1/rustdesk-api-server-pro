package app

import (
	"rustdesk-api-server-pro/app/model"
	"rustdesk-api-server-pro/config"
	"rustdesk-api-server-pro/db"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/golang-module/carbon/v2"
)

func StartJobs(cfg *config.ServerConfig) {

	dbEngine, err := db.NewEngine(cfg.Db)
	if err != nil {
		panic(err)
	}

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}
	
	// Job: Check device online status
	s.NewJob(gocron.DurationJob(time.Duration(cfg.JobsConfig.DeviceCheckJob.Duration)*time.Second), gocron.NewTask(func() {
		expired := carbon.Now(cfg.Db.TimeZone).SubSeconds(30).ToDateTimeString()
		dbEngine.Where("is_online = 1 and updated_at <= ?", expired).Cols("is_online").Update(&model.Device{
			IsOnline: false,
		})
	}))

	// Job: Close orphaned audit sessions (sessions without closed_at that are older than 2 hours)
	// This handles disconnections due to network issues, Alt+F4, crashes, etc.
	s.NewJob(gocron.DurationJob(5*time.Minute), gocron.NewTask(func() {
		twoHoursAgo := carbon.Now(cfg.Db.TimeZone).SubHours(2).ToDateTimeString()
		now := time.Now()
		
		// Find sessions that are open (closed_at is zero) and created more than 2 hours ago
		dbEngine.Where("(closed_at IS NULL OR closed_at = '0001-01-01 00:00:00') AND created_at <= ?", twoHoursAgo).
			Cols("closed_at", "note").
			Update(&model.Audit{
				ClosedAt: now,
				Note:     "Auto-closed: session timeout",
			})
	}))

	s.Start()
}
