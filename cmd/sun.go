package main

import (
	"math"
)

//type controllerTime struct {
//	year  int
//	month int
//	day   int
//	hour  float64
//	min   float64

//Original from WSJT-X by Joe Taylor K1JT sun.f90
//Translated to Go by Saied Seghatoleslami AD2CC
//subroutine sun(y,m,DD,UT,lon,lat,RA,Dec,LST,Az,El,mjd,day)

func sun(y, m, DD int, UT, lon, lat float64) (RA, Dec, LST, Az, El, mjd, day float64) {

	//! NB: Double caps here are single caps in the writeup.

	// ! Orbital elements of the Sun (also N=0, i=0, a=1):
	//
	var w float64
	var e float64
	var MM float64
	var Ls float64

	// ! Other standard variables:
	//
	var v float64
	var EE float64
	var ecl float64
	var d float64
	var r float64
	var xv, yv float64
	var lonsun float64
	// ! Ecliptic coords of sun (geocentric)
	//
	var xs, ys float64
	// ! Equatorial coords of sun (geocentric)
	//
	var xe, ye, ze float64

	// real GMST0,LST,HA
	var GMST0, HA float64
	var xx, yy, zz float64
	var xhor, yhor, zhor float64

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
	HA = 15.0*LST - RA                           //!HA in degrees
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
}
