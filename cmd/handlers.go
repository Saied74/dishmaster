package main

import (
	"fmt"
	//	"image/color"
	"log"
	//	"runtime/debug"
	//	"strconv"
	"fyne.io/fyne/v2"
	"strconv"
	"time"
	//	ap "fyne.io/fyne/v2/app"
	//  "fyne.io/fyne/v2/data/binding"
	//	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	//	"fyne.io/fyne/v2/theme"
	//	"fyne.io/fyne/v2/widget"
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
				app.currAz = az
				app.currEl = el
				fmt.Printf("Moon Azimuth: %5.2f\tMoon Elevation: %5.2f\n", app.currAz, app.currEl)
                app.reSync()				
			case TRACKING_SUN:
                ct.getTime()
				_, _, _, az, el, _, _ := sun(ct.year, ct.month, ct.day, ct.ut, app.lon, app.lat)
				app.currAz = az
				app.currEl = el
				fmt.Printf("Sun Azimuth: %5.2f\tSun Elevation: %5.2f\n", app.currAz, app.currEl)
                app.reSync()	
			case PARKED:
                app.reSync()
				//continue
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
    fmt.Println("Month: ", ct.month)
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
	app.saveMasterData()
	fmt.Println("Grid: ", grid)
	fmt.Println("Lattitude: ", lat)
	fmt.Println("Longitude: ", lon)
}

func (app *application) updatePark(azimuth, elevation string) {
	var msg string
	az, err := strconv.ParseFloat(azimuth, 64)
	if err != nil {
		msg = fmt.Sprintf("Azimuth entry error: %v.  It is not a number, you entered \"%s\", try again", err, azimuth)
		app.handleError(msg)
		return
	}
	if az < 0.0 {
		msg = fmt.Sprintf("Azimuth cannot be less than 0, you entered \"%s\" try again", azimuth)
		app.handleError(msg)
		return
	}
	if az > 90.0 {
		msg = fmt.Sprintf("Azimuth cannot be more than 90, you entered \"%s\" try again", azimuth)
		app.handleError(msg)
		return
	}
	el, err := strconv.ParseFloat(elevation, 64)
	if err != nil {
		msg = fmt.Sprintf("Elevation enitry error: %v.  It is not a number, you entered \"%s\", try again", err, elevation)
		app.handleError(msg)
		return
	}
	if el < 0.0 {
		msg = fmt.Sprintf("Elevation cannot be less than 0, you entered \"%s\" try again", elevation)
		app.handleError(msg)
		return
	}
	if el > 360.0 {
		msg = fmt.Sprintf("Elevation cannot be more than 360, you entered \"%s\" try again", elevation)
		app.handleError(msg)
		return
	}
	app.parkAz = az
	app.parkEl = el
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
		msg = fmt.Sprintf("Target El. cannot be more than 90, you entered \"%s\" try again", elevation)
		app.handleError(msg)
		return
	}
	app.currAz = az
	app.currEl = el
	app.saveDishData()
    app.reSync()

	fmt.Println("Target Azimuth: ", azimuth)
	fmt.Println("Target Elevattion: ", elevation)
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
    if tst < app.maxAz && tst < 360.0  {
        app.currAz = tst
        app.reSync()
        return
    }
    return
}

func (app *application) adjustLeft() {
    tst := app.currAz - 0.5
    if tst > app.minAz && tst > 0.0  {
        app.currAz = tst
        app.reSync()
        return
    }
    return
}





func (app *application) trackModeSelect(value string) {
	fmt.Printf("Selected Value: %s\n", value)
	switch value {
	case SUN:
		app.selection = TRACKING_SUN
	case MOON:
		app.selection = TRACKING_MOON
	default:
		app.selection = IDLE
	}
}

func(app *application) pushedTrack() {
	fmt.Println("pushed track")
    ct := controllerTime{}
	switch app.selection {
	case TRACKING_SUN:

    	app.state = TRACKING_SUN
        ct.getTime()
        _, _, _, az, el, _, _ := sun(ct.year, ct.month, ct.day, ct.ut, app.lon, app.lat)
		app.currAz = az
		app.currEl = el
        app.reSync()	
	case TRACKING_MOON:
		app.state = TRACKING_MOON
        ct.getTime()
		_, _, _, _, _, _, az, el, _ := moon2(ct.year, ct.month, ct.day, ct.ut, app.lon, app.lat)
		app.currAz = az
		app.currEl = el
        app.reSync()	
	default:
		app.state = IDLE
	}
}

func (app *application) parkDish() {
	app.state = PARKED
	app.currAz = app.parkAz
	app.currEl = app.parkEl
    app.reSync()
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
	errW.Resize(fyne.NewSize(100, 50))

	errW.Show()

}

func (app *application) reSync() {

     err := app.azBind.Set(fmt.Sprintf("%5.2f", app.currAz))
	if err != nil {
	    log.Fatal("resync data failed: ", err)
	}
	err = app.elBind.Set(fmt.Sprintf("%5.2f", app.currEl))
	if err != nil {
		log.Fatal("resync data failed: ", err)
	}

}

