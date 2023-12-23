package main

//
//import (
//
//    "fmt"
//
//)
//
////simpleCmd has just one argument
//func (app *application) simpleCmd(c byte) error {
//    wBuff := make([]byte, 5)
//    rBuff := make([]byte, 5)
//
//    wBuff[0] = byte(c)
//    wBuff[1] = byte('\r')
//
//    n, err := app.port.Write(wBuff)
//    if err != nil {
//        return fmt.Errorf("failed write to usb port %v", err)
//    }
//    if n != 2 {
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
//}
//
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
//
////func crc16(message []byte) uint16 {
////    crc := uint16(0xFFFF) // Initial value
////    for _, b := range message {
////        crc ^= uint16(b) << 8
////        for i := 0; i < 8; i++ {
////            if crc&0x8000 != 0 {
////                crc = (crc << 1) ^ 0x1021
////            } else {
////                crc <<= 1
////            }
////        }
////    }
////    return crc
////}
//
//// Example usage
////func main() {
////    data := []byte("Hello, world!")
////    checksum := crc16(data)
////    println(f"CRC16 checksum for '{data}': {checksum:04x}")
////}
