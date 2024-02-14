package main

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

const (
	TRACKING_SUN  = "trackingSun"
	TRACKING_MOON = "trackingMoon"
	PARKED        = "parked"
	IDLE          = "idle"
	basicMicro    = "03EB"
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
	currAz     float64 //this is really target position
	currEl     float64
	azPosition float64 //this is current position
	elPosition float64
	masterPath string
	dishPath   string
	port       serial.Port
	azBind     binding.String
	elBind     binding.String
	azPosBind  binding.String
	elPosBind  binding.String
	gridBind   binding.String
	parkAzBind binding.String
	parkElBind binding.String
	maxAzBind  binding.String
	minAzBind  binding.String
	maxElBind  binding.String
	minElBind  binding.String
	modeBind   binding.String
	sDa        *scaleData
	sDe        *scaleData
}

func main() {
	dataPath := "./"
	masterPath := filepath.Join(dataPath, "master.json")
	dishPath := filepath.Join(dataPath, "dish.json")

	sDa := &scaleData{
		centerX: 250.0,
		centerY: 180.0, //250.0,
		endX:    200.0,
		endY:    100.0,
	}

	sDe := &scaleData{
		centerX: 200.0,
		centerY: 180.0, //250.0,
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
		currAz:     100.0,
		currEl:     30.0,
		azPosition: 100.0,
		elPosition: 30.0,
		masterPath: masterPath,
		dishPath:   dishPath,
		azBind:     binding.NewString(),
		elBind:     binding.NewString(),
		azPosBind:  binding.NewString(),
		elPosBind:  binding.NewString(),
		gridBind:   binding.NewString(),
		parkAzBind: binding.NewString(),
		parkElBind: binding.NewString(),
		maxAzBind:  binding.NewString(),
		minAzBind:  binding.NewString(),
		maxElBind:  binding.NewString(),
		minElBind:  binding.NewString(),
		modeBind:   binding.NewString(),
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
		BaudRate: 460800,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	go func() {
		firstTime := true
		for {
			usbPort, err := findPort(basicMicro)
			if err != nil {
				if firstTime {
					log.Printf("port finding error %v", err)
					firstTime = false
				}
				app.port = nil
				continue
			}
			if app.port == nil {
				port, err := serial.Open(usbPort, mode) //   tty.usbmodemF412FA9C9C682", mode)
                port.SetReadTimeout(time.Duration(2) * time.Second)
				if err != nil {
					log.Printf("failed to open the usb connecttion %s: %v", usbPort, err)
					continue
				}
				app.port = port
				app.initApp()
				firstTime = true
                log.Printf("port %v reopened", usbPort) 
			}
			time.Sleep(time.Duration(1) * time.Second)
		}
	}()

	app.azBind.Set(fmt.Sprintf("%5.2f", app.currAz))
	app.elBind.Set(fmt.Sprintf("%5.2f", app.currEl))
	app.azPosBind.Set(fmt.Sprintf("%5.2f", app.azPosition))
	app.elPosBind.Set(fmt.Sprintf("%5.2f", app.elPosition))
	app.gridBind.Set(fmt.Sprintf("%s", app.grid))
	app.parkAzBind.Set(fmt.Sprintf("%5.2f", app.parkAz))
	app.parkElBind.Set(fmt.Sprintf("%5.2f", app.parkEl))
	app.maxAzBind.Set(fmt.Sprintf("%5.2f", app.maxAz))
	app.minAzBind.Set(fmt.Sprintf("%5.2f", app.minAz))
	app.maxElBind.Set(fmt.Sprintf("%5.2f", app.maxEl))
	app.minElBind.Set(fmt.Sprintf("%5.2f", app.minEl))

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
	app.move()
	app.screen()

}
func findPort(vid string) (port string, err error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return "", err
	}
	if len(ports) == 0 {
		return "", fmt.Errorf("no ports were found")
	}
	for _, port := range ports {
		if port.IsUSB {
			if port.VID == vid {
				return port.Name, nil
			}
		}
	}
	return "", fmt.Errorf("right usb port not found")
}

func (app *application) initApp() {

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

	var packetSerial uint16 = 0x0003
	err = app.setStdConfig(packetSerial)
	if err != nil {
		log.Printf("packet serial configuration failed %v", err)
	}
	mode := &roboClaw{cmd: azEncMode, value: revMot | revEnc}
	err = app.writeCmd(mode)
	if err != nil {
		log.Printf("reversing az motor failed %v", err)
	}
	azQPID := &pid{q: 2, p: 1, i: 0, d: 0} //defined in the comms.go file
	err = app.setVelocityPID(azQPID, "az")
	if err != nil {
		log.Printf("setting azimuth pid failed %v", err)
	}
	elQPID := &pid{q: 2, p: 1, i: 0, d: 0}
	err = app.setVelocityPID(elQPID, "el")
	if err != nil {
		log.Printf("setting elevation pid failed %v", err)
	}
	azRegister := uint32(app.currAz / azMul)
	elRegister := uint32(app.currEl / elMul)
	err = app.writeQuadRegister(azRegister, "az")
	if err != nil {
		log.Printf("Updating Az register failed: %v", err)
	}
	err = app.writeQuadRegister(elRegister, "el")
	if err != nil {
		log.Printf("Updating El register failed: %v", err)
	}
}
