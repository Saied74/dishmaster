package main

import (
	"fmt"
	"log"
	"path/filepath"
	//	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	//    "fyne.io/fyne/v2/canvas"

	"go.bug.st/serial"
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
	port       serial.Port
//	moveLimit  float64
	azBind     binding.String
	elBind     binding.String
	gridBind   binding.String
	parkAzBind binding.String
	parkElBind binding.String
	maxAzBind  binding.String
	minAzBind  binding.String
	maxElBind  binding.String
	minElBind  binding.String
	modeBind   binding.String
//	azDialBind fyne.CanvasObject
//	elDialBind binding.FloatList
	sDa        *scaleData
	sDe        *scaleData
}

func main() {

	dataPath := "./"
	masterPath := filepath.Join(dataPath, "master.json")
	dishPath := filepath.Join(dataPath, "dish.json")

	sDa := &scaleData{
		centerX: 250.0,
		centerY: 250.0,
		endX:    250.0,
		endY:    100.0,
	}

	sDe := &scaleData{
		centerX: 200.0,
		centerY: 250.0,
		endX:    200.0,
		endY:    100.0,
	}

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
		minEl:      0.0,
		currAz:     125.0,
		currEl:     30.0,
//		moveLimit:  0.5,
		masterPath: masterPath,
		dishPath:   dishPath,
		azBind:     binding.NewString(),
		elBind:     binding.NewString(),
		gridBind:   binding.NewString(),
		parkAzBind: binding.NewString(),
		parkElBind: binding.NewString(),
		maxAzBind:  binding.NewString(),
		minAzBind:  binding.NewString(),
		maxElBind:  binding.NewString(),
		minElBind:  binding.NewString(),
		modeBind:   binding.NewString(),

//		elDialBind: binding.NewFloatList(),
		sDa:        sDa,
		sDe:        sDe,
	}
	err := app.getMasterData()
	if err != nil {
		log.Printf("System failed to initialize master data because: %v\n", err)
		log.Printf("Initializing the file in the current directory with default data")
		log.Printf("File name is master.json")
		app.saveMasterData()
	}
	err = app.getDishData()
	if err != nil {
		log.Printf("system failed to initiatlize dish data because: %v", err)
		log.Printf("Initializing the file in the current directory with default data")
		log.Printf("File name is dish.json")
		app.saveDishData()
	}

	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	usbPort := "/dev/tty.usbserial-A10L2L39"
	port, err := serial.Open(usbPort, mode) //   tty.usbmodemF412FA9C9C682", mode)
	if err != nil {
		log.Printf("failed to open the usb connecttion %s: %v", usbPort, err)
	}
	app.port = port
	app.azBind.Set(fmt.Sprintf("%5.2f", app.currAz))
	app.elBind.Set(fmt.Sprintf("%5.2f", app.currEl))
	app.gridBind.Set(fmt.Sprintf("%s", app.grid))
	app.parkAzBind.Set(fmt.Sprintf("%5.2f", app.parkAz))
	app.parkElBind.Set(fmt.Sprintf("%5.2f", app.parkEl))
	app.maxAzBind.Set(fmt.Sprintf("%5.2f", app.maxAz))
	app.minAzBind.Set(fmt.Sprintf("%5.2f", app.minAz))
	app.maxElBind.Set(fmt.Sprintf("%5.2f", app.maxEl))
	app.minElBind.Set(fmt.Sprintf("%5.2f", app.minEl))

	//    app.azDial = app.makeAzDial()

	switch app.state {
	case TRACKING_SUN:
		app.modeBind.Set("Tracking the Sun")
	case TRACKING_MOON:
		app.modeBind.Set("Tracking the Moon")
	case PARKED:
		app.modeBind.Set("Parked")
	case IDLE:
		app.modeBind.Set("Idle")
	}
	app.mooner()
	app.screen()
}
