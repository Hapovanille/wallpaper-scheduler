package wpscheduler

import (
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"os"
	"path/filepath"
	"strings"
)

var (
	wallpaperPath string
	timeInterval  int
)

const configPath = "config.yaml"

func readFromConfig() {
	// add driver for support yaml content
	config.AddDriver(yamlv3.Driver)

	err := config.LoadFiles(GetConfigRelativePath())
	if err != nil {
		panic(err)
	}
}

// GetConfigRelativePath works when app is built, gets the config relative path
func GetConfigRelativePath() string {
	path, _ := filepath.Abs(os.Args[0])
	index := strings.LastIndex(path, string(os.PathSeparator))
	path = filepath.Join(path[:index], configPath)
	return path
}

// GetWallpaperPath returns the wallpaper path
func GetWallpaperPath() string {
	readFromConfig()
	wallpaperPath = config.String("wallpaper.path")

	return wallpaperPath
}

// GetTimeInterval returns the time interval
func GetTimeInterval() int {
	readFromConfig()
	timeInterval = config.Int("refresh.frequency")

	return timeInterval
}
