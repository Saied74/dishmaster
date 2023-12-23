package main

import (
	"fmt"
	"log"
	"path/filepath"
//	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
)

const (
	TRACKING_SUN  = "trackingSun"
	TRACKING_MOON = "trackingMoon"
	PARKED        = "parked"
	IDLE          = "idle"
)

// a button push is needed to cause state to chage based on the selection.
type application struct {
	state      string //trackingSun, trackingMoon, parked, idle
	selection  string //sun, moon, neither (idle)
	ap         fyne.App
	grid       string
	lat        float64
	lon        float64
	parkAz     float64
	parkEl     float64
	maxAz      float64
	minAz      float64
	maxEl      float64
	minEl      float64
	currAz     float64
	currEl     float64
	masterPath string
	dishPath   string
	azBind     binding.String
	elBind     binding.String
}

func main() {

	dataPath := "./"
	masterPath := filepath.Join(dataPath, "master.json")
	dishPath := filepath.Join(dataPath, "dish.json")

	app := &application{
		state:      IDLE,
		selection:  IDLE,
		grid:       "FN20RH",
		lat:        40.321490,
		lon:        -74.510240,
		parkAz:     90.0,
		parkEl:     20.0,
		maxAz:      315.0,
		minAz:      45.0,
		maxEl:      90.0,
		minEl:      28.0,
		currAz:     125.0,
		currEl:     30.0,
		masterPath: masterPath,
		dishPath:   dishPath,
		azBind:     binding.NewString(),
		elBind:     binding.NewString(),
	}
    app.saveMasterData()
    app.saveDishData()
    err := app.getMasterData()
    if err != nil {
        log.Fatalf("system failed to initialize master data because: %v", err)
    }
    err = app.getDishData()
    if err != nil {
        log.Fatalf("system failed to initiatlize dish data because: %v", err)
    }
	app.azBind.Set(fmt.Sprintf("%5.2f", app.currAz))
	app.elBind.Set(fmt.Sprintf("%5.2f", app.currEl))

	//	app.incAz()
	//	app.incEl()
	app.mooner()
	app.screen()
}


// this is for test purposes only
//func (app *application) incAz() {
//	go func() {
//		var err error
//		for {
//			time.Sleep(time.Duration(3) * time.Second)
//
//			app.currAz++
//			if app.currAz > 360.0 {
//				app.currAz = 0.0
//			}
//			err = app.azBind.Set(fmt.Sprintf("%5.2f", app.currAz))
//			if err != nil {
//				log.Fatal("from the inside the go routine: ", err)
//			}
//		}
//	}()
//}
//
//func (app *application) incEl() {
//	var err error
//	go func() {
//		for {
//			time.Sleep(time.Duration(3) * time.Second)
//
//			app.currEl++
//			if app.currEl > 90.0 {
//				app.currEl = 0
//			}
//			err = app.elBind.Set(fmt.Sprintf("%5.2f", app.currEl))
//			if err != nil {
//				log.Fatal("from the inside the go routine: ", err)
//			}
//		}
//	}()
//}
