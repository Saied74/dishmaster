package main

import (
	"errors"
	"fmt"
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
	timeOut
	azEncMode
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
	moonMoveLimit = 0.2   // for the Sub Lunar // 0.2
	sunMoveLimit  = 0.2   // for the Sub Lunar // 0.2
	azMul         = 0.150 //1.67  // for the Sub Lunar // 200.0  //3m dish only
	elMul         = 0.150 //1.67  // for the Sub Lunar // 200.0  //3m dish only
	dT            = 500   //milliseconds
	azFast        = 120.0 //ramp velocity up and down steps (out of 127)
	elFast        = 120.0 //fir the sublunar
	azSlow        = 30.0  //65.0
	elSlow        = 30.0  //65.0
	azEndLimit    = 0.21  //degrees of travel needed to stop the motor
	elEndLimit    = 0.21
	near          = 1.5    // for the Sub Lunar //1.0 //distance to be considered on target
	revEnc        = 1 << 6 //see pages 73 and 74 of roboclaw user manual
	revMot        = 1 << 5
	faultLimit    = 6
)

type roboClaw struct {
	cmd   int
	value byte
}

func (app *application) move() {
	var limit float64 = 0.2 // for the Sub Lunar // 0.2
	var azPhase, elPhase int = phaseZero, phaseZero
	var err error
	go func() {
		for {
			//if app.faultCnt > faultLimit {
			//	continue
			//}
			switch app.state { //Select one of 3 limits, moon, sun or default.
			case TRACKING_MOON:
				limit = moonMoveLimit
			case TRACKING_SUN:
				limit = sunMoveLimit
			}
			err = app.setAzPosition()    //Read the quadrature encoder registers and set the position variables
			if errors.Is(err, noReadN) { //noReadN is caused by the remote falling behind we just need to wait
				continue
			}
			err = app.setElPosition()
			if errors.Is(err, noReadN) { //noReadN is caused by the remote falling behind we just need to wait
				continue
			}
			azPosition, _ := app.getPosition() //race condiition issue
			currAz, _ := app.getCurr()         //race condition issue
			az := math.Abs(currAz - azPosition)
			switch {
			case az > near:
				azPhase = phaseTwo
			case az > limit:
				azPhase = phaseOne
			default:
				azPhase = phaseZero
			}
			azPosition, _ = app.getPosition() //race condirtion issue
			currAz, _ = app.getCurr()         //race condition issue
			azD := currAz - azPosition        //+ means move Fwd, - means move Bwd
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
			_, elPosition := app.getPosition()  //race condition issue
			_, currEl := app.getCurr()          //race iissue
			el := math.Abs(currEl - elPosition) //race condition issue/
			switch {
			case el > near:
				elPhase = phaseTwo
			case el > limit:
				elPhase = phaseOne
			default:
				elPhase = phaseZero
			}
			_, elPosition = app.getPosition() //race condition issue
			_, currEl = app.getCurr()         //race condition issue
			elD := currEl - elPosition        // + means move up, - means move down
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
			app.saveDishData()
			app.reSync()
			time.Sleep(time.Duration(dT) * time.Millisecond)
		}
	}()
}

func (app *application) moveError(err error, reason string) {
	if app.port == nil {
		return
	}
	//app.faultCnt++
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
	//to handle negative numbers (instead of really large positive ones)
	p := float64(int32(int32(azQuad))) * azMul
	if p > app.maxAz || p < app.minAz {
		//app.faultCnt = faultLimit + 1
		//app.pushedStop()
		return fmt.Errorf("read crazy number %5.0f from azimuth quadrature register", p)
	}
	//to handle the dish position initialization issue.
	if err == nil {
		app.writeAzPosition(p) //race condition issue
	}
	//app.azPosition = p                  //race condition issue
	//	az := float64(azQuad) * azMul
	//	switch {
	//	case az >= app.maxAz:
	//		app.azPosition = app.maxAz
	//	case az <= app.minAz:
	//		app.azPosition = app.minAz
	//	default:
	//		app.azPosition = az
	//	}
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
	//to handle negative numbers (instead of really large positive ones)
	p := float64(int32(elQuad)) * elMul
	if p > app.maxEl || p < app.minEl {
		//app.faultCnt = faultLimit + 1
		//app.pushedStop()
		return fmt.Errorf("read crazy number %5.0f from elevation quadrature register", p)
	}
	//to handle the dish position initialization issue
	if err == nil {
		app.writeElPosition(p) //race condition issue
	}
	//app.elPosition = p            //race condition issue
	//	el := float64(elQuad) * elMul
	//	switch {
	//	case el >= app.maxEl:
	//		app.elPosition = app.maxEl
	//	case el <= app.minEl:
	//		app.elPosition = app.minEl
	//	default:
	//		app.elPosition = el
	//	}
	return nil
}
