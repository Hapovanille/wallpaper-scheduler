# Wallpaper Scheduler


[![Go Report Card](https://goreportcard.com/badge/github.com/Hapovanille/wpscheduler)](https://goreportcard.com/report/github.com/Hapovanille/wpscheduler)

The idea of this small project is because the IT department of my company resets my laptop wallpaper randomly and frequently. Writing a go app is definitely not the most straight forward solution. Some more straightforward approaches could be cronjobs, Automator workflows, systemd, etc.

## Permission requirement

<img src="./resources/access-required-storage.png" width=200>
<img src="./resources/access-required-event.png" width=200>

## Usage

### Build

```bash
$ make build
```

### Run
Double-click the app and grant the permission of "Accessibility" when first time running the app.

Or run the app from terminal:

```bash
./bin/wpscheduler.app/Contents/MacOS/wpscheduler
```


## Todo
1. make the app configurable by introducing a yaml config file.
2. maybe show a red `!` on the system tray icon to show the app is not functioning
