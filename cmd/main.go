package main

import (
    "flag"
	"fmt"
	"log"
	"path/filepath"
	"sync"
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
	fdi           = "0403"
	timeout       = 20
    upLimit       = 3 //how many times to test remote being up
)

// a button push is needed to cause state to chage based on the selection.
type application struct {
	mu         sync.Mutex
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
	//counts the number of serial communicatio faults - see issue 9
	//it is reset by the monitoring thread in the main function
	faultCnt int
	//it indicates that the remote is not functioning.  See issue 9
	remFault   bool
	port       serial.Port
	azBind     binding.String
	elBind     binding.String
	azPosBind  binding.String
	elPosBind  binding.String
	azDiffBind binding.String
	elDiffBind binding.String
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

    vid := flag.String("vid",fdi , "USB Vendor ID")
    flag.Parse()

	sDa := &scaleData{
		centerX: 125.0, //250.0,
		centerY: 105.0, //180.0, //250.0,
		endX:    150.0, //200.0, //I don't think endX and endY are used anymore
		endY:    50.0,  //100.0,
	}

	sDe := &scaleData{
		centerX: 85.0,  //200.0,
		centerY: 125.0, //250.0,
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
		faultCnt:   0,
        port:       nil,
		remFault:   true, //start with the remote fault being true
		masterPath: masterPath,
		dishPath:   dishPath,
		azBind:     binding.NewString(),
		elBind:     binding.NewString(),
		azPosBind:  binding.NewString(),
		elPosBind:  binding.NewString(),
		azDiffBind: binding.NewString(),
		elDiffBind: binding.NewString(),
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
		BaudRate: 38400, //460800,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	go func() {
		firstTime := true //so we don't have repeat reporting of the same event
		for {
			//constantly monitor the port and see if a port is available
            usbPort, err := findPort(*vid) //basicMicro)
			if err != nil {
				if firstTime {
					//app.handleError(fmt.Sprintf("port finding error %v", err))
					log.Printf("port finding error %v", err)
					firstTime = false
				}
				app.setPort(nil) //no portis available
				//app.port = nil
				continue //keep looking
			}
			//if there is a port, check to see if a port is needed (app.port not nil)
			if app.getPort() == nil {
				//berore trying to open it, wait for slow computers
				time.Sleep(time.Duration(500) * time.Millisecond)
				port, err := serial.Open(usbPort, mode) //   tty.usbmodemF412FA9C9C682", mode)
				port.SetReadTimeout(time.Duration(2) * time.Second)
				if err != nil {
					log.Printf("failed to open the usb connecttion %s: %v", usbPort, err)
					continue
				}
				app.setPort(port)
                log.Println("found and set port", usbPort)
				//app.initApp()  //based on issue 9 discussin, this is too soon
				firstTime = true //since we have the port and opened it, if it fails, it will be first time again
				log.Printf("port %s reopened", usbPort)
			}
			//once a port is identified and opened, check to see if there is also a remote fault
            if app.getRemoteFault() || app.getFaultCnt() >= faultLimit {
                log.Printf("Remote Fault: %v\tFault Count: %d", app.remFault, app.faultCnt)
				if app.remoteUp() {
                    log.Printf("Remote Up")
                    app.initApp()
					app.setFaultCnt(0)
					app.setRemoteFault(false)
				}
				if app.getFaultCnt() >= faultLimit {
					if !app.remoteUp() {
						app.setRemoteFault(true)
					} else {
                        app.initApp()
						app.setFaultCnt(0)
                        app.setRemoteFault(false)
					}
				}
				time.Sleep(time.Duration(1) * time.Second)
			}
		}
	}()

	app.azBind.Set(fmt.Sprintf("%5.2f", app.currAz))
	app.elBind.Set(fmt.Sprintf("%5.2f", app.currEl))
	app.azPosBind.Set(fmt.Sprintf("%5.2f", app.azPosition))
	app.elPosBind.Set(fmt.Sprintf("%5.2f", app.elPosition))
	app.azDiffBind.Set(fmt.Sprintf("%5.2f", app.currAz-app.azPosition))
	app.elDiffBind.Set(fmt.Sprintf("%5.2f", app.currEl-app.elPosition))
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

func (app *application) remoteUp() bool {
	if app.port == nil {
		return false
	}
	for i := 0; i < upLimit; i++ {
		_, err := app.readQuadRegister("az")
		if err == nil {
			return true
		}
        time.Sleep(time.Duration(500) * time.Millisecond)
	}
	return false
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

	//	var packetSerial uint16 = 0x0003
	//	err = app.setStdConfig(packetSerial)
	//	if err != nil {
	//		log.Printf("packet serial configuration failed %v", err)
	//	}
	mode := &roboClaw{cmd: azEncMode, value: revMot | revEnc}
	err = app.writeCmd(mode)
	if err != nil {
		log.Printf("reversing az motor failed %v", err)
	}

	tOut := &roboClaw{cmd: timeOut, value: timeout}
	err = app.writeCmd(tOut)
	if err != nil {
		log.Printf("setting timeout to %d did not work because: %v", timeout, err)
	}
	azQPID := &pid{q: 8, p: 1, i: 0, d: 0} //defined in the comms.go file
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

// The folliwng 8 functions were created to address the race condition issue
func (app *application) getPosition() (az, el float64) {
	app.mu.Lock()
	az, el = app.azPosition, app.elPosition
	app.mu.Unlock()
	return az, el
}

func (app *application) getCurr() (az, el float64) {
	app.mu.Lock()
	az, el = app.currAz, app.currEl
	app.mu.Unlock()
	return az, el
}

func (app *application) writeAzPosition(az float64) {
	app.mu.Lock()
	app.azPosition = az
	app.mu.Unlock()
}

func (app *application) writeElPosition(el float64) {
	app.mu.Lock()
	app.elPosition = el
	app.mu.Unlock()
}

func (app *application) writeCurrAz(az float64) {
	app.mu.Lock()
	app.currAz = az
	app.mu.Unlock()
}

func (app *application) writeCurrEl(el float64) {
	app.mu.Lock()
	app.currEl = el
	app.mu.Unlock()
}

func (app *application) writeAzElPosition(az, el float64) {
	app.mu.Lock()
	app.azPosition = az
	app.elPosition = el
	app.mu.Unlock()
}

func (app *application) writeCurrAzEl(az, el float64) {
	app.mu.Lock()
	app.currAz = az
	app.currEl = el
	app.mu.Unlock()
}

func (app *application) setPort(port serial.Port) {
	app.mu.Lock()
	app.port = port
	app.mu.Unlock()
}

func (app *application) getPort() serial.Port {
	app.mu.Lock()
	port := app.port
	app.mu.Unlock()
	return port
}

func (app *application) getFaultCnt() int {
	app.mu.Lock()
	faultCnt := app.faultCnt
	app.mu.Unlock()
	return faultCnt
}

func (app *application) setFaultCnt(cnt int) {
	app.mu.Lock()
	app.faultCnt = cnt
	app.mu.Unlock()
}

func (app *application) getRemoteFault() bool {
	app.mu.Lock()
	flt := app.remFault
	app.mu.Unlock()
	return flt
}

func (app *application) setRemoteFault(flt bool) {
	app.mu.Lock()
	app.remFault = flt
	app.mu.Unlock()
}
