package main

import (
	"fmt"
	"image/color"
	//	"log"
	//	"runtime/debug"
	//	"strconv"
	// 	"time"
	"fyne.io/fyne/v2"
	ap "fyne.io/fyne/v2/app"
	//	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type buttonWrap struct {
	txt      string
	txtClr   color.Color
	txtBld   bool
	bgClr    color.Color
	callBack func()
}

type textWrap struct {
	txt    string
	txtClr color.Color
	txtBld bool
	bgClr  color.Color
}

type labelWrap struct {
	txt    string
	txtClr color.Color
	txtBld bool
	bgClr  color.Color
	bind   binding.String
}

var (
	red         = color.NRGBA{R: 255, G: 0, B: 0, A: 125}
	green       = color.NRGBA{R: 0, G: 255, B: 0, A: 150}
	yellow      = color.NRGBA{R: 255, G: 255, B: 0, A: 125}
	grey        = color.NRGBA{R: 224, G: 224, B: 224, A: 255}
	pink        = color.NRGBA{R: 255, G: 179, B: 203, A: 1}
	azure       = color.NRGBA{R: 240, G: 255, B: 255, A: 255}
	beige       = color.NRGBA{R: 245, G: 245, B: 220, A: 255}
	bisque      = color.NRGBA{R: 255, G: 228, B: 196, A: 255}
	black       = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	white       = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	purple      = color.NRGBA{R: 55, G: 0, B: 179, A: 175}
	lightPurple = color.NRGBA{R: 55, G: 0, B: 179, A: 125}
	teal        = color.NRGBA{R: 3, G: 218, B: 198, A: 200}
)

var state string = STATE_OPERATE

func (app *application) screen() {
	var row4 fyne.CanvasObject

	a := ap.New()
	app.ap = a
	w := a.NewWindow(EME_TITLE)

	row4 = app.operatePage(w)

	w.SetContent(row4)
	w.Resize(fyne.NewSize(800, 400))

	w.Show()

	a.Run()
	fmt.Println("Exited")

}

func (app *application) setupPage(w fyne.Window) fyne.CanvasObject {

	t := &textWrap{txt: TITLE_SETUP, txtClr: white, txtBld: true, bgClr: purple}
	row0 := t.makeText()

	t = &textWrap{txt: TEXT_CURR_GRID, txtClr: black, txtBld: false, bgClr: white}
	latLabel := t.makeText()

	t.txt = TEXT_NEW_GRID
	lonLabel := t.makeText()

	row1 := container.New(layout.NewGridLayout(2), latLabel, lonLabel)

	t.txt = TEXT_GRID_VALUE + app.grid
	currGrid := t.makeText()

	//    lattitude := widget.NewEntry()
	//	lattitude.SetPlaceHolder(ENTER_LATTITUDE)
	newGrid := widget.NewEntry()
	newGrid.SetPlaceHolder(ENTER_GRID)
	row2 := container.New(layout.NewGridLayout(2), currGrid, newGrid)

	b := &buttonWrap{
		txt:    BUTTON_UPDATE_GRID,
		txtClr: black,
		txtBld: false,
		bgClr:  white,
		callBack: func() {
			app.updateGrid(newGrid.Text)
		},
	}
	enterSetUp := b.makeButton()

	row3 := container.New(layout.NewGridLayout(1), enterSetUp)
	row35 := seperator()

	t.txt = TEXT_PARK_AZIMUTH
	parkAzLabel := t.makeText()
	t.txt = TEXT_PARK_ELEVATION
	parkElLabel := t.makeText()

	row4 := container.New(layout.NewGridLayout(2), parkAzLabel, parkElLabel)

	parkAz := widget.NewEntry()
	parkAz.SetPlaceHolder(ENTER_PARK_AZIMUTH)
	parkEl := widget.NewEntry()
	parkEl.SetPlaceHolder(ENTER_PARK_ELEVATION)
	row5 := container.New(layout.NewGridLayout(2), parkAz, parkEl)

	b.txt = BUTTON_UPDATE_PARK
	b.callBack = func() {
		app.updatePark(parkAz.Text, parkEl.Text)
	}
	enterPark := b.makeButton()

	row6 := container.New(layout.NewGridLayout(1), enterPark)
	row65 := seperator()

	t.txt = TEXT_MAX_AZ
	maxAzLabel := t.makeText()
	t.txt = TEXT_MIN_AZ
	minAzLabel := t.makeText()
	t.txt = TEXT_MAX_EL
	maxElLabel := t.makeText()
	t.txt = TEXT_MIN_EL
	minElLabel := t.makeText()

	row7 := container.New(layout.NewGridLayout(4), maxAzLabel, minAzLabel, maxElLabel, minElLabel)

	maxAz := widget.NewEntry()
	maxAz.SetPlaceHolder(ENTER_MAX_AZ)
	minAz := widget.NewEntry()
	minAz.SetPlaceHolder(ENTER_MIN_AZ)
	maxEl := widget.NewEntry()
	maxEl.SetPlaceHolder(ENTER_MAX_EL)
	minEl := widget.NewEntry()
	minEl.SetPlaceHolder(ENTER_MIN_EL)

	row8 := container.New(layout.NewGridLayout(4), maxAz, minAz, maxEl, minEl)

	b.txt = BUTTON_UPDATE_MINMAX
	b.callBack = func() {
		app.updateMinMax(minAz.Text, maxAz.Text, minEl.Text, maxEl.Text)
	}
	enterMinMax := b.makeButton()

	row9 := container.New(layout.NewGridLayout(1), enterMinMax)
	row95 := seperator()

	rowN := app.basePage(w)
	setupGrid := container.NewBorder(nil, rowN, nil, nil, container.New(layout.NewVBoxLayout(),
		row0, row1, row2, row3, row35, row4, row5, row6, row65, row7, row8, row9, row95))
	return setupGrid

}

func (app *application) operatePage(w fyne.Window) fyne.CanvasObject {

	t := &textWrap{txt: TITLE_OPERATE, txtClr: white, txtBld: true, bgClr: purple}
	row0 := t.makeText()

	sunMoon := widget.NewRadioGroup([]string{MOON, SUN}, func(value string) {
		app.trackModeSelect(value)
	})
	sunMoon1 := container.New(layout.NewCenterLayout(), sunMoon)

	park := widget.NewButtonWithIcon(BUTTON_PARK, theme.HomeIcon(), func() {
		app.parkDish()
	})
	track := widget.NewButtonWithIcon(BUTTON_TRACK, theme.HomeIcon(), func() {
		app.pushedTrack()
	})
	row1 := container.New(layout.NewGridLayout(3), sunMoon1, track, park)

	row15 := seperator()
    
	t.txt = TEXT_TRACKING + TEXT_MOON
	t.bgClr = lightPurple
	tracking := t.makeText()
	row3 := container.New(layout.NewGridLayout(1), tracking)

	t = &textWrap{txt: TEXT_CURR_AZ, txtClr: black, txtBld: false, bgClr: white}
	currAzLabel := t.makeText()
	t.txt = TEXT_CURR_EL
	currElLabel := t.makeText()
	row4 := container.New(layout.NewGridLayout(2), currAzLabel, currElLabel)

	l := &labelWrap{
		txt:    fmt.Sprintf("%5.2f", app.currAz),
		txtClr: black,
		txtBld: false,
		bgClr:  white,
	}
	l.bind = app.azBind
	currAz := l.makeLabel()

	l.txt = fmt.Sprintf("%5.2f", app.currEl)
	l.bind = app.elBind
	currEl := l.makeLabel()
	row5 := container.New(layout.NewGridLayout(2), currAz, currEl)

	row55 := seperator()

	rowN := app.basePage(w)
	operateGrid := container.NewBorder(nil, rowN, nil, nil, container.New(layout.NewVBoxLayout(),
		row0, row1, row15, row3, row4, row5, row55))
	return operateGrid
}

func (app *application) pointPage(w fyne.Window) fyne.CanvasObject {

	t := &textWrap{txt: TITLE_POINT, txtClr: white, txtBld: true, bgClr: purple}
	row0 := t.makeText()

	t = &textWrap{txt: TEXT_TARGET_AZ, txtClr: black, txtBld: false, bgClr: white}
	targetAzLabel := t.makeText()
	t.txt = TEXT_TARGET_EL
	targetElLabel := t.makeText()

	row1 := container.New(layout.NewGridLayout(2), targetAzLabel, targetElLabel)

	targetAz := widget.NewEntry()
	targetAz.SetPlaceHolder(ENTER_TARGET_AZ)
	targetEl := widget.NewEntry()
	targetEl.SetPlaceHolder(ENTER_TARGET_EL)
	row2 := container.New(layout.NewGridLayout(2), targetAz, targetEl)

	b := &buttonWrap{
		txt:    BUTTON_UPDATE_TARGET,
		txtClr: black,
		txtBld: false,
		bgClr:  white,
		callBack: func() {
			app.updateTarget(targetAz.Text, targetEl.Text)
		},
	}
	enterTarget := b.makeButton()
	row3 := container.New(layout.NewGridLayout(1), enterTarget)

	row35 := seperator()

	t.txt = TEXT_ADJ_SIZE
	adj := t.makeText()
	row36 := container.New(layout.NewGridLayout(1), adj)

	t.txt = TEXT_CURR_AZ
	currAzLabel := t.makeText()
	t.txt = TEXT_CURR_EL
	currElLabel := t.makeText()

	row4 := container.New(layout.NewGridLayout(2), currAzLabel, currElLabel)

	l := &labelWrap{
		txt:    fmt.Sprintf("%5.2f", app.currAz),
		txtClr: black,
		txtBld: false,
		bgClr:  white,
	}
	l.bind = app.azBind
	currAz := l.makeLabel()

	l.txt = fmt.Sprintf("%5.2f", app.currEl)
	l.bind = app.elBind
	currEl := l.makeLabel()
	row5 := container.New(layout.NewGridLayout(2), currAz, currEl)

	b.txt = BUTTON_ADJ_UP
	b.callBack = app.adjustUp
	adjUp := b.makeButton()
	b.txt = BUTTON_ADJ_DN
	b.callBack = app.adjustDn
	adjDn := b.makeButton()
	b.txt = BUTTON_ADJ_RIGHT
	b.callBack = app.adjustRight
	adjRight := b.makeButton()
	b.txt = BUTTON_ADJ_LEFT
	b.callBack = app.adjustLeft
	adjLeft := b.makeButton()
	row6 := container.New(layout.NewGridLayout(4), adjRight, adjLeft, adjUp, adjDn)

	b.txt = BUTTON_RECALIBRATE
	b.callBack = func() { fmt.Println("Recalibrate") }
	reCalib := b.makeButton()
	row7 := container.New(layout.NewGridLayout(1), reCalib)

	row75 := seperator()

	rowN := app.basePage(w)
	setupGrid := container.NewBorder(nil, rowN, nil, nil, container.New(layout.NewVBoxLayout(),
		row0, row1, row2, row3, row35, row36, row4, row5, row6, row7, row75))
	return setupGrid

}

func (app *application) manualPage(w fyne.Window) fyne.CanvasObject {

	t := &textWrap{txt: TITLE_MANUAL, txtClr: white, txtBld: true, bgClr: purple}
	row0 := t.makeText()

	t = &textWrap{txt: TEXT_TARGET_AZ, txtClr: black, txtBld: false, bgClr: white}
	targetAzLabel := t.makeText()
	t.txt = TEXT_TARGET_EL
	targetElLabel := t.makeText()

	row1 := container.New(layout.NewGridLayout(2), targetAzLabel, targetElLabel)

	targetAz := widget.NewEntry()
	targetAz.SetPlaceHolder(ENTER_TARGET_AZ)
	targetEl := widget.NewEntry()
	targetEl.SetPlaceHolder(ENTER_TARGET_EL)
	row2 := container.New(layout.NewGridLayout(2), targetAz, targetEl)

	b := &buttonWrap{
		txt:    BUTTON_UPDATE_TARGET,
		txtClr: black,
		txtBld: false,
		bgClr:  white,
		callBack: func() {
			app.updateTarget(targetAz.Text, targetEl.Text)
		},
	}
	enterTarget := b.makeButton()

	row3 := container.New(layout.NewGridLayout(1), enterTarget)
	row35 := seperator()

	t.txt = TEXT_CURR_AZ
	currAzLabel := t.makeText()
	t.txt = TEXT_CURR_EL
	currElLabel := t.makeText()

	row4 := container.New(layout.NewGridLayout(2), currAzLabel, currElLabel)

	l := &labelWrap{
		txt:    fmt.Sprintf("%5.2f", app.currAz),
		txtClr: black,
		txtBld: false,
		bgClr:  white,
	}
	l.bind = app.azBind
	currAz := l.makeLabel()

	l.txt = fmt.Sprintf("%5.2f", app.currEl)
	l.bind = app.elBind
	currEl := l.makeLabel()
	row5 := container.New(layout.NewGridLayout(2), currAz, currEl)

	b.txt = BUTTON_ADJ_UP
	b.callBack = app.adjustUp
	adjUp := b.makeButton()
	b.txt = BUTTON_ADJ_DN
	b.callBack = app.adjustDn
	adjDn := b.makeButton()
	b.txt = BUTTON_ADJ_RIGHT
	b.callBack = app.adjustRight
	adjRight := b.makeButton()
	b.txt = BUTTON_ADJ_LEFT
	b.callBack = app.adjustLeft
	adjLeft := b.makeButton()
	row6 := container.New(layout.NewGridLayout(4), adjRight, adjLeft, adjUp, adjDn)

	row65 := seperator()

	rowN := app.basePage(w)
	setupGrid := container.NewBorder(nil, rowN, nil, nil, container.New(layout.NewVBoxLayout(),
		row0, row1, row2, row3, row35, row4, row5, row6, row65))
	return setupGrid
}

func (app *application) basePage(w fyne.Window) fyne.CanvasObject {
	var bOper *buttonWrap
	var bMan *buttonWrap
	var bPoint *buttonWrap
	var bSet *buttonWrap

	bOper = &buttonWrap{
		txt:    "Operate",
		txtClr: black,
		txtBld: false,
		bgClr:  opColor(),
		callBack: func() {
			bMan.bgClr = grey
			bPoint.bgClr = grey
			bSet.bgClr = grey
			state = STATE_OPERATE
			fmt.Println("Operate")
			operateGrid := app.operatePage(w)
			w.SetContent(operateGrid)
			w.Show()
		},
	}
	operate := bOper.makeButton()

	bMan = &buttonWrap{
		txt:    "Manual",
		txtClr: black,
		txtBld: false,
		bgClr:  manColor(),
		callBack: func() {
			state = STATE_MANUAL
			fmt.Println("Manual")
			manualGrid := app.manualPage(w)
			w.SetContent(manualGrid)
			w.Show()
		},
	}
	manual := bMan.makeButton()

	bPoint = &buttonWrap{
		txt:    "Point",
		txtClr: black,
		txtBld: false,
		bgClr:  pointColor(),
		callBack: func() {
			state = STATE_POINT
			fmt.Println("Point")
			pointGrid := app.pointPage(w)
			w.SetContent(pointGrid)
			w.Show()
		},
	}
	point := bPoint.makeButton()

	bSet = &buttonWrap{
		txt:    "Set Up",
		txtClr: black,
		txtBld: false,
		bgClr:  setupColor(),
		callBack: func() {
			state = STATE_SETUP
			fmt.Println("Set Up")
			setupGrid := app.setupPage(w)
			w.SetContent(setupGrid)
			w.Show()
		},
	}
	setup := bSet.makeButton()

	row4 := container.New(layout.NewGridLayout(4), operate, manual, point, setup)
	row4 = colorize(row4, black)
	fmt.Println("State: ", state)
	return row4
}
