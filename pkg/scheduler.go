package wpscheduler

import (
	"errors"
	"github.com/go-vgo/robotgo"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	instance     *WallpaperScheduler
	TimeInterval = 60
	Wallpaper    = "/System/Library/Desktop Pictures/Monterey Graphic.heic"
)

const (
	logDir      = "log"
	logFileName = "log_file_wps"
)

// Start the main app
func (wp *WallpaperScheduler) Start() {
	if wp.isRunning() {
		return
	}
	logger := getLogger(wp, false, logFileName)
	logger.Infof("scheduler: running status set")

	if _, err := os.Stat(GetConfigRelativePath()); errors.Is(err, os.ErrNotExist) {
		logger.Infof("scheduler: config file does not exist")
	} else {
		logger.Infof("scheduler: config file exists")
		if GetTimeInterval() > 0 {
			TimeInterval = GetTimeInterval()
			logger.Infof("config: time interval set to %d", TimeInterval)
		}
		if GetWallpaperPath() != "" {
			Wallpaper = GetWallpaperPath()
			logger.Infof("config: wallpaper set to %s", Wallpaper)
		}
	}

	wp.quit = make(chan bool)
	wp.ticker = time.NewTicker(time.Duration(TimeInterval) * time.Second)
	wp.updateRunningStatus(true)

	success, err := SetWallpaper(Wallpaper)
	if err != nil {
		logger.Infof("scheduler: set failed by %s", err)
		go func() {
			robotgo.ShowAlert("scheduler Update wallpaper fails", time.Now().String())
			errCount := wp.getErrCount()
			wp.setErrCount(errCount + 1)
		}()
	}

	go func() {
		for {
			select {
			case t := <-wp.ticker.C:
				success, err = SetWallpaper(Wallpaper)
				if err != nil {
					logger.Infof("scheduler: set failed by %s", err)

					//show the error alert at most three times
					if wp.getErrCount() < 3 {
						go func() {
							robotgo.ShowAlert("scheduler: update wallpaper fails at %s", t.String())
						}()
						errCount := wp.getErrCount()
						wp.setErrCount(errCount + 1)
					}
				}
				if success {
					logger.Infof("scheduler: set succeeded, lastUpdateTime updated from %s to %s", wp.getLastUpdateTime().String(), t.String())
					wp.setLastUpdateTime(t)
					wp.setErrCount(0)
				} else {
					logger.Infof("scheduler: unchanged %s", t.String())
				}
			case <-wp.quit:
				wp.updateRunningStatus(false)
				return
			}

		}
	}()
}

// Quit the main app
func (wp *WallpaperScheduler) Quit() {
	if wp != nil && wp.isRunning() {
		wp.ticker.Stop()
		wp.quit <- true
	}
	if wp.logFile != nil {
		_ = wp.logFile.Close()

	}
}

// GetInstance gets the singleton instance for wallpaper scheduler app
func GetInstance() *WallpaperScheduler {
	if instance == nil {
		instance = &WallpaperScheduler{}
	}
	return instance
}

// Utility functions
func (wp *WallpaperScheduler) isRunning() bool {
	wp.mutex.RLock()
	defer wp.mutex.RUnlock()
	return wp.isAppRunning
}

func (wp *WallpaperScheduler) updateRunningStatus(isRunning bool) {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()
	wp.isAppRunning = isRunning
}

func (wp *WallpaperScheduler) getLastUpdateTime() time.Time {
	wp.mutex.RLock()
	defer wp.mutex.RUnlock()
	return wp.lastUpdateTime
}

func (wp *WallpaperScheduler) setLastUpdateTime(time time.Time) {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()
	wp.lastUpdateTime = time
}

func (wp *WallpaperScheduler) getErrCount() int {
	wp.mutex.RLock()
	defer wp.mutex.RUnlock()
	return wp.errCount
}

func (wp *WallpaperScheduler) setErrCount(count int) {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()
	wp.errCount = count
}

func getLogger(wp *WallpaperScheduler, doWriteToFile bool, filename string) *log.Logger {
	logger := log.New()
	logger.Formatter = &log.TextFormatter{
		FullTimestamp: true,
	}

	if doWriteToFile {
		_, err := os.Stat(logDir)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.Mkdir(logDir, os.ModePerm)
				if err != nil {
					log.Fatalf("error creating dir: %v", err)
				}
			}
		}

		logFile, err := os.OpenFile(logDir+"/"+filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		logger.SetOutput(logFile)
		wp.logFile = logFile
	}

	return logger
}
