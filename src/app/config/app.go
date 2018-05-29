package config

import (
	"time"
)

type App struct {
	LocalTimeZone *time.Location
}

func Get() App {
	loc, _ := time.LoadLocation("Asia/Bangkok")
	appData := App{
		LocalTimeZone: loc,
	}
	return appData
}
