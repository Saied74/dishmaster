package main

import (
	"errors"
	"log"
	"math"
	"time"
)

/*
move file contains the information and algorithms for moving the antenna.
The 5 meter dish rotator produces 200 pulses per degree or 72,000 pulses
in total

The 3 meter dish encoder produces 200 pules per degree on the azimuth,
or 72,000 pulses in total (actually 71920 which maps into 199.778 -
we will have to see about this)
On the 3 meter dish encoder produces 200 pulses per degree on the elevation
or 36,000 pules in total (not that we are going to go 360 degrees on elevation).

When calibrating the dish, the current position register is set to the
azimuth or elevation degrees times 200.  That is the same as the distance
(in counter counts) from the origin (zero elevation or azimuth).

Zero elevation is the local horizon and zero azymuth is true north.

I will assume channel 1 is Azimuth and channel 2 is the elevation control
*/

const (
	moveFwd = iota
	moveBwd
	moveUp
	moveDn
    phaseZero
    phaseOne
    phaseTwo
	idle
	rampUpForward
	rampUpBackward
	rampUpUp
	rampUpDwn
	fastFwd
	fastBwd
	slowFwd
	slowBwd
	fastUp
	fastDwn
	slowUp
	slowDwn
)

const (
	moonMoveLimit = 0.2
	sunMoveLimit  = 0.2
	azMul         = 200.0 //3m dish only
	elMul         = 200.0 //3m dish only
	dT            = 500   //milliseconds
	azFast        = 127.0  //ramp velocity up and down steps (out of 127)
	elFast        = 127.0
    azSlow        = 65.0
    elSlow        = 65.0
	azEndLimit    = 0.21 //degrees of travel needed to stop the motor
	elEndLimit    = 0.21
	near          = 1.0 //distance to be considered on target
)

type roboClaw struct {
	cmd   int
	value byte
}

func (app *application) move() {
	var limit float64 = 0.2
	var azPhase, elPhase int = idle, idle
	var err error
	go func() {
		for {
			switch app.state { //Select one of 3 limits, moon, sun or default.
			case TRACKING_MOON:
				limit = moonMoveLimit
			case TRACKING_SUN:
				limit = sunMoveLimit
			}
			err = app.setAzPosition() //Read the quadrature encoder registers and set the position variables
			err = app.setElPosition()
			if errors.Is(err, noReadN) { //noReadN is caused by the remote falling behind we just need to wait
				continue
			}
			//move when there is a gap between the current and desired position
            az := math.Abs(app.currAz - app.azPosition)
            switch {
            case az > near:
                azPhase = phaseTwo
            case az > limit:
                azPhase = phaseOne
            default:
                azPhase = phaseZero
            }
            azD := app.currAz - app.azPosition //+ means move Fwd, - means move Bwd
            if azPhase == phaseTwo {
                if azD >= 0 {
                    err = app.writeCmd(&roboClaw{moveFwd, azFast})
				    if err != nil {
					    app.moveError(err, "forward cmd fast")
				    }
                } else {
                    err = app.writeCmd(&roboClaw{moveBwd, azFast})
				    if err != nil {
					    app.moveError(err, "backwards cmd fast")
				    }
                }
            }
            if azPhase == phaseOne {
                if azD >= 0 {
                    err = app.writeCmd(&roboClaw{moveFwd, azSlow})
				    if err != nil {
					    app.moveError(err, "forward cmd slow")
				    }
                } else {
                    err = app.writeCmd(&roboClaw{moveBwd, azSlow})
				    if err != nil {
					    app.moveError(err, "backwards cmd slow")
				    }
                }
            }
            if azPhase == phaseZero {
                if azD >= 0 {
                    err = app.writeCmd(&roboClaw{moveFwd, 0})
				    if err != nil {
					    app.moveError(err, "forward cmd stop")
				    }
                } else {
                    err = app.writeCmd(&roboClaw{moveBwd, 0})
				    if err != nil {
					    app.moveError(err, "backwards cmd stop")
				    }
                }
            }

            el := math.Abs(app.currEl - app.elPosition)
            switch {
            case el > near:
                elPhase = phaseTwo
            case el > limit:
                elPhase = phaseOne
            default:
                elPhase = phaseZero
            }
            elD := app.currEl - app.elPosition // + means move up, - means move down
            if elPhase == phaseTwo {
                if elD >= 0 {
                    err = app.writeCmd(&roboClaw{moveUp, elFast})
				    if err != nil {
					    app.moveError(err, "up cmd fast")
				    }
                } else {
                    err = app.writeCmd(&roboClaw{moveDn, elFast})
				    if err != nil {
					    app.moveError(err, "down cmd fast")
				    }
                }
            }
            if elPhase == phaseOne {
                if elD >= 0 {
                    err = app.writeCmd(&roboClaw{moveUp, elSlow})
				    if err != nil {
					    app.moveError(err, "up cmd slow")
				    }
                } else {
                    err = app.writeCmd(&roboClaw{moveDn, elSlow})
				    if err != nil {
					    app.moveError(err, "down cmd slow")
				    }
                }
            }
            if elPhase == phaseZero {
                if elD >= 0 {
                    err = app.writeCmd(&roboClaw{moveUp, 0})
				    if err != nil {
					    app.moveError(err, "up cmd stop")
				    }
                } else {
                    err = app.writeCmd(&roboClaw{moveDn, 0})
				    if err != nil {
					    app.moveError(err, "Down cmd stop")
				    }
                }
            }
            app.reSync()
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
    az := float64(azQuad) / azMul
    switch {
    case az >= app.maxAz:
        app.azPosition = app.maxAz
    case az <= app.minAz:
        app.azPosition = app.minAz
    default: 
	    app.azPosition = az
    }	
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
    el :=  float64(elQuad) / elMul
    switch {
    case el >= app.maxEl:
        app.elPosition = app.maxEl
    case el <= app.minEl:
        app.elPosition = app.minEl
    default:
	    app.elPosition = el
    }	
    return nil
}