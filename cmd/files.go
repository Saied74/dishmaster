package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type MasterData struct {
	Grid   string  `json:"grid"`
	Lat    float64 `json:"lat"`
	Lon    float64 `json:"lon"`
	ParkAz float64 `json:"parkAz"`
	ParkEl float64 `json:"parkEl"`
	MaxAz  float64 `json:"maxAz"`
	MinAz  float64 `json:"minAz"`
	MaxEl  float64 `json:"maxEl"`
	MinEl  float64 `json:"minEL"`
}

type dishData struct {
	CurrAz     float64 `json:"currAz"`
	CurrEl     float64 `json:"currEl"`
	AzPosition float64 `json:"azPosition"`
	ElPosition float64 `json:"elPosition"`
}

func (app *application) getMasterData() error {
	m := make([]byte, 300)
	j := &MasterData{}

	f, err := os.OpenFile(app.masterPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("Err openning master file: %v", err)
	}
	_, err = f.Read(m)
	if err != nil {
		return fmt.Errorf("Err reading master file %v:", err)
	}

	before, _, _ := bytes.Cut(m, []byte{0x00})
	err = json.Unmarshal(before, j)
	if err != nil {
		return fmt.Errorf("Err unmarshaling master file: %v", err)
	}
	app.grid = j.Grid
	app.lat = j.Lat
	app.lon = j.Lon
	app.parkAz = j.ParkAz
	app.parkEl = j.ParkEl
	app.maxAz = j.MaxAz
	app.minAz = j.MinAz
	app.maxEl = j.MaxEl
	app.minEl = j.MinEl
	return nil
}

func (app *application) getDishData() error {
	m := make([]byte, 200)
	j := &dishData{}

	f, err := os.Open(app.dishPath)
	if err != nil {
		return fmt.Errorf("Err openning dish file: %v", err)
	}
	_, err = f.Read(m)
	if err != nil {
		return fmt.Errorf("Err reading dish file: %v", err)
	}
	before, _, _ := bytes.Cut(m, []byte{0x00})
	err = json.Unmarshal(before, j)
	if err != nil {
		return fmt.Errorf("Err unmarshaling dish file: %v", err)
	}
	app.currAz = j.CurrAz
	app.currEl = j.CurrEl
	app.azPosition = j.CurrAz //j.AzPosition //for debugging ease
	app.elPosition = j.CurrEl //j.ElPosition
	return nil
}

func (app *application) saveMasterData() error {

	j := &MasterData{
		Grid:   app.grid,
		Lat:    app.lat,
		Lon:    app.lon,
		ParkAz: app.parkAz,
		ParkEl: app.parkEl,
		MaxAz:  app.maxAz,
		MinAz:  app.minAz,
		MaxEl:  app.maxEl,
		MinEl:  app.minEl,
	}
	m, err := json.Marshal(j)
	if err != nil {
		return fmt.Errorf("Err marshaling master data: %v", err)
	}
	m = append(m, 0x00)

	f, err := os.OpenFile(app.masterPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("Err openning master file: %v", err)
	}
	defer f.Close()
	_, err = f.Write(m)
	if err != nil {
		return fmt.Errorf("Err wirting master file: %v", err)
	}
	fmt.Println(j)
	return nil
}

func (app *application) saveDishData() error {
	m := []byte{}
	j := &dishData{
		CurrAz:     app.currAz,
		CurrEl:     app.currEl,
		AzPosition: app.azPosition,
		ElPosition: app.elPosition,
	}

	m, err := json.Marshal(j)
	if err != nil {
		return fmt.Errorf("Err marshaling dish data: %v", err)
	}
	m = append(m, 0x00)

	f, err := os.OpenFile(app.dishPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return fmt.Errorf("Err openning dish file: %v", err)
	}
	defer f.Close()
	_, err = f.Write(m)
	if err != nil {
		return fmt.Errorf("Err writing dish file: %v", err)
	}
	fmt.Println(j)
	return nil
}
