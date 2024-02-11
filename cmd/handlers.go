package main

import (
	"fmt"
	"log"
	//	"math"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	//    "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

type controllerTime struct {
	year  int
	month int
	day   int
	hour  float64
	min   float64
	sec   float64
	ut    float64
}

const (
	azPulses = azMul //for the sub lunar rotator  200 //for the 3m dish
	elPulses = elMul //for the sub lunar rotator 200 //for the 3m dish
)

func (app *application) mooner() {
	go func() {
		ct := controllerTime{}
		for {
			switch app.state {
			case IDLE:
				continue
			case TRACKING_MOON:
				ct.getTime()
				_, _, _, _, _, _, az, el, _ := moon2(ct.year, ct.month, ct.day, ct.ut, app.lon, app.lat)
				if checkLimits(az, app.maxAz, app.minAz) && checkLimits(el, app.maxEl, app.minEl) {
					app.currAz = az
					app.currEl = el
					fmt.Println(az, el)
					app.reSync()
				}
				fmt.Printf("Moon Azimuth: %5.2f\tMoon Elevation: %5.2f\n", az, el)
			case TRACKING_SUN:
				ct.getTime()
				_, _, _, az, el, _, _ := sun(ct.year, ct.month, ct.day, ct.ut, app.lon, app.lat)
				if checkLimits(az, app.maxAz, app.minAz) && checkLimits(el, app.maxEl, app.minEl) {
					app.currAz = az
					app.currEl = el
					app.reSync()
				}
				fmt.Printf("Sun Azimuth: %5.2f\tSun Elevation: %5.2f\n", az, el)
			case PARKED:
				app.reSync()
			}
			time.Sleep(time.Duration(30) * time.Second)
		}
	}()
}

func (ct *controllerTime) getTime() {

	localTime := time.Now()
	t := localTime.UTC()
	ct.year = t.Year()
	ct.month = int(t.Month())
	ct.day = t.Day()
	ct.hour = float64(t.Hour())
	ct.min = float64(t.Minute())
	ct.sec = float64(t.Second())
	ct.ut = ct.hour + (ct.min / 60.0) + (ct.sec / 3600.0)
}

func (app *application) updateGrid(grid string) {
	lat, lon, err := grid2Deg(grid)
	if err != nil {
		msg := fmt.Sprintf("Grid value err: %v: You entered: \"%s\" try again", err, grid)
		app.handleError(msg)
		return
	}
	app.grid = grid
	app.lat = lat
	app.lon = lon
	app.reSync()
	app.saveMasterData()

	fmt.Println("Lattitude: ", lat)
	fmt.Println("Longitude: ", lon)
}

func (app *application) updatePark(azimuth, elevation string) {
	var msg string
	az, err := strconv.ParseFloat(azimuth, 64)
	if err != nil {
		msg = fmt.Sprintf("Park azimuth entry error: %v.  It is not a number, you entered \"%s\", try again", err, azimuth)
		app.handleError(msg)
		return
	}
	if az < 0.0 {
		msg = fmt.Sprintf("Park azimuth cannot be less than 0, you entered \"%s\" try again", azimuth)
		app.handleError(msg)
		return
	}
	if az > 360.0 {
		msg = fmt.Sprintf("Park azimuth cannot be more than 360, you entered \"%s\" try again", azimuth)
		app.handleError(msg)
		return
	}
	if az > app.maxAz {
		msg = fmt.Sprintf("Park azimuth cannot be larger than max azimuth \"%s\" try again", azimuth)
		app.handleError(msg)
		return
	}
	if az < app.minAz {
		msg = fmt.Sprintf("Park azimuth cannot be smaller than min azimuth \"%s\" try again", azimuth)
		app.handleError(msg)
		return
	}
	el, err := strconv.ParseFloat(elevation, 64)
	if err != nil {
		msg = fmt.Sprintf("Park elevation enitry error: %v  It is not a number, you entered \"%s\", try again", err, elevation)
		app.handleError(msg)
		return
	}
	if el < 0.0 {
		msg = fmt.Sprintf("Park elevation cannot be less than 0, you entered \"%s\" try again", elevation)
		app.handleError(msg)
		return
	}
	if el > 90.0 {
		msg = fmt.Sprintf("Park elevation cannot be more than 90, you entered \"%s\" try again", elevation)
		app.handleError(msg)
		return
	}
	if el > app.maxEl {
		msg = fmt.Sprintf("Park elevation cannot be larger than max elevation \"%s\" try again", elevation)
		app.handleError(msg)
		return
	}
	if el < app.minEl {
		msg = fmt.Sprintf("Park elevation cannot be smaller than min elevation \"%s\" try again", elevation)
		app.handleError(msg)
		return
	}
	app.parkAz = az
	app.parkEl = el
	app.reSync()
	app.saveMasterData()
	app.state = PARKED

	fmt.Println("Park Azimuth: ", azimuth)
	fmt.Println("Park Elevation: ", elevation)
}

func (app *application) updateMinMax(minAz, maxAz, minEl, maxEl string) {
	var msg string
	mina, err := strconv.ParseFloat(minAz, 64)
	if err != nil {
		msg = fmt.Sprintf("Min. Az. entry error: %v.  It is not a number, you entered \"%s\", try again", err, minAz)
		app.handleError(msg)
		return
	}
	if mina < 0.0 {
		msg = fmt.Sprintf("Min. Az. cannot be less than 0, you entered \"%s\" try again", minAz)
		app.handleError(msg)
		return
	}
	if mina > 360.0 {
		msg = fmt.Sprintf("Min. Az. cannot be more than 360, you entered \"%s\" try again", minAz)
		app.handleError(msg)
		return
	}
	maxa, err := strconv.ParseFloat(maxAz, 64)
	if err != nil {
		msg = fmt.Sprintf("Max. Az. entry error: %v.  It is not a number, you entered \"%s\", try again", err, maxAz)
		app.handleError(msg)
		return
	}

	if maxa < 0.0 {
		msg = fmt.Sprintf("Max. Az. cannot be less than 0, you entered \"%s\" try again", maxAz)
		app.handleError(msg)
		return
	}
	if maxa > 360.0 {
		msg = fmt.Sprintf("Max. Az. cannot be more than 360, you entered \"%s\" try again", maxAz)
		app.handleError(msg)
		return
	}

	mine, err := strconv.ParseFloat(minEl, 64)
	if err != nil {
		msg = fmt.Sprintf("Min. El. entry error: %v.  It is not a number, you entered \"%s\", try again", err, minEl)
		app.handleError(msg)
		return
	}
	if mine < 0.0 {
		msg = fmt.Sprintf("Min El. cannot be less than 0, you entered \"%s\" try again", minEl)
		app.handleError(msg)
		return
	}
	if mine > 90.0 {
		msg = fmt.Sprintf("Min. El. cannot be more than 90, you entered \"%s\" try again", minEl)
		app.handleError(msg)
		return
	}

	maxe, err := strconv.ParseFloat(maxEl, 64)
	if err != nil {
		msg = fmt.Sprintf("Max. El. entry error: %v.  It is not a number, you entered \"%s\", try again", err, maxEl)
		app.handleError(msg)
		return
	}
	if maxe < 0.0 {
		msg = fmt.Sprintf("Max. El. cannot be less than 0, you entered \"%s\" try again", maxEl)
		app.handleError(msg)
		return
	}
	if maxe > 90.0 {
		msg = fmt.Sprintf("Max. El. cannot be more than 90, you entered \"%s\" try again", maxEl)
		app.handleError(msg)
		return
	}

	if maxe < mine {
		msg = fmt.Sprintf("Max. El. \"%s\" must be larger than Min. El. \"%s\" try again", maxEl, minEl)
		app.handleError(msg)
		return
	}
	if maxa < mina {
		msg = fmt.Sprintf("Max. Az. \"%s\" must be larger than Min. Az. \"%s\" try again", maxAz, minAz)
		app.handleError(msg)
		return
	}
	app.minAz = mina
	app.maxAz = maxa
	app.minEl = mine
	app.maxEl = maxe
	app.reSync()
	app.saveMasterData()

	fmt.Println("Min Azimuth: ", minAz)
	fmt.Println("Max Azimuth: ", maxAz)
	fmt.Println("Min Elevation: ", minEl)
	fmt.Println("Max Elevation: ", maxEl)
}

func (app *application) updateTarget(azimuth, elevation string) {
	var msg string

	az, err := strconv.ParseFloat(azimuth, 64)
	if err != nil {
		msg = fmt.Sprintf("Target Az. entry error: %v.  It is not a number, you entered \"%s\", try again", err, azimuth)
		app.handleError(msg)
		return
	}
	if az < 0.0 {
		msg = fmt.Sprintf("Target Az. cannot be less than 0, you entered \"%s\" try again", azimuth)
		app.handleError(msg)
		return
	}
	if az > 360.0 {
		msg = fmt.Sprintf("Target Az. cannot be more than 360, you entered \"%s\" try again", azimuth)
		app.handleError(msg)
		return
	}
	if az > app.maxAz {
		msg = fmt.Sprintf("Target Az cannot be larger than max Az, you entered \"%s\", try again", azimuth)
		app.handleError(msg)
		return
	}
	if az < app.minAz {
		msg = fmt.Sprintf("Target Az cannot be smaller than min Az, you entered \"%s\", try again", azimuth)
		app.handleError(msg)
		return
	}

	el, err := strconv.ParseFloat(elevation, 64)
	if err != nil {
		msg = fmt.Sprintf("Target El. entry error: %v.  It is not a number, you entered \"%s\", try again", err, elevation)
		app.handleError(msg)
		return
	}
	if el < 0.0 {
		msg = fmt.Sprintf("Target El. cannot be less than 0, you entered \"%s\" try again", elevation)
		app.handleError(msg)
		return
	}
	if el > 90.0 {
		msg = fmt.Sprintf("Target El cannot be more than 90, you entered \"%s\" try again", elevation)
		app.handleError(msg)
		return
	}
	if el > app.maxEl {
		msg = fmt.Sprintf("Target El cannot be more than max El, you entered \"%s\" try again", elevation)
		app.handleError(msg)
		return
	}
	if el < app.minEl {
		msg = fmt.Sprintf("Target El cannot be less than min El, you entered \"%^s\", try again", elevation)
		app.handleError(msg)
		return
	}
	app.currAz = az
	app.currEl = el
	app.saveDishData()
	app.reSync()

	fmt.Println("Target Azimuth: ", app.currAz, app.azPosition)
	fmt.Println("Target Elevattion: ", app.currEl, app.elPosition)
}

func (app *application) adjustUp() {
	tst := app.currEl + 0.5
	if tst < app.maxEl && tst < 90.0 {
		app.currEl = tst
		app.reSync()
		return
	}
	return
}

func (app *application) adjustDn() {
	tst := app.currEl - 0.5
	if tst > app.minEl && tst > 0.0 {
		app.currEl = tst
		app.reSync()
		return
	}
	return
}

func (app *application) adjustRight() {
	tst := app.currAz + 0.5
	if tst < app.maxAz && tst < 360.0 {
		app.currAz = tst
		app.reSync()
		return
	}
	return
}

func (app *application) adjustLeft() {
	tst := app.currAz - 0.5
	if tst > app.minAz && tst > 0.0 {
		app.currAz = tst
		app.reSync()
		return
	}
	return
}

func (app *application) trackModeSelect(value string) {
	switch value {
	case SUN:
		app.selection = TRACKING_SUN
	case MOON:
		app.selection = TRACKING_MOON
	default:
		app.selection = IDLE
	}
}

func (app *application) pushedTrack() {
	var testUpdate bool
	ct := controllerTime{}
	switch app.selection {
	case TRACKING_SUN:

		app.state = TRACKING_SUN
		ct.getTime()
		_, _, _, az, el, _, _ := sun(ct.year, ct.month, ct.day, ct.ut, app.lon, app.lat)
		testUpdate = checkLimits(az, app.maxAz, app.minAz) && checkLimits(el, app.maxEl, app.minEl)
		if testUpdate {
			app.currAz = az
			app.currEl = el
			app.reSync()
		} else {
			app.state = IDLE
			msg := fmt.Sprintf("Sun Az: %5.1f\t El: %5.1f\t is outside the system limits", az, el)
			app.handleError(msg)
		}
	case TRACKING_MOON:
		app.state = TRACKING_MOON
		ct.getTime()
		_, _, _, _, _, _, az, el, _ := moon2(ct.year, ct.month, ct.day, ct.ut, app.lon, app.lat)
		testUpdate = checkLimits(az, app.maxAz, app.minAz) && checkLimits(el, app.maxEl, app.minEl)
		if testUpdate {
			app.currAz = az
			app.currEl = el
			app.reSync()
		} else {
			app.state = IDLE
			msg := fmt.Sprintf("Moon Az: %5.1f\t El: %5.1f\t is outside the system limits", az, el)
			app.handleError(msg)
		}
	default:
		app.state = IDLE
	}
}

func (app *application) pushedPark() {
	app.state = PARKED
	app.currAz = app.parkAz
	app.currEl = app.parkEl
	app.reSync()
}

func (app *application) pushedStop() {
	app.currAz = app.azPosition
	app.currEl = app.elPosition
	app.state = IDLE
	app.reSync()
}

func checkLimits(azel, highLimit, lowLimit float64) bool {
	if azel > highLimit {
		return false
	}
	if azel < lowLimit {
		return false
	}
	return true
}

func (app *application) handleError(msg string) {
	errW := app.ap.NewWindow(TITLE_ERR)
	t := &textWrap{txt: msg, txtClr: black, txtBld: false, bgClr: white}
	errText := t.makeText()
	row1 := container.New(layout.NewGridLayout(1), errText)
	row2 := seperator()

	b := &buttonWrap{
		txt:    BUTTON_OK,
		txtClr: black,
		txtBld: false,
		bgClr:  white,
		callBack: func() {
			errW.Close()
		},
	}
	okButton := b.makeButton()

	row3 := container.New(layout.NewGridLayout(1), okButton)

	sug := container.NewBorder(nil, row3, nil, nil,
		container.New(layout.NewVBoxLayout(), row1, row2))

	errW.SetContent(sug)
	errW.Resize(fyne.NewSize(300, 80))

	errW.Show()

}

func (app *application) reSync() {
	err := app.azBind.Set(fmt.Sprintf("%5.2f", app.currAz))
	if err != nil {
		log.Fatal("resync data failed in currAz: ", err)
	}
	err = app.elBind.Set(fmt.Sprintf("%5.2f", app.currEl))
	if err != nil {
		log.Fatal("resync data failed in currEl: ", err)
	}
	err = app.azPosBind.Set(fmt.Sprintf("%5.2f", app.azPosition))
	if err != nil {
		log.Fatal("resync datta failed in azPosition")
	}
	err = app.elPosBind.Set(fmt.Sprintf("%5.2f", app.elPosition))
	if err != nil {
		log.Fatal("resync datta failed in elPosition")
	}
	err = app.gridBind.Set(fmt.Sprintf("%s%s", TEXT_GRID_VALUE, app.grid))
	if err != nil {
		log.Fatal("resync data failed in grid: ", err)
	}
	err = app.parkAzBind.Set(fmt.Sprintf("%s%5.2f", TEXT_CURRENT_VALUE, app.parkAz))
	if err != nil {
		log.Fatal("resync data failed in parkAz: ", err)
	}
	err = app.parkElBind.Set(fmt.Sprintf("%s%5.2f", TEXT_CURRENT_VALUE, app.parkEl))
	if err != nil {
		log.Fatal("resync data failed in parkEl: ", err)
	}
	err = app.maxAzBind.Set(fmt.Sprintf("%s%5.2f", TEXT_CURRENT_VALUE, app.maxAz))
	if err != nil {
		log.Fatal("resync data failed in maxAz: ", err)
	}
	err = app.minAzBind.Set(fmt.Sprintf("%s%5.2f", TEXT_CURRENT_VALUE, app.minAz))
	if err != nil {
		log.Fatal("resync data failed in minAz: ", err)
	}
	err = app.maxElBind.Set(fmt.Sprintf("%s%5.2f", TEXT_CURRENT_VALUE, app.maxEl))
	if err != nil {
		log.Fatal("resync data failed maxEl: ", err)
	}
	err = app.minElBind.Set(fmt.Sprintf("%s%5.2f", TEXT_CURRENT_VALUE, app.minEl))
	if err != nil {
		log.Fatal("resync data failed minEl: ", err)
	}
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
}

func (app *application) recalibrate(azimuth, elevation string) {
	fmt.Println("recalibratte")
	var msg string

	az, err := strconv.ParseFloat(azimuth, 64)
	if err != nil {
		msg = fmt.Sprintf("Target Az. entry error: %v.  It is not a number, you entered \"%s\", try again", err, azimuth)
		app.handleError(msg)
		return
	}
	if az < 0.0 {
		msg = fmt.Sprintf("Target Az. cannot be less than 0, you entered \"%s\" try again", azimuth)
		app.handleError(msg)
		return
	}
	if az > 360.0 {
		msg = fmt.Sprintf("Target Az. cannot be more than 360, you entered \"%s\" try again", azimuth)
		app.handleError(msg)
		return
	}
	if az > app.maxAz {
		msg = fmt.Sprintf("Target Az cannot be larger than max Az, you entered \"%s\", try again", azimuth)
		app.handleError(msg)
		return
	}
	if az < app.minAz {
		msg = fmt.Sprintf("Target Az cannot be smaller than min Az, you entered \"%s\", try again", azimuth)
		app.handleError(msg)
		return
	}

	el, err := strconv.ParseFloat(elevation, 64)
	if err != nil {
		msg = fmt.Sprintf("Target El. entry error: %v.  It is not a number, you entered \"%s\", try again", err, elevation)
		app.handleError(msg)
		return
	}
	if el < 0.0 {
		msg = fmt.Sprintf("Target El. cannot be less than 0, you entered \"%s\" try again", elevation)
		app.handleError(msg)
		return
	}
	if el > 90.0 {
		msg = fmt.Sprintf("Target El cannot be more than 90, you entered \"%s\" try again", elevation)
		app.handleError(msg)
		return
	}
	if el > app.maxEl {
		msg = fmt.Sprintf("Target El cannot be more than max El, you entered \"%s\" try again", elevation)
		app.handleError(msg)
		return
	}
	if el < app.minEl {
		msg = fmt.Sprintf("Target El cannot be less than min El, you entered \"%^s\", try again", elevation)
		app.handleError(msg)
		return
	}

	azRegister := uint32(az * azPulses)
	elRegister := uint32(el * elPulses)

	err = app.writeQuadRegister(azRegister, "az")
	if err != nil {
		log.Printf("Updating Az register failed: %v", err)
	}
	err = app.writeQuadRegister(elRegister, "el")
	if err != nil {
		log.Printf("Updating El register failed: %v", err)
	}
	app.azPosition = az
	app.currAz = az
	app.elPosition = el
	app.currEl = el
	app.reSync()

}
