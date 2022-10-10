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
		wpStart := systray.AddMenuItem("Start", "Start the app")
		wpForceUpdate := systray.AddMenuItem("Force Update", "Force update the default wallpaper")
		wpStop := systray.AddMenuItem("Stop", "Stop the app")
		wpStop.Disable()
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

		wpScheduler := wpscheduler.GetInstance()
		wpScheduler.Start()
		wpStart.Disable()
		wpStop.Enable()
		wpForceUpdate.Enable()
		for {
			select {
			case <-wpStart.ClickedCh:
				log.Infof("main: Starting the app")
				wpScheduler.Start()
				wpStart.Disable()
				wpStop.Enable()
				log.Infof("main: App started")

			case <-wpForceUpdate.ClickedCh:
				log.Infof("main: Force updating the default wallpaper")
				_, _ = wpscheduler.SetWallpaper(wpscheduler.DefaultWallpaper) //ignore return values
				log.Infof("main: Wallpaper update attempted")

			case <-wpStop.ClickedCh:
				log.Infof("main: Stopping the app")
				wpStart.Enable()
				wpStop.Disable()
				wpScheduler.Quit()
				log.Infof("main: App stopped")

			case <-mQuit.ClickedCh:
				log.Infof("main: Quiting the app")
				wpScheduler.Quit()
				systray.Quit()
				return

			case <-about.ClickedCh:
				log.Infof("main: Requesting about page")
				robotgo.ShowAlert("Wallpaper Scheduler app v1.2.0", "Developed by Hapovanille. \n\nMore info at: https://github.com/Hapovanille/wallpaper-scheduler")
			}
		}
	}()
}

func onExit() {
	// Cleaning stuff here.
	log.Infof("main: App quited")
}
