package main

import (
	"fmt"
	"log"
)

var cmdLen map[cmdType]int = map[cmdType]int{
	moveFwd: 5,
	moveBwd: 5,
	moveUp:  5,
	moveDn:  5,
}

const address byte = 0x80
const ack byte = 0xFF

// roboClaw command to command byte mapping
var cmds map[cmdType]byte = map[cmdType]byte{
	moveFwd: 0,
	moveBwd: 1,
	moveUp:  4,
	moveDn:  5,
}


func (app *application) writeCmd(rc *roboClaw) error {
	wBuff := make([]byte, cmdLen[rc.cmd])
	rBuff := make([]byte, 1)

	wBuff[0] = address
	wBuff[1] = cmds[rc.cmd]
	wBuff[2] = rc.value[0]
	crc := crc16(wBuff[:3])
	wBuff[3] = byte(crc >> 8)
	wBuff[4] = byte(crc)

	n, err := app.port.Write(wBuff)
	if err != nil {
		return fmt.Errorf("failed write to usb port %v", err)
	}
	log.Println("Just wrote", n, wBuff)
	if n != cmdLen[rc.cmd] {
		return fmt.Errorf("did not write %d bytes, it wrote %d", cmdLen[rc.cmd], n)
	}
	n, err = app.port.Read(rBuff)
	if err != nil {
		return fmt.Errorf("failed to read from usb port: %v", err)
	}
	if n != 1 {
		return fmt.Errorf("did not read one byte, read %d:", err)
	}
	if rBuff[0] != ack {
		return fmt.Errorf("did not get a %X in return, got %v", ack, rBuff)
	}
	return nil
}

////midCmd has two arguments in form of a byte slice
//func (app *application) midCmd(buff []byte) error {
//    b[2] = byte('\r')
//
//    n, err := app.port.Write(wBuff)
//    if err != nil {
//        return fmt.Errorf("failed write to usb port %v", err)
//    }
//    if n != 3 {
//        return fmt.Errorf("did not write 2 bytes, it wrote %d", n)
//    }
//    n, err = app.port.Read(rBuff)
//	if err != nil {
//		return fmt.Errorf("failed to read from usb port: %v", err)
//    }
//    if n != 1 {
//        return fmt.Errorf("did not read one byte, read %d:", err)
//	}
//    if rBuff[0] != '\r' {
//        return fmt.Errorf("did not get a 0x)D in return, got %v", rBuff)
//    }
//    return nil
//
//
//}

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

// Example usage
//func main() {
//    data := []byte("Hello, world!")
//    checksum := crc16(data)
//    println(f"CRC16 checksum for '{data}': {checksum:04x}")
//}
