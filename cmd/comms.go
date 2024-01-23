package main

import (
    "errors"
	"fmt"
	// "log"
//    "time"

)

var cmdLen map[int]int = map[int]int{
	moveFwd: 5,
	moveBwd: 5,
	moveUp:  5,
	moveDn:  5,
}

const (
	address   byte = 0x80
	ack       byte = 0xFF
	underFlow      = 0x01
	overFlow       = 0x04
)

// roboClaw command to command byte mapping
var cmds map[int]byte = map[int]byte{
	moveFwd: 0,
	moveBwd: 1,
	moveUp:  4,
	moveDn:  5,
}

var noReadN = errors.New("readN did not read anything")

// for commands that return a single byte ack
func (app *application) writeCmd(rc *roboClaw) error {
	l := cmdLen[rc.cmd]
	wBuff := make([]byte, l)
	rBuff := make([]byte, 1)

	wBuff[0] = address
	wBuff[1] = cmds[rc.cmd]
	wBuff[2] = rc.value
	crc := crc16(wBuff[:l-2])
	wBuff[l-2] = byte(crc >> 8)
	wBuff[l-1] = byte(crc)
	if app.port == nil {
		return fmt.Errorf("USB port is nil")
	}
	n, err := app.port.Write(wBuff)
	if err != nil {
		return fmt.Errorf("failed write to usb port %v", err)
	}
	if n != l {
		return fmt.Errorf("did not write %d bytes, it wrote %d", l, n)
	}
	n, err = app.port.Read(rBuff)
	if err != nil {
		return fmt.Errorf("failed to read from usb port: %v", err)
	}
	if n != 1 {
		return fmt.Errorf("did not read one byte, read %d:", n)
	}
	if rBuff[0] != ack {
		return fmt.Errorf("did not get a %X in return, got %v", ack, rBuff)
	}
	return nil
}

func (app *application) writeQuadRegister(c int32, s string) error {
	wBuff := make([]byte, 6)
	rBuff := make([]byte, 1)
	wBuff[0] = address
	switch s {
	case "az":
		wBuff[1] = 22
	case "el":
		wBuff[1] = 23
	default:
		return fmt.Errorf("Bad command \"%s\" in writeRegister", s)
	}
	wBuff[2] = byte(c >> 24)
	wBuff[3] = byte(c >> 16)
	wBuff[4] = byte(c >> 8)
	wBuff[5] = byte(c)

	crc := crc16(wBuff)
	wBuff = append(wBuff, byte(crc >> 8))
	wBuff = append(wBuff, byte(crc))
	if app.port == nil {
		return fmt.Errorf("Port is not open %v", wBuff)
	}
	n, err := app.port.Write(wBuff)
	if err != nil {
		return fmt.Errorf("failed write to usb port %v", err)
	}
	if n != 8 {
		return fmt.Errorf("did not write 8 bytes, it wrote %d", n)
	}
	n, err = app.port.Read(rBuff)
	if err != nil {
		return fmt.Errorf("failed to read from usb port: %v", err)
	}
	if n != 1 {
		return fmt.Errorf("did not read one byte, read %d:", n)
	}
	if rBuff[0] != ack {
		return fmt.Errorf("did not get a %X in return, got %v", ack, rBuff)
	}
	return nil
}

func (app *application) readQuadRegister(s string) (int32, error) {
	var r int32
	wBuff := make([]byte, 2)
	//rBuff := make([]byte, 7)
	wBuff[0] = address
	switch s {
	case "az":
		wBuff[1] = 16
	case "el":
		wBuff[1] = 17
	default:
		return 0, fmt.Errorf("Bad command \"%s\" in writeQuadRegister", s)
	}
	if app.port == nil {
		return 0, fmt.Errorf("Port is not open %v", wBuff)
	}
	n, err := app.port.Write(wBuff)
	if err != nil {
		return 0, fmt.Errorf("failed write to usb port %v", err)
	}
	if n != 2 {
		return 0, fmt.Errorf("did not write 2 bytes, it wrote %d", n)
	}
    //time.Sleep(time.Duration(120) * time.Millisecond)
	//n, err = app.port.Read(rBuff)
    rBuff, err := app.readN(7)
	if err != nil {
        if errors.Is(err, noReadN) {
            return 0, err
        }
		return 0, fmt.Errorf("failed to read from usb port: %v", err)
	}
	//if n != 7 {
	//	return 0, fmt.Errorf("did not read 9 byte, read %d:", n)
	//}
	crc := crc16(rBuff[0:5])
	highByte := byte(crc >> 8)
	lowByte := byte(crc)
	if highByte != rBuff[5] && lowByte != rBuff[6] {
		return 0, fmt.Errorf("CRC mismtach on read quad registers %v", rBuff)
	}
	if rBuff[4]&underFlow != 0x00 {
		return 0, fmt.Errorf("Quad counter underflowed")
	}
	if rBuff[4]&overFlow != 0x00 {
		return 0, fmt.Errorf("Quad counter overflowed")
	}

	r = r | int32(rBuff[0])<<24
	r = r | int32(rBuff[1])<<16
	r = r | int32(rBuff[2])<<8
	r = r | int32(rBuff[3])

	return r, nil
}

// call this function with the exact slice you want processed
func crc16(message []byte) uint16 {
	crc := uint16(0xFFFF) // Initial value
	for _, b := range message {
		crc ^= uint16(b) << 8
		for i := 0; i < 8; i++ {
			if crc&0x8000 != 0 {
				crc = (crc << 1) ^ 0x1021
			} else {
				crc <<= 1
			}
		}
	}
	return crc
}

func (app *application) readN(n int) ([]byte, error) {
    rBuff := []byte{}
    b := make([]byte, n)
    k := 0
    m := 0 //read byte counter
    for i := 0; i < n; i += k {
        k, err := app.port.Read(b)
        //fmt.Println("K: ", k, n)
	    if err != nil {
		    return []byte{}, fmt.Errorf("failed readN")
	    }
        if k == 0 {
            return []byte{}, noReadN
        }
	    if k == n {
            return b, nil
        }
        m += k
        rBuff = append(rBuff, b[:k]...)
        //fmt.Println("rBuff: ", rBuff)
	}
    if m != n {
        return []byte{}, fmt.Errorf("did not read %d bytes, read %d", n, m)
    } 
    return rBuff, nil

} 
