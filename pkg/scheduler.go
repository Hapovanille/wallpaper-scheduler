package wpscheduler

import (
	"github.com/go-vgo/robotgo"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var instance *WallpaperScheduler

const (
	logDir           = "log"
	logFileName      = "log_file_wps"
	DefaultWallpaper = "/System/Library/Desktop Pictures/Monterey Graphic.heic"
)

//Start the main app
func (wp *WallpaperScheduler) Start() {
	if wp.isRunning() {
		return
	}
	wp.quit = make(chan bool)
	wp.ticker = time.NewTicker(1 * time.Second)
	wp.updateRunningStatus(true)
	logger := getLogger(wp, false, logFileName)
	logger.Infof("scheduler: running status set")
	success, err := SetWallpaper(DefaultWallpaper)
	if err != nil {
		logger.Infof("scheduler: set failed by %s", err)
		go func() {
			robotgo.ShowAlert("scheduler Update wallpaper fails", time.Now().String())
		}()
	}

	go func() {
		for {
			select {
			case t := <-wp.ticker.C:
				success, err = SetWallpaper(DefaultWallpaper)
				if err != nil {
					logger.Infof("scheduler: set failed by %s", err)

					//show the error alert only three times
					if wp.getErrCount() <= 3 {
						go func() {
							robotgo.ShowAlert("scheduler: update wallpaper fails at %s", t.String())
						}()
						errCount := wp.getErrCount()
						wp.setErrCount(errCount + 1)
					}
				}
				if success {
					logger.Infof("scheduler: set succeeded, lastUpdateTime updated from %s to %s", wp.getlastUpdateTime().String(), t.String())
					wp.setlastUpdateTime(t)
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

//Quit the main app
func (wp *WallpaperScheduler) Quit() {
	if wp != nil && wp.isRunning() {
		wp.ticker.Stop()
		wp.quit <- true
	}
	if wp.logFile != nil {
		_ = wp.logFile.Close()

	}
}

//GetInstance gets the singleton instance for wallpaper scheduler app
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

func (wp *WallpaperScheduler) getlastUpdateTime() time.Time {
	wp.mutex.RLock()
	defer wp.mutex.RUnlock()
	return wp.lastUpdateTime
}

func (wp *WallpaperScheduler) setlastUpdateTime(time time.Time) {
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
