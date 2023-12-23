package main

import (
	//	"fmt"
	"image/color"
	"log"
	//	"runtime/debug"
	//	"strconv"
	//	"time"
	"fyne.io/fyne/v2"
	//	"strconv"

	//	ap "fyne.io/fyne/v2/app"
	//  "fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	//	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	EME_TITLE = "EME Dish Controller"

	STATE_OPERATE = "operate"
	STATE_MANUAL  = "manual"
	STATE_POINT   = "point"
	STATE_SETUP   = "setup"

	TITLE_SETUP = "Set Up"

	TEXT_CURR_GRID     = "Current Grid Square"
	TEXT_NEW_GRID      = "Enter New Grid"
	TEXT_GRID_VALUE    = "You are on grid: "
	ENTER_GRID         = "Enter new grid..."
	BUTTON_UPDATE_GRID = "Update Grid"

	TEXT_PARK_AZIMUTH    = "Park Azimuth"
	TEXT_PARK_ELEVATION  = "Park Elevation"
	ENTER_PARK_AZIMUTH   = "Enter park azimuth..."
	ENTER_PARK_ELEVATION = "Enter park elevaton..."
	BUTTON_UPDATE_PARK   = "Update Park Az/El"

	TEXT_MAX_AZ          = "Max Azimuth"
	TEXT_MIN_AZ          = "Min Azimuth"
	TEXT_MAX_EL          = "Max Elevation"
	TEXT_MIN_EL          = "Min Elevation"
	ENTER_MAX_AZ         = "Enter Max Az..."
	ENTER_MIN_AZ         = "Enter Min Az..."
	ENTER_MAX_EL         = "Enter Max El..."
	ENTER_MIN_EL         = "Enter Min El..."
	BUTTON_UPDATE_MINMAX = "Update Min/Max Az/El"

	TITLE_OPERATE = "Operate"

	SUN  = "Sun"
	MOON = "Moon"

	BUTTON_TRACK  = "Track"
	BUTTON_PARK   = "Park"
	TEXT_TRACKING = "Tracking: "
	TEXT_MOON     = "The Moon"
	TEXT_CURR_AZ  = "Current Azimuth"
	TEXT_CURR_EL  = "Current Elevation"

	TITLE_POINT = "Point and Adjust"

	TEXT_TARGET_AZ       = "Target Azimuth"
	TEXT_TARGET_EL       = "Target Elevation"
	ENTER_TARGET_AZ      = "Enter target azimuth..."
	ENTER_TARGET_EL      = "Enter target elevation..."
	BUTTON_UPDATE_TARGET = "Update Target Az/El"

	TEXT_ADJ_SIZE      = "All adjustments are by 0.5 degree increments"
	BUTTON_ADJ_UP      = "Adjust Up"
	BUTTON_ADJ_DN      = "Adjust Down"
	BUTTON_ADJ_RIGHT   = "Adjust Right"
	BUTTON_ADJ_LEFT    = "Adjust Left"
	BUTTON_RECALIBRATE = "Recalibrate"

	TITLE_MANUAL = "Manual"

	TITLE_ERR = "ERROR"
	BUTTON_OK = "OK"
)

func seperator() fyne.CanvasObject {
	s := canvas.NewRectangle(black)
	return container.New(layout.NewMaxLayout(), s)
}

func colorize(item *fyne.Container, clr color.Color) *fyne.Container {
	bgColor := canvas.NewRectangle(clr)
	return container.New(layout.NewMaxLayout(), bgColor, item)
}

func (b *buttonWrap) makeButton() *fyne.Container {
	bg := canvas.NewRectangle(b.bgClr)
	button := widget.NewButton(b.txt, b.callBack)
	button1 := container.New(layout.NewCenterLayout(), button)
	return container.New(layout.NewMaxLayout(), bg, button1)
}

func (t *textWrap) makeText() *fyne.Container {
	bg := canvas.NewRectangle(t.bgClr)
	txtArea := canvas.NewText(t.txt, t.txtClr)
	if t.txtBld {
		txtArea.TextStyle.Bold = true
	}
	con := container.New(layout.NewCenterLayout(), txtArea)
	return container.New(layout.NewMaxLayout(), bg, con)
}

func (l *labelWrap) makeLabel() *fyne.Container {
	bg := canvas.NewRectangle(l.bgClr)

	err := l.bind.Set(l.txt)
	if err != nil {
		log.Fatalf("in creating bound label: %v", err)
	}
	label := widget.NewLabelWithData(l.bind)
	con := container.New(layout.NewCenterLayout(), label)
	return container.New(layout.NewMaxLayout(), bg, con)
}

func opColor() color.Color {
	if state == STATE_OPERATE {
		return teal
	}
	return grey
}

func manColor() color.Color {
	if state == STATE_MANUAL {
		return teal
	}
	return grey
}

func pointColor() color.Color {
	if state == STATE_POINT {
		return teal
	}
	return grey
}

func setupColor() color.Color {
	if state == STATE_SETUP {
		return teal
	}
	return grey
}
