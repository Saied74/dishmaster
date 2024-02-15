package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	ap "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"math"
	"time"
)

type buttonWrap struct {
	txt      string
	txtClr   color.Color
	txtBld   bool
	bgClr    color.Color
	callBack func()
}

type textWrap struct {
	txt     string
	txtClr  color.Color
	txtSize float32
	txtBld  bool
	bgClr   color.Color
	bind    binding.String
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

const (
	windowWidth  = 430
	windowHeight = 515
)

func (app *application) screen() {
	var row4 fyne.CanvasObject

	a := ap.New()
	app.ap = a
	w := a.NewWindow(EME_TITLE)

	row4 = app.operatePage(w)

	w.SetContent(row4)
	w.Resize(fyne.NewSize(windowWidth, windowHeight)) //850, 700))

	w.Show()

	a.Run()
	fmt.Println("Exited")

}

func (app *application) setupPage(w fyne.Window) fyne.CanvasObject {

	t := &textWrap{
		txt:     TITLE_SETUP,
		txtClr:  white,
		txtSize: SIZE_PAGE_TITLE,
		txtBld:  true,
		bgClr:   purple,
	}
	row0 := t.makeText()

	t = &textWrap{txt: TEXT_CURR_GRID, txtClr: black, txtBld: false, bgClr: white}
	latLabel := t.makeText()

	t.txt = TEXT_NEW_GRID
	lonLabel := t.makeText()

	row1 := container.New(layout.NewGridLayout(2), latLabel, lonLabel)
	l := &labelWrap{
		txtClr: black,
		txtBld: false,
		bgClr:  white,
	}
	l.txt = fmt.Sprintf("%s%s", TEXT_GRID_VALUE, app.grid)
	l.bind = app.gridBind
	currGrid := l.makeLabel()

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

	l = &labelWrap{
		txtClr: black,
		txtBld: false,
		bgClr:  white,
	}
	l.txt = fmt.Sprintf("%s%5.2f", TEXT_CURRENT_VALUE, app.parkAz)
	l.bind = app.parkAzBind
	currParkAz := l.makeLabel()
	l.txt = fmt.Sprintf("%s%5.2f", TEXT_CURRENT_VALUE, app.parkEl)
	l.bind = app.parkElBind
	currParkEl := l.makeLabel()
	row45 := container.New(layout.NewGridLayout(2), currParkAz, currParkEl)

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

	l.txt = fmt.Sprintf("%s%5.2f", TEXT_CURRENT_VALUE, app.maxAz)
	l.bind = app.maxAzBind
	currMaxAz := l.makeLabel()
	l.txt = fmt.Sprintf("%s%5.2f", TEXT_CURRENT_VALUE, app.minAz)
	l.bind = app.minAzBind
	currMinAz := l.makeLabel()
	l.txt = fmt.Sprintf("%s%5.2f", TEXT_CURRENT_VALUE, app.maxEl)
	l.bind = app.maxElBind
	currMaxEl := l.makeLabel()
	l.txt = fmt.Sprintf("%s%5.2f", TEXT_CURRENT_VALUE, app.minEl)
	l.bind = app.minElBind
	currMinEl := l.makeLabel()
	row75 := container.New(layout.NewGridLayout(4), currMaxAz, currMinAz, currMaxEl, currMinEl)

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
		row0, row1, row2, row3, row35, row4, row45, row5, row6, row65, row7, row75, row8, row9,
		row95))
	return setupGrid

}

func (app *application) operatePage(w fyne.Window) fyne.CanvasObject {

	t := &textWrap{
		txt:     TITLE_OPERATE,
		txtClr:  white,
		txtSize: SIZE_PAGE_TITLE,
		txtBld:  true,
		bgClr:   purple,
	}
	row0 := t.makeText()

	sunMoon := widget.NewRadioGroup([]string{MOON, SUN}, func(value string) {
		app.trackModeSelect(value)
	})
	sunMoon1 := container.New(layout.NewCenterLayout(), sunMoon)
	parkSign, err := fyne.LoadResourceFromPath("./assets/park2.jpg")
	if err != nil {
		log.Printf("Failed to load park sign: %v", err)
	}
	park := widget.NewButtonWithIcon(BUTTON_PARK, parkSign, func() {
		app.pushedPark()
	})
	trackSign, err := fyne.LoadResourceFromPath("./assets/track2.jpg")
	if err != nil {
		log.Printf("Failed to load track sign: %v", err)
	}
	track := widget.NewButtonWithIcon(BUTTON_TRACK, trackSign, func() {
		app.pushedTrack()
	})
	redCircle, err := fyne.LoadResourceFromPath("./assets/red_circle.png")
	if err != nil {
		log.Printf("Failed to load red circle: %v", err)
	}
	stop := widget.NewButtonWithIcon(BUTTON_STOP, redCircle, func() {
		app.pushedStop()
	})

	row1 := container.New(layout.NewGridLayout(4), sunMoon1, track, stop, park)

	row15 := seperator()

	l := &labelWrap{
		txtClr: black,
		txtBld: true,
		bgClr:  white,
	}
	l.bind = app.modeBind
	switch app.state {
	case TRACKING_SUN:
		l.txt = "Tracking the Sun"
	case TRACKING_MOON:
		l.txt = "Tracking the Moon"
	case PARKED:
		l.txt = "Parked"
	case IDLE:
		l.txt = "Idle"
	}
	opMode := l.makeLabel()
	row3 := container.New(layout.NewGridLayout(1), opMode)

	//var smallText float32 = 12.0

	t = &textWrap{txt: "Target", txtClr: black, txtBld: false, bgClr: white}
	currAzLabel := t.makeText()
	//currAzLabel.TextSize = smallText
	t.txt = "Current"          //TEXT_CURRENT_AZ
	azPosLabel := t.makeText() //canvas.NewText(t.txt, t.txtClr)
	//azPosLabel.TextSize = smallText
	t.txt = "Difference"        //TEXT_TARGET_EL
	currElLabel := t.makeText() //canvas.NewText(t.txt, t.txtClr)
	//currElLabel.TextSize = smallText

	t.txt = ""
	blank := t.makeText() //canvas.NewText(t.txt, t.txtClr)
	//blank.TextSize = smallText

	row4 := container.New(layout.NewGridLayout(4), blank, currAzLabel, azPosLabel, currElLabel) //, elPosLabel)

	t.txt = "Azimuth"
	azRowLabel := t.makeText() //canvas.NewText(t.txt, t.txtClr)
	//azRowLabel.TextSize = smallText

	t.txt = "Elevation"
	elRowLabel := t.makeText() //canvas.NewText(t.txt, t.txtClr)
	// elRowLabel.TextSize = smallText

	l = &labelWrap{
		txt:    fmt.Sprintf("%5.2f", app.currAz),
		txtClr: black,
		txtBld: false,
		bgClr:  white,
	}
	l.bind = app.azBind
	currAz := l.makeLabel()

	l.txt = fmt.Sprintf("%5.2f", app.azPosition)
	l.bind = app.azPosBind
	azPosition := l.makeLabel()

	l.txt = fmt.Sprintf("%5.2f", app.currAz-app.azPosition)
	l.bind = app.azDiffBind
	azDiff := l.makeLabel()

	l.txt = fmt.Sprintf("%5.2f", app.currEl)
	l.bind = app.elBind
	currEl := l.makeLabel()

	l.txt = fmt.Sprintf("%5.2f", app.elPosition)
	l.bind = app.elPosBind
	elPosition := l.makeLabel()

	l.txt = fmt.Sprintf("%5.2f", app.currEl-app.elPosition)
	l.bind = app.elDiffBind
	elDiff := l.makeLabel()

	row5 := container.New(layout.NewGridLayout(4), azRowLabel, currAz, azPosition, azDiff)
	row51 := container.New(layout.NewGridLayout(4), elRowLabel, currEl, elPosition, elDiff)

	hashMarksa := app.makeScale("az")

	//define the dial line and locate it
	la := canvas.NewLine(red)
	la.StrokeWidth = 3
	la.Position1 = fyne.Position{float32(app.sDa.centerX), float32(app.sDa.centerY)}

	go func() {
		for {
			radThetaA := ((2.0 * math.Pi) * (app.azPosition - 90.0)) / 360.0
			endXa := innerX*math.Cos(radThetaA) + app.sDa.centerX
			endYa := innerY*math.Sin(radThetaA) + app.sDa.centerY
			lp2a := fyne.Position{float32(endXa), float32(endYa)}
			la.Position2 = lp2a
			canvas.Refresh(la)
			time.Sleep(time.Duration(500) * time.Millisecond)
		}
	}()

	hashMarksa = append(hashMarksa, la)
	ita := container.NewWithoutLayout(hashMarksa...)

	hashMarkse := app.makeScale("el")

	//define the dial line and locate it
	le := canvas.NewLine(red)
	le.StrokeWidth = 3
	le.Position1 = fyne.Position{float32(app.sDe.centerX), float32(app.sDe.centerY)}

	go func() {
		for {
			radThetaE := ((2.0 * math.Pi) * (360.0 - app.elPosition)) / 360.0
			endXe := innerX*math.Cos(radThetaE) + app.sDe.centerX
			endYe := innerY*math.Sin(radThetaE) + app.sDe.centerY
			le.Position2 = fyne.Position{float32(endXe), float32(endYe)}
			canvas.Refresh(le)
			time.Sleep(time.Duration(500) * time.Millisecond)
		}
	}()

	hashMarkse = append(hashMarkse, le)
	ite := container.NewWithoutLayout(hashMarkse...)

	row54 := container.New(layout.NewGridLayout(2), ita, ite)

	row55 := seperator()

	//	t = &textWrap{txt: TEXT_TARGET_AZ, txtClr: black, txtBld: false, bgClr: white}
	//	targetAzLabel := t.makeText()
	//	t.txt = TEXT_TARGET_EL
	//	targetElLabel := t.makeText()
	//
	//	row6 := container.New(layout.NewGridLayout(2), targetAzLabel, targetElLabel)
	//
	//	targetAz := widget.NewEntry()
	//	targetAz.SetPlaceHolder(ENTER_TARGET_AZ)
	//	targetEl := widget.NewEntry()
	//	targetEl.SetPlaceHolder(ENTER_TARGET_EL)
	//	row7 := container.New(layout.NewGridLayout(2), targetAz, targetEl)
	//
	//	b := &buttonWrap{
	//		txt:    BUTTON_UPDATE_TARGET,
	//		txtClr: black,
	//		txtBld: false,
	//		bgClr:  white,
	//		callBack: func() {
	//			app.updateTarget(targetAz.Text, targetEl.Text)
	//		},
	//	}
	//	enterTarget := b.makeButton()
	//	b.txt = BUTTON_ADJ_UP
	//	b.callBack = app.adjustUp
	//	adjUp := b.makeButton()
	//	b.txt = BUTTON_ADJ_DN
	//	b.callBack = app.adjustDn
	//	adjDn := b.makeButton()
	//	b.txt = BUTTON_ADJ_RIGHT
	//	b.callBack = app.adjustRight
	//	adjRight := b.makeButton()
	//	b.txt = BUTTON_ADJ_LEFT
	//	b.callBack = app.adjustLeft
	//	adjLeft := b.makeButton()
	//	b.txt = BUTTON_RECALIBRATE
	//	b.callBack = func() {
	//		app.recalibrate(targetAz.Text, targetEl.Text)
	//	}
	//	reCalib := b.makeButton()
	//
	//	row8 := container.New(layout.NewGridLayout(6), enterTarget, adjRight, adjLeft, adjUp, adjDn, reCalib)

	rowN := app.basePage(w)
	operateGrid := container.NewBorder(nil, rowN, nil, nil, container.New(layout.NewVBoxLayout(),
		row0, row1, row15, row3, row4, row5, row51, row55, row54)) //took out rows 6, 7, and 8
	return operateGrid
}

func (app *application) pointPage(w fyne.Window) fyne.CanvasObject {

	t := &textWrap{
		txt:     TITLE_POINT,
		txtClr:  white,
		txtSize: SIZE_PAGE_TITLE,
		txtBld:  true,
		bgClr:   purple,
	}
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
	b.callBack = func() {
		app.recalibrate(targetAz.Text, targetEl.Text)
	}

	reCalib := b.makeButton()
	row7 := container.New(layout.NewGridLayout(1), reCalib)

	row75 := seperator()

	rowN := app.basePage(w)
	setupGrid := container.NewBorder(nil, rowN, nil, nil, container.New(layout.NewVBoxLayout(),
		row0, row1, row2, row3, row35, row36, row4, row5, row6, row7, row75))
	return setupGrid

}

//
//func (app *application) manualPage(w fyne.Window) fyne.CanvasObject {
//
//	t := &textWrap{
//		txt:     TITLE_MANUAL,
//		txtClr:  white,
//		txtSize: SIZE_PAGE_TITLE,
//		txtBld:  true,
//		bgClr:   purple,
//	}
//	row0 := t.makeText()
//
//	t = &textWrap{txt: TEXT_TARGET_AZ, txtClr: black, txtBld: false, bgClr: white}
//	targetAzLabel := t.makeText()
//	t.txt = TEXT_TARGET_EL
//	targetElLabel := t.makeText()
//
//	row1 := container.New(layout.NewGridLayout(2), targetAzLabel, targetElLabel)
//
//	targetAz := widget.NewEntry()
//	targetAz.SetPlaceHolder(ENTER_TARGET_AZ)
//	targetEl := widget.NewEntry()
//	targetEl.SetPlaceHolder(ENTER_TARGET_EL)
//	row2 := container.New(layout.NewGridLayout(2), targetAz, targetEl)
//
//	b := &buttonWrap{
//		txt:    BUTTON_UPDATE_TARGET,
//		txtClr: black,
//		txtBld: false,
//		bgClr:  white,
//		callBack: func() {
//			app.updateTarget(targetAz.Text, targetEl.Text)
//		},
//	}
//	enterTarget := b.makeButton()
//
//	row3 := container.New(layout.NewGridLayout(1), enterTarget)
//	row35 := seperator()
//
//	t.txt = TEXT_ADJ_SIZE
//	adj := t.makeText()
//	row36 := container.New(layout.NewGridLayout(1), adj)
//
//	t.txt = TEXT_CURR_AZ
//	currAzLabel := t.makeText()
//	t.txt = TEXT_CURR_EL
//	currElLabel := t.makeText()
//
//	row4 := container.New(layout.NewGridLayout(2), currAzLabel, currElLabel)
//
//	l := &labelWrap{
//		txt:    fmt.Sprintf("%5.2f", app.currAz),
//		txtClr: black,
//		txtBld: false,
//		bgClr:  white,
//	}
//	l.bind = app.azBind
//	currAz := l.makeLabel()
//
//	l.txt = fmt.Sprintf("%5.2f", app.currEl)
//	l.bind = app.elBind
//	currEl := l.makeLabel()
//	row5 := container.New(layout.NewGridLayout(2), currAz, currEl)
//
//	b.txt = BUTTON_ADJ_UP
//	b.callBack = app.adjustUp
//	adjUp := b.makeButton()
//	b.txt = BUTTON_ADJ_DN
//	b.callBack = app.adjustDn
//	adjDn := b.makeButton()
//	b.txt = BUTTON_ADJ_RIGHT
//	b.callBack = app.adjustRight
//	adjRight := b.makeButton()
//	b.txt = BUTTON_ADJ_LEFT
//	b.callBack = app.adjustLeft
//	adjLeft := b.makeButton()
//	row6 := container.New(layout.NewGridLayout(4), adjRight, adjLeft, adjUp, adjDn)
//
//	row65 := seperator()
//
//	rowN := app.basePage(w)
//	setupGrid := container.NewBorder(nil, rowN, nil, nil, container.New(layout.NewVBoxLayout(),
//		row0, row1, row2, row3, row35, row36, row4, row5, row6, row65))
//	return setupGrid
//}

func (app *application) basePage(w fyne.Window) fyne.CanvasObject {
	var bOper *buttonWrap
	//	var bMan *buttonWrap
	var bPoint *buttonWrap
	var bSet *buttonWrap

	bOper = &buttonWrap{
		txt:    "Operate",
		txtClr: black,
		txtBld: false,
		bgClr:  opColor(),
		callBack: func() {
			//			bMan.bgClr = grey
			//			bPoint.bgClr = grey
			bSet.bgClr = grey
			state = STATE_OPERATE
			fmt.Println("Operate")
			operateGrid := app.operatePage(w)
			w.SetContent(operateGrid)
			w.Resize(fyne.NewSize(windowWidth, windowHeight)) //850, 700))
			w.Show()
		},
	}
	operate := bOper.makeButton()

	//	bMan = &buttonWrap{
	//		txt:    "Manual",
	//		txtClr: black,
	//		txtBld: false,
	//		bgClr:  manColor(),
	//		callBack: func() {
	//			state = STATE_MANUAL
	//			fmt.Println("Manual")
	//			manualGrid := app.manualPage(w)
	//			w.SetContent(manualGrid)
	//			w.Show()
	//		},
	//	}
	//	manual := bMan.makeButton()
	//
	bPoint = &buttonWrap{
		txt:    "Calibrate",
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

	row4 := container.New(layout.NewGridLayout(3), operate, point, setup) //manual, point, setup)
	row4 = colorize(row4, black)
	fmt.Println("State: ", state)
	return row4
}
