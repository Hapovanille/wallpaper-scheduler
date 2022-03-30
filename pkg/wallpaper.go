package wpscheduler

// Forked from https://github.com/reujab/wallpaper/blob/master/darwin.go
import (
	"os/exec"
	"strconv"
	"strings"
)

// GetWallpaper Get returns the path to the current wallpaper.
func GetWallpaper() (string, error) {
	stdout, err := exec.Command("osascript", "-e", `tell application "Finder" to get POSIX path of (get desktop picture as alias)`).Output()
	if err != nil {
		return "", err
	}

	// is calling strings.TrimSpace() necessary?
	return strings.TrimSpace(string(stdout)), nil
}

// SetFromFile uses AppleScript to tell Finder to set the desktop wallpaper to specified file.
func SetFromFile(file string) error {
	return exec.Command("osascript", "-e", `tell application "System Events" to tell every desktop to set picture to `+strconv.Quote(file)).Run()
}

// SetWallpaper Check if the wallpaper is changed. If so, update it to the default one
func SetWallpaper(wallpaper string) (bool, error) {

	background, err := GetWallpaper()
	if err != nil {
		return false, err
	}

	// check if current wallpaper is different than the default one
	if background != wallpaper {
		err = SetFromFile(wallpaper)
		if err != nil {
			return false, err
		}
		// SetFromFile succeed
		return true, nil
	}

	// wallpaper never changed
	return false, nil
}
