package main

import (
// "fmt"
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

type cmdType int

const (
	moveFwd cmdType = iota
	moveBwd
	moveUp
	moveDn
)

const (
    moonMoveLimit = 0.5
    sunMoveLimit  = 0.5
)

type roboClaw struct {
	cmd   cmdType
	value []byte
}

var azVelProfile = []byte{0x10, 0x20, 0x30, 0x50}
var elVelProfile = []byte{0x10, 0x20, 0x30, 0x50}



func (app *application) moveAz(az float64) (err error) {
	rc := &roboClaw{value: make([]byte, 1)}
	go func() {
		for _, v := range azVelProfile {
			switch {
			case az >= app.currAz:
				rc.cmd = moveFwd
			case az < app.currAz:
				rc.cmd = moveBwd
			}
			rc.value[0] = v
			err := app.writeCmd(rc)
			if err != nil {
				return
			}
		}
	}()
	return err
}

func (app *application) moveEl(el float64) (err error) {
	rc := &roboClaw{value: make([]byte, 1)}
	go func() {
		for _, v := range elVelProfile {
			switch {
			case el >= app.currEl:
				rc.cmd = moveUp
			case el < app.currEl:
				rc.cmd = moveDn
			}
			rc.value[0] = v
			err := app.writeCmd(rc)
			if err != nil {
				return
			}
		}
	}()
	return nil
}
