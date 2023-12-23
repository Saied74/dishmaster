package main

import (
//	"fmt"
	"math"
//    "time"
)

//type controllerTime struct {
//	year  int
//	month int
//	day   int
//	hour  float64
//	min   float64
//	sec   float64
//	ut    float64
//}
//
//
//func main() {
//
//	t := controllerTime{}
//	t.getTime()
//
//	lat := 40.321490 //North
//	lon := -74.51024 //West
//
//	_, _, _, Az, El, _, _ := sun(t.year, t.month, t.day, t.ut, lon, lat)
//
//	fmt.Printf("Az: %5.2f\tEl: %5.2f\n", Az, El)
//
//}
//
//func (ct *controllerTime) getTime() {
//
//	localTime := time.Now()
//	t := localTime.UTC()
//	ct.year = t.Year()
//	//	m := t.Month()
//	//	mIndex := monthIndex(string(m))
//	//    if mIndex == -1 {
//	//        log.Fatal("Failed to index: ", m)
//	//    }
//	ct.month = 12
//	ct.day = t.Day()
//	ct.hour = float64(t.Hour())
//	ct.min = float64(t.Minute())
//	ct.sec = float64(t.Second())
//	ct.ut = ct.hour + (ct.min / 60.0) + (ct.sec / 3600.0)
//}


//subroutine sun(y,m,DD,UT,lon,lat,RA,Dec,LST,Az,El,mjd,day)

func sun(y, m, DD int, UT, lon, lat float64) (RA, Dec, LST, Az, El, mjd, day float64) {

	//  implicit none

	//  integer y                         !Year
	//  integer m                         !Month
	//  integer DD                        !Day
	//  integer mjd                       !Modified Julian Date
	//var mjd int
	//  real UT                           !UTC in hours
	//  real RA,Dec                       !RA and Dec of sun

	//! NB: Double caps here are single caps in the writeup.

	// ! Orbital elements of the Sun (also N=0, i=0, a=1):
	//
	//	real w                            !Argument of perihelion
	var w float64
	// real e                            !Eccentricity
	var e float64
	// real MM                           !Mean anomaly
	var MM float64
	// real Ls                           !Mean longitude
	var Ls float64

	// ! Other standard variables:
	//
	//	real v                            !True anomaly
	var v float64
	// real EE                           !Eccentric anomaly
	var EE float64
	// real ecl                          !Obliquity of the ecliptic
	var ecl float64
	// real d                            !Ephemeris time argument in days
	var d float64
	// real r                            !Distance to sun, AU
	var r float64
	// real xv,yv                        !x and y coords in ecliptic
	var xv, yv float64
	// real lonsun                       !Ecliptic long and lat of sun
	var lonsun float64
	// ! Ecliptic coords of sun (geocentric)
	//
	//	real xs,ys
	var xs, ys float64
	// ! Equatorial coords of sun (geocentric)
	//
	//	real xe,ye,ze
	var xe, ye, ze float64
	//  real lon,lat

	// real GMST0,LST,HA
	var GMST0, HA float64
	// real xx,yy,zz
	var xx, yy, zz float64
	// real xhor,yhor,zhor
	var xhor, yhor, zhor float64
	//  real Az,El

	// real day
	// real rad
	// data rad/57.2957795/
	var rad float64 = 57.2957795

	//! Time in days, with Jan 0, 2000 equal to 0.0:
	dd := 367*y - 7*(y+(m+9)/12)/4 + 275*m/9 + DD - 730530
	d = float64(dd) + UT/24.0
	mjd = d + 51543.0
	ecl = 23.4393 - 3.563e-7*d

	//! Compute updated orbital elements for Sun:
	w = 282.9404 + 4.70935e-5*d
	e = 0.016709 - 1.151e-9*d
	MM = math.Mod((float64(356.0470) + float64(0.9856002585)*d + float64(360000.0)), (float64(360.0)))
	Ls = math.Mod((w + MM + 720.0), (360.0))

	EE = MM + e*rad*math.Sin(MM/rad)*(1.0+e*math.Cos(MM/rad))
	EE = EE - (EE-e*rad*math.Sin(EE/rad)-MM)/(1.0-e*math.Cos(EE/rad))

	xv = math.Cos(EE/rad) - e
	yv = math.Sqrt(1.0-e*e) * math.Sin(EE/rad)
	v = rad * math.Atan2(yv, xv)
	r = math.Sqrt(xv*xv + yv*yv)
	lonsun = math.Mod((v + w + 720.0), (360.0))
	//! Ecliptic coordinates of sun (rectangular):
	xs = r * math.Cos(lonsun/rad)
	ys = r * math.Sin(lonsun/rad)

	//! Equatorial coordinates of sun (rectangular):
	xe = xs
	ye = ys * math.Cos(ecl/rad)
	ze = ys * math.Sin(ecl/rad)

	//! RA and Dec in degrees:
	RA = rad * math.Atan2(ye, xe)
	Dec = rad * math.Atan2(ze, math.Sqrt(xe*xe+ye*ye))

	GMST0 = (Ls + 180.0) / 15.0
	LST = math.Mod(GMST0+UT+lon/15.0+48.0, 24.0) //!LST in hours
	HA = 15.0*LST - RA                      //!HA in degrees
	xx = math.Cos(HA/rad) * math.Cos(Dec/rad)
	yy = math.Sin(HA/rad) * math.Cos(Dec/rad)
	zz = math.Sin(Dec / rad)
	xhor = xx*math.Sin(lat/rad) - zz*math.Cos(lat/rad)
	yhor = yy
	zhor = xx*math.Cos(lat/rad) + zz*math.Sin(lat/rad)
	Az = math.Mod((rad*math.Atan2(yhor, xhor) + 180.0 + 360.0), (360.0))
	El = rad * math.Asin(zhor)
	day = d - 1.5

	return RA, Dec, LST, Az, El, mjd, day
	// end subroutine sun
}
