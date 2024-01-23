package main

import (
    "errors"
//    "fmt"
	"log"
	"math"
	//    "os"
	"time"
)

/*
move file contains the information and algorithms for moving the antenna.
The 5 meter dish rotator produces 200 pulses per degree or 72,000 pulses
in total

The 3 meter dish encoder produces 200 pules per degree on the azimuth,
or 72,000 pulses in total (actually 71920 which maps into 199.778 -
we will have to see about this)
On the 3 meter dish encoder produces 100 pulses per degree on the elevation
or 36,000 pules in total (not that we are going to go 360 degrees on elevation).

When calibrating the dish, the current position register is set to the
azimuth or elevation degrees times 200.  That is the same as the distance
(in counter counts) from the origin (zero elevation or azimuth).

I will assume channel 1 is Azimuth and channel 2 is the elevation control
*/

//type cmdType int

const (
	moveFwd = iota
	moveBwd
	moveUp
	moveDn
	idle
	rampUpForward
	rampUpBackward
	rampUpUp
	rampUpDwn
	constSpeedFwd
	constSpeedBwd
	rampDwnFwd
	rampDwnBwd
	constSpeedUp
	constSpeedDwn
	rampDwnUp
	rampDwnDwn
)

const (
	moonMoveLimit = 0.5
	sunMoveLimit  = 0.5
	azMul         = 200.0 //3m dish only
	elMul         = 200.0 //3m dish only
	dT            = 500   //milliseconds
	azInc         = 21.0  //ramp velocity up and down steps (out of 127)
	elInc         = 21.0
	azEndLimit    = 0.21 //degrees of travel needed to stop the motor
	elEndLimit    = 0.21
	near          = 3.5 //distance to be considered on target
)

type roboClaw struct {
	cmd   int
	value byte
}

var azVelProfile = []byte{0x10, 0x20, 0x30, 0x50}
var elVelProfile = []byte{0x10, 0x20, 0x30, 0x50}

func (app *application) move() {
	var limit float64 = 0.5
	var azPhase, elPhase int = idle, idle
	var azVelocity, elVelocity byte
    var err error
	go func() {
		for {
			//t1 := time.Now()
			switch app.state {
			case TRACKING_MOON:
				limit = moonMoveLimit
			case TRACKING_SUN:
				limit = sunMoveLimit
			}
            err = app.setAzPosition()
			err = app.setElPosition()
            if errors.Is(err, noReadN) {
                continue
            }
			if azPhase == idle && math.Abs(app.currAz-app.azPosition) > limit {
				if app.currAz > app.azPosition {
					azPhase = rampUpForward
				} else {
					azPhase = rampUpBackward
				}
			}
			if elPhase == idle && math.Abs(app.currEl-app.elPosition) > limit {
				if app.currEl > app.elPosition {
					elPhase = rampUpUp
				} else {
					elPhase = rampUpDwn
				}
			}
			switch {
			case azPhase == rampUpForward && azVelocity < 128:
				azVelocity = azInc
				err = app.writeCmd(&roboClaw{moveFwd, azVelocity})
				if err != nil {
					app.moveError(err, "forward cmd speed up")
				}
                azPhase = constSpeedFwd
			case azPhase == rampUpBackward && azVelocity < 128:
				azVelocity = azInc
				err = app.writeCmd(&roboClaw{moveBwd, azVelocity})
				if err != nil {
					app.moveError(err, "backwards cmd speed up")
				}
                azPhase = constSpeedBwd
			case azPhase == constSpeedFwd && math.Abs(app.currAz-app.azPosition) <= near:
				azVelocity = 0
				err = app.writeCmd(&roboClaw{moveFwd, azVelocity})
				if err != nil {
					app.moveError(err, "backward cmd speed down")
				}
				azPhase = idle
			case azPhase == constSpeedBwd && math.Abs(app.currAz-app.azPosition) <= near:
				azVelocity = 0
				err = app.writeCmd(&roboClaw{moveBwd, azVelocity})
				if err != nil {
					app.moveError(err, "backward cmd speed down")
				}
				azPhase = idle
			default:
				azPhase = idle
				azVelocity = 0
				//app.moveError(nil, "default on Az move, should not be here")
			}

			switch {
			case elPhase == rampUpUp && elVelocity < 128:
				elVelocity = elInc
				err = app.writeCmd(&roboClaw{moveUp, elVelocity})
				if err != nil {
					app.moveError(err, "up cmd speed up")
				}
                elPhase = constSpeedUp
			case elPhase == rampUpDwn && elVelocity < 128:
				elVelocity = elInc
				err = app.writeCmd(&roboClaw{moveDn, elVelocity})
				if err != nil {
					app.moveError(err, "Down cmd speed up")
				}
                elPhase = constSpeedDwn
			case elPhase == constSpeedUp && math.Abs(app.currEl-app.elPosition) <= near:
				elVelocity = 0
				err = app.writeCmd(&roboClaw{moveUp, elVelocity})
				if err != nil {
					app.moveError(err, "backward cmd speed down")
				}
				elPhase = idle
			case elPhase == constSpeedDwn && math.Abs(app.currEl-app.elPosition) <= near:
				elVelocity = 0
				err = app.writeCmd(&roboClaw{moveDn, elVelocity})
				if err != nil {
					app.moveError(err, "backward cmd speed down")
				}
				elPhase = idle
			default:
				elPhase = idle
				elVelocity = 0
				//app.moveError(nil, "default on El move, should not be here")
			}
			//t2 := time.Now()
			//t3 := t2.Sub(t1)
			//time.Sleep(time.Duration(dT*time.Millisecond-t3) * time.Millisecond)
			time.Sleep(time.Duration(dT) * time.Millisecond)

		}
	}()
}

func (app *application) moveError(err error, reason string) {
	if app.port == nil {
		return
	}
	log.Printf("comm error %s: %v \n", reason, err)
}

func (app *application) setAzPosition() error {
    azQuad, err := app.readQuadRegister("az")
	if err != nil {
        if errors.Is(err, noReadN) {
            return err
        }
		app.moveError(err, "az quad read")
	}
	app.azPosition = float64(azQuad) / azMul
    return nil
}

func (app *application) setElPosition() error {
    elQuad, err := app.readQuadRegister("el")
	if err != nil {
        if errors.Is(err, noReadN) {
            return err
        }
		app.moveError(err, "el quad read")
	}
	app.elPosition = float64(elQuad) / elMul
    return nil
}
