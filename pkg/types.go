package wpscheduler

import (
	"os"
	"sync"
	"time"
)

//WallpaperScheduler is the main struct for the app
type WallpaperScheduler struct {
	isAppRunning   bool
	quit           chan bool
	ticker         *time.Ticker
	lastUpdateTime time.Time
	mutex          sync.RWMutex
	logFile        *os.File
}
