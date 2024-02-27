package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"image/color"
	"log"
	"math"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	EME_TITLE = "EME Dish Controller Version v0.7"

	STATE_OPERATE = "operate"
	STATE_MANUAL  = "manual"
	STATE_POINT   = "point"
	STATE_SETUP   = "setup"

	TITLE_SETUP = "Set Up"

	TEXT_CURRENT_VALUE = "Current Value: "

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

	BUTTON_TRACK    = "Track"
	BUTTON_PARK     = "Park"
	BUTTON_STOP     = "Stop"
	TEXT_TRACKING   = "Tracking: "
	TEXT_MOON       = "The Moon"
	TEXT_SUN        = "The Sun"
	TEXT_IDLE       = "Not Tracking"
	TEXT_CURRENT_AZ = "Current Azimuth"
	TEXT_CURRENT_EL = "Current Elevation"
	TEXT_CURR_AZ    = "Current Azimuth"
	TEXT_CURR_EL    = "Current Elevation"

	TITLE_POINT = "Point and Adjust"

	TEXT_TARGET_AZ       = "Target Azimuth"
	TEXT_TARGET_EL       = "Target Elevation"
	ENTER_TARGET_AZ      = "Enter target azimuth..."
	ENTER_TARGET_EL      = "Enter target elevation..."
	BUTTON_UPDATE_TARGET = "Update Target Az/El"

	TEXT_ADJ_SIZE      = "All adjustments are by 0.5 degree increments"
	BUTTON_ADJ_UP      = "Up"
	BUTTON_ADJ_DN      = "Dn"
	BUTTON_ADJ_RIGHT   = "CW"
	BUTTON_ADJ_LEFT    = "CCW"
	BUTTON_RECALIBRATE = "Recalibrate"

	TITLE_MANUAL = "Manual"

	TITLE_ERR = "ERROR"
	BUTTON_OK = "OK"

	SIZE_PAGE_TITLE = 18.0
)

const (
	//	centerX   = 400.0 //center of the scale geometry
	//	centerY   = 250.0
	innerX    = 75.0  //100.0 //140.0 //hashmark start point from the center
	innerY    = 75.0  //100.0 //140.0
	outerX    = 82.0  //110.0 //150.0 //hasmark end point from the center
	outerY    = 82.0  //110.0 //150.0
	txtRadX   = 105.0 //140.0 //180.0 //text distance from the center
	txtRadY   = 105.0 //140.0 //180.0
	limitRadX = 120.0 //160.0 //205.0 //for the upper and lower limit labels
	limitRadY = 120.0 //160.0 //205.0
)

type scaleData struct {
	ll         float64
	ul         float64
	lowerLimit float64
	upperLimit float64
	centerX    float64
	centerY    float64
	endX       float64
	endY       float64
}

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
	if t.txtSize == 0 {
		t.txtSize = 14
	}
	if t.bind != nil {
		err := t.bind.Set(t.txt)
		if err != nil {
			log.Fatalf("creating bound text failed: %v", err)
		}
	}
	txtArea.TextSize = t.txtSize
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
	label.TextStyle.Bold = l.txtBld
	label.Resize(fyne.Size{3.0, 3.0})
	con := container.New(layout.NewCenterLayout(), label)
	return container.New(layout.NewMaxLayout(), bg, con)
}

func opColor() color.Color {
	if state == STATE_OPERATE {
		return purple //teal
	}
	return grey
}

func manColor() color.Color {
	if state == STATE_MANUAL {
		return purple //teal
	}
	return grey
}

func pointColor() color.Color {
	if state == STATE_POINT {
		return purple //teal
	}
	return grey
}

func setupColor() color.Color {
	if state == STATE_SETUP {
		return purple //teal
	}
	return grey
}

func calcLetter(l int) (w, h float64) {
	if l == 0 {
		return -6., -8.0
	}
	if l < 100 {
		return -12., -8.
	}
	return -15, -5 //-18., -8.

}

func roundUp(x float64) float64 {
	if int(x)%10 == 0 {
		return x
	} else {
		return 10. * math.Round((x/10.)+0.5)
	}
}

func checkLimit(azel string, i, endPoint float64) bool {
	switch azel {
	case "az":
		if i <= endPoint {
			return true
		}
	case "el":
		if i >= endPoint {
			return true
		}
	default:
		log.Printf("program bug, did not ask for az or el, asked for %s", azel)
	}
	return false
}

func (app *application) makeScale(azel string) []fyne.CanvasObject {

	var ll, ul, txtX, txtY float64
	var txt *canvas.Text
	var inc float64
	hashMarks := []fyne.CanvasObject{}
	sD := scaleData{}
	var smallText float32 = 10.0

	switch azel {
	case "az":
		sD.ll = app.minAz
		sD.ul = app.maxAz
		sD.lowerLimit = sD.ll - 90.0
		sD.upperLimit = sD.ul - 90.0
		sD.centerX = app.sDa.centerX
		sD.centerY = app.sDa.centerY
		inc = 10.0 //azimuth rotates clockwise
	case "el":
		sD.ll = app.minEl
		sD.ul = app.maxEl
		sD.upperLimit = 360.0 - sD.ul
		sD.lowerLimit = 360.0 - sD.ll
		sD.centerX = app.sDe.centerX
		sD.centerY = app.sDe.centerY
		inc = -10.0 //elevation rotates counterclockwise
	}

	//locate the first hash mark at the lowerLimit
	radAlpha := (2.0 * math.Pi) * (sD.lowerLimit / 360.0)
	cA := math.Cos(radAlpha)
	sA := math.Sin(radAlpha)
	hashX := innerX*cA + sD.centerX
	hashY := innerY*sA + sD.centerY
	endX := outerX*cA + sD.centerX
	endY := outerY*sA + sD.centerY
	ln := canvas.NewLine(red)
	ln.StrokeWidth = 2
	ln.Position1 = fyne.Position{float32(hashX), float32(hashY)}
	ln.Position2 = fyne.Position{float32(endX), float32(endY)}
	hashMarks = append(hashMarks, ln)

	//label the first hash mark - note that it is 10 units further out than the rest
	ll, ul = calcLetter(int(sD.ll)) //to compensate for text width
	txtX = limitRadX*cA + sD.centerX + ll
	txtY = limitRadY*sA + sD.centerY + ul
	txt = canvas.NewText(fmt.Sprintf("%4.0f", sD.ll), color.Black)
	txt.TextSize = smallText
	txt.Move(fyne.Position{float32(txtX), float32(txtY)})
	txt.TextStyle.Bold = true //just for the limit labels
	hashMarks = append(hashMarks, txt)

	//calculate the bounds of the hash marks
	var startPoint float64

	startPoint = roundUp(sD.lowerLimit)
	endPoint := 10 * math.Round(sD.upperLimit/10.) //truncate to 10s

	//plot the hash marks and labels around the scale excluding the limit hashmarks
	for i := startPoint; checkLimit(azel, i, endPoint); i += inc {
		switch azel {
		case "az":
			if i > sD.upperLimit {
				continue
			}
		case "el":
			if i > sD.lowerLimit {
				continue
			}
		default:
			log.Printf("program bug, did not ask for az or el, asked for %s", azel)
		}

		radAlpha = (2.0 * math.Pi * i) / 360.0
		cA = math.Cos(radAlpha)
		sA = math.Sin(radAlpha)
		hashX := innerX*cA + sD.centerX
		hashY := innerY*sA + sD.centerY
		endX := outerX*cA + sD.centerX
		endY := outerY*sA + sD.centerY

		ll, ul := calcLetter(int(i) + 90) //+90 is to correct for the fyne geometry
		txtX = txtRadX*cA + sD.centerX + ll
		txtY = txtRadY*sA + sD.centerY + ul

		switch azel {
		case "az":
			txt = canvas.NewText(fmt.Sprintf("%4.0f", i+90), color.Black)
		case "el":
			txt = canvas.NewText(fmt.Sprintf("%4.0f", 360.-i), color.Black)
		default:
			log.Printf("program bug, did not ask for az or el, asked for %s", azel)
		}
		txt.TextSize = smallText
		txt.Move(fyne.Position{float32(txtX), float32(txtY)})
		hashMarks = append(hashMarks, txt)
		ln := canvas.NewLine(red)
		ln.StrokeWidth = 2
		ln.Position1 = fyne.Position{float32(hashX), float32(hashY)}
		ln.Position2 = fyne.Position{float32(endX), float32(endY)}
		hashMarks = append(hashMarks, ln)
	}

	//Upper limit hashmark and label
	radAlpha = (2.0 * math.Pi) * (sD.upperLimit / 360.0)
	cA = math.Cos(radAlpha)
	sA = math.Sin(radAlpha)
	hashX = innerX*cA + sD.centerX
	hashY = innerX*sA + sD.centerY
	endX = outerX*cA + sD.centerX
	endY = outerY*sA + sD.centerY
	ln = canvas.NewLine(red)
	ln.StrokeWidth = 2
	ln.Position1 = fyne.Position{float32(hashX), float32(hashY)}
	ln.Position2 = fyne.Position{float32(endX), float32(endY)}
	hashMarks = append(hashMarks, ln)

	ll, ul = calcLetter(int(sD.upperLimit))
	txtX = limitRadX*cA + sD.centerX + ll
	txtY = limitRadY*sA + sD.centerY + ul
	switch azel {
	case "az":
		txt = canvas.NewText(fmt.Sprintf("%4.0f", sD.ul), color.Black)
	case "el":
		txt = canvas.NewText(fmt.Sprintf("%4.0f", sD.ul), color.Black)
	default:
		log.Printf("program bug, did not ask for az or el, asked for %s", azel)
	}
	txt.TextSize = smallText
	txt.Move(fyne.Position{float32(txtX), float32(txtY)})
	txt.TextStyle.Bold = true
	hashMarks = append(hashMarks, txt)

	return hashMarks

}

//func (app *application) makeAzDial() fyne.CanvasObject {

//}

//func (app *application) azimuthDial() *fyne.Container{
//
//    sDa := &scaleData{
//		ll:      app.minAz,
//		ul:      app.maxAz,
//		centerX: 250.0,
//		centerY: 250.0,
//		endX:    250.0,
//		endY:    100.0,
//	}
//	hashMarksa := sDa.makeScale("az")
//
//    //define the dial line and locate it
//    la := canvas.NewLine(red)
//    la.StrokeWidth = 3
//	lp1a := fyne.Position{float32(sDa.centerX), float32(sDa.centerY)}
//    la.Position1 = lp1a
//    radThetaE :=( (2.0 * math.Pi) * (app.currAz - 90.0)) / 360.0
//	endX := innerX * math.Cos(radThetaE) + sD.centerX
//	endY := innerY * math.Sin(radThetaE) + sD.centerY
//	lp2a := fyne.Position{float32(endX), float32(endY)}
//	la.Position2 = lp2a
//    app.azDialBind.Set(dialAzPos)
//	hashMarksa = append(hashMarksa, la)
//    return container.NewWithoutLayout(hashMarksa...)
//
//}
