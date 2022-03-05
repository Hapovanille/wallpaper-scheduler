package main

import (
	"github.com/Hapovanille/wpscheduler/assets/icon"
	wpscheduler "github.com/Hapovanille/wpscheduler/pkg"
	"github.com/getlantern/systray"
	"github.com/go-vgo/robotgo"
	log "github.com/sirupsen/logrus"
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	go func() {
		systray.SetIcon(icon.Data)
		about := systray.AddMenuItem("About WP", "Information about the app")
		systray.AddSeparator()
		wpStart := systray.AddMenuItem("Start", "start the app")
		wpStop := systray.AddMenuItem("Stop", "stop the app")
		wpStop.Disable()
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

		wpScheduler := wpscheduler.GetInstance()
		wpScheduler.Start()
		wpStart.Disable()
		wpStop.Enable()
		for {
			select {
			case <-wpStart.ClickedCh:
				log.Infof("main: Starting the app")
				wpScheduler.Start()
				wpStart.Disable()
				wpStop.Enable()

			case <-wpStop.ClickedCh:
				log.Infof("main: Stopping the app")
				wpStart.Enable()
				wpStop.Disable()
				wpScheduler.Quit()

			case <-mQuit.ClickedCh:
				log.Infof("main: Requesting quit")
				wpScheduler.Quit()
				systray.Quit()
				return

			case <-about.ClickedCh:
				log.Infof("main: requesting about")
				robotgo.ShowAlert("Wallpaper Scheduler app v1.0.0", "Developed by Hapovanille. \n\nMore info at: https://github.com/Hapovanille/wp-scheduler")
			}
		}
	}()
}

func onExit() {
	// Cleaning stuff here.
	log.Infof("main: finished quitting")
}
