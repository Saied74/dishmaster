package main

import (
	"fmt"
	"math"
	"strings"
	"unicode"
	"unicode/utf8"
)

//func main() {
//
//	x := "FN20rh"
//	lat, lon, err := grid2Deg(x)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("%6.3f\t%6.3f\n", lat, lon)
//	lat, lon, err = otherGrid2Deg(x)
//	if err != nil {
//		log.Fatal(err)

// Translated from the WSJT Fortran code by Pete VE5VA
// Traslated to Go by Saied Seghatoleslami AD2CC

/////////////////////////////////////////////////////////
/////////////        G R I D 2 D E G        /////////////
/////////////////////////////////////////////////////////
// Choose which version of the grid conversion to use.
// Defining WSJT uses the one from WSJT
// Removing the define of WSJT uses the code based on
// the PERL script at wikipedia which seems to be
// slightly more accurate.
//#define WSJT

func grid2Deg(grid string) (dlat, dlong float64, err error) {

	newGrid := strings.ToUpper(grid)
	var nlong, n20d int
	var nlat int
	var xminlong, xminlat float64

	gp, err := validateGrid(newGrid)
	if err != nil {
		return 0.0, 0.0, err
	}
	nlong = 180 - 20*gp.pos0
	n20d = 2 * gp.pos2
	xminlong = 5*float64(gp.pos4) + 0.5

	dlong = -(float64(nlong-n20d) - xminlong/60.0)
	nlat = -90 + 10*gp.pos1 + gp.pos3
	xminlat = 2.5 * (float64(gp.pos5) + 0.5)
	dlat = float64(nlat) + xminlat/60.0
	return dlat, dlong, nil
}

// From: http://en.wikipedia.org/wiki/Maidenhead_Locator_System
func otherGrid2Deg(grid string) (dlat, dlong float64, err error) {

	newGrid := strings.ToUpper(grid)

	gp, err := validateGrid(newGrid)
	if err != nil {
		return 0.0, 0.0, err
	}

	dlong = 20.0*float64(gp.pos0) - 180.0
	dlat = 10.0*float64(gp.pos1) - 90.0
	dlong += float64(gp.pos2) * 2.0
	dlat += float64(gp.pos3)

	dlong += float64(gp.pos4) * (5.0 / 60.0)
	dlat += float64(gp.pos5) * (2.5 / 60.0)
	dlong += 2.5 / 60.0
	dlat += 1.25 / 60.0
	return dlat, dlong, nil
}

/////////////////////////////////////////////////////////
//////////////        D C O O R D         ///////////////
/////////////////////////////////////////////////////////
// In WSJT this is used in various places to do coordinate
// system conversions but moon2 only uses it once.
//void DCOORD(double xA0,double xB0,double AP,double BP,
//              double xA1,double xB1,double *xA2,double *B2)
//{

func DCOORD(xA0, xB0, AP, BP, xA1, xB1 float64) (xA2, B2 float64) {

	var TA2O2 float64
	var SB0, CB0, SBP, CBP, SB1, CB1, SB2, CB2 float64
	var SAA, CAA, SBB, CBB, CA2, SA2 float64

	SB0 = math.Sin(xB0)
	CB0 = math.Cos(xB0)
	SBP = math.Sin(BP)
	CBP = math.Cos(BP)
	SB1 = math.Sin(xB1)
	CB1 = math.Cos(xB1)
	SB2 = SBP*SB1 + CBP*CB1*math.Cos(AP-xA1)
	CB2 = math.Sqrt(1.0 - (SB2 * SB2))
	B2 = math.Atan(SB2 / CB2)
	SAA = math.Sin(AP-xA1) * CB1 / CB2
	CAA = (SB1 - SB2*SBP) / (CB2 * CBP)
	CBB = SB0 / CBP
	SBB = math.Sin(AP-xA0) * CB0
	SA2 = SAA*CBB - CAA*SBB
	CA2 = CAA*CBB + SAA*SBB
	TA2O2 = 0.0
	if CA2 <= 0.0 {
		TA2O2 = (1. - CA2) / SA2
	}
	if CA2 > 0.0 {
		TA2O2 = SA2 / (1. + CA2)
	}
	xA2 = 2.0 * math.Atan(TA2O2)
	if xA2 < 0.0 {
		xA2 = xA2 + 6.2831853071795864
	}
	return xA2, B2
}

/////////////////////////////////////////////////////////
////////////////        M O O N 2        ////////////////
/////////////////////////////////////////////////////////

// You can derive the lat/long from the grid square umath.Sing
// the grid2deg function to translate from grid to lat/long
// Example call to this function for 2014/01/04 1709Z:
//   moon2(2014,1,4,17+9/60.,-106.625,52.104168,
//       &RA, &Dec, &topRA, &topDec, &LST, &HA, &Az, &El, &dist);
// I have not used or checked any of the outputs other than Az and El
//void moon2(int y,int m,int Day,
//        double UT,
//        double lon,double lat,
//        double *RA,double *Dec,
//        double *topRA,double *topDec,
//        double *LST,double *HA,
//        double *Az,double *El,double *dist)
//{

func moon2(y, m, Day int, UT, lon, lat float64) (RA, Dec, topRA, topDec, LST, HA, Az, El, dist float64) {

	// The strange position of some of the semicolons is because some
	// of this was translated umath.Sing a TCL script that I wrote - it isn't
	// too smart but it saved some typing
	var NN float64 //Longitude of ascending node
	var i float64  //Inclination to the ecliptic
	var w float64  //Argument of perigee
	var a float64  //Semi-major axis
	var e float64  //Eccentricity
	var MM float64 //Mean anomaly

	var v float64   //True anomaly
	var EE float64  //Eccentric anomaly
	var ecl float64 //Obliquity of the ecliptic

	var d float64              //Ephemeris time argument in days
	var r float64              //Distance to sun, AU
	var xv, yv float64         //x and y coords in ecliptic
	var lonecl, latecl float64 //Ecliptic long and lat of moon
	var xg, yg, zg float64     //Ecliptic rectangular coords
	var Ms float64             //Mean anomaly of sun
	var ws float64             //Argument of perihelion of sun
	var Ls float64             //Mean longitude of sun (Ns=0)
	var Lm float64             //Mean longitude of moon
	var DD float64             //Mean elongation of moon
	var FF float64             //Argument of latitude for moon
	var xe, ye, ze float64     //Equatorial geocentric coords of moon
	var mpar float64           //Parallax of moon (r_E / d)
	//  double lat,lon            float64//Station coordinates on earth
	var gclat float64 //Geocentric latitude
	var rho float64   //Earth radius factor
	var GMST0 float64 //,LST,HA;
	var g float64

	var rad float64 = 57.2957795131
	var twopi float64 = 6.283185307
	var pi, pio2 float64

	dint := int32(367)*int32(y) - int32(7)*(int32(y)+(int32(m)+int32(9))/int32(12))/int32(4) + int32(275)*int32(m)/int32(9) + int32(Day)
	d = float64(dint) - 730530.0 + UT/24.0

	ecl = 23.4393 - 3.563e-7*d

	NN = 125.1228 - 0.0529538083*d
	i = 5.1454
	w = math.Mod(318.0634+0.1643573223*d+360000., 360.0)
	a = 60.2666
	e = 0.054900
	MM = math.Mod(115.3654+13.0649929509*d+360000., 360.0)

	EE = MM + e*rad*math.Sin(MM/rad)*(1.+e*math.Cos(MM/rad))
	EE = EE - (EE-e*rad*math.Sin(EE/rad)-MM)/(1.-e*math.Cos(EE/rad))
	EE = EE - (EE-e*rad*math.Sin(EE/rad)-MM)/(1.-e*math.Cos(EE/rad))

	xv = a * (math.Cos(EE/rad) - e)
	yv = a * (math.Sqrt(1.-e*e) * math.Sin(EE/rad))

	v = math.Mod(rad*math.Atan2(yv, xv)+720.0, 360.0)
	r = math.Sqrt(xv*xv + yv*yv)

	//  Get geocentric position in ecliptic recmath.Tangular coordinates:

	xg = r * (math.Cos(NN/rad)*math.Cos((v+w)/rad) - math.Sin(NN/rad)*math.Sin((v+w)/rad)*math.Cos(i/rad))
	yg = r * (math.Sin(NN/rad)*math.Cos((v+w)/rad) + math.Cos(NN/rad)*math.Sin((v+w)/rad)*math.Cos(i/rad))
	zg = r * (math.Sin((v+w)/rad) * math.Sin(i/rad))

	//  Ecliptic longitude and latitude of moon:
	lonecl = math.Mod(rad*math.Atan2(yg/rad, xg/rad)+720.0, 360.0)
	latecl = rad * math.Atan2(zg/rad, math.Sqrt(xg*xg+yg*yg)/rad)

	//  Now include orbital perturbations:
	Ms = math.Mod(356.0470+0.9856002585*d+3600000.0, 360.0)
	ws = 282.9404 + 4.70935e-5*d
	Ls = math.Mod(Ms+ws+720.0, 360.0)
	Lm = math.Mod(MM+w+NN+720., 360.)
	DD = math.Mod(Lm-Ls+360.0, 360.0)
	FF = math.Mod(Lm-NN+360.0, 360.0)

	lonecl = lonecl -
		1.274*math.Sin((MM-2.0*DD)/rad) +
		0.658*math.Sin(2.0*DD/rad) -
		0.186*math.Sin(Ms/rad) -
		0.059*math.Sin((2.0*MM-2.0*DD)/rad) -
		0.057*math.Sin((MM-2.0*DD+Ms)/rad) +
		0.053*math.Sin((MM+2.0*DD)/rad) +
		0.046*math.Sin((2.0*DD-Ms)/rad) +
		0.041*math.Sin((MM-Ms)/rad) -
		0.035*math.Sin(DD/rad) -
		0.031*math.Sin((MM+Ms)/rad) -
		0.015*math.Sin((2.0*FF-2.0*DD)/rad) +
		0.011*math.Sin((MM-4.0*DD)/rad)

	latecl = latecl -
		0.173*math.Sin((FF-2.0*DD)/rad) -
		0.055*math.Sin((MM-FF-2.0*DD)/rad) -
		0.046*math.Sin((MM+FF-2.0*DD)/rad) +
		0.033*math.Sin((FF+2.0*DD)/rad) +
		0.017*math.Sin((2.0*MM+FF)/rad)

	r = 60.36298 -
		3.27746*math.Cos(MM/rad) -
		0.57994*math.Cos((MM-2.0*DD)/rad) -
		0.46357*math.Cos(2.0*DD/rad) -
		0.08904*math.Cos(2.0*MM/rad) +
		0.03865*math.Cos((2.0*MM-2.0*DD)/rad) -
		0.03237*math.Cos((2.0*DD-Ms)/rad) -
		0.02688*math.Cos((MM+2.0*DD)/rad) -
		0.02358*math.Cos((MM-2.0*DD+Ms)/rad) -
		0.02030*math.Cos((MM-Ms)/rad) +
		0.01719*math.Cos(DD/rad) +
		0.01671*math.Cos((MM+Ms)/rad)

	dist = r * 6378.140

	//  Geocentric coordinates:
	//  Recmath.Tangular ecliptic coordinates of the moon:

	xg = r * math.Cos(lonecl/rad) * math.Cos(latecl/rad)
	yg = r * math.Sin(lonecl/rad) * math.Cos(latecl/rad)
	zg = r * math.Sin(latecl/rad)

	//  Recmath.Tangular equatorial coordinates of the moon:
	xe = xg
	ye = yg*math.Cos(ecl/rad) - zg*math.Sin(ecl/rad)
	ze = yg*math.Sin(ecl/rad) + zg*math.Cos(ecl/rad)

	//  Right Ascension, Declination:
	RA = math.Mod(rad*math.Atan2(ye, xe)+360.0, 360.0)
	Dec = rad * math.Atan2(ze, math.Sqrt(xe*xe+ye*ye))

	//  Now convert to topocentric system:
	mpar = rad * math.Sin(1.0/r)
	//      alt_topoc = alt_geoc - mpar*math.Cos(alt_geoc)
	gclat = lat - 0.1924*math.Sin(2.0*lat/rad)
	rho = 0.99883 + 0.00167*math.Cos(2.0*lat/rad)
	GMST0 = (Ls + 180.0) / 15.0
	LST = math.Mod(GMST0+UT+lon/15.0+48.0, 24.0) //LST in hours

	HA = 15.0*LST - RA //HA in degrees
	g = rad * math.Atan(math.Tan(gclat/rad)/math.Cos(HA/rad))
	topRA = RA - mpar*rho*math.Cos(gclat/rad)*math.Sin(HA/rad)/math.Cos(Dec/rad)
	topDec = Dec - mpar*rho*math.Sin(gclat/rad)*math.Sin((g-Dec)/rad)/math.Sin(g/rad)

	HA = 15.0*LST - topRA //HA in degrees
	if HA > 180.0 {
		HA = HA - 360.
	}
	if HA < -180.0 {
		HA = HA + 360.
	}
	pi = 0.5 * twopi
	pio2 = 0.5 * pi
	Az, El = DCOORD(pi, pio2-lat/rad, 0.0, lat/rad, HA*twopi/360.0, topDec/rad) //, Az, El)
	Az = Az * rad
	El = El * rad

	return RA, Dec, topRA, topDec, LST, HA, Az, El, dist
}

type gridPositions struct {
	pos0 int
	pos1 int
	pos2 int
	pos3 int
	pos4 int
	pos5 int
}

var alphaError = fmt.Errorf("letter is not in the English upper case alphabet")
var numberError = fmt.Errorf("letter is not a number in the range of 0-9")

func validateGrid(grid string) (gP *gridPositions, err error) {
	abc := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	oneTwo := "0123456789"
	gP = &gridPositions{}
	var j int
	n := 0
	start := 0
	upGrid := strings.ToUpper(grid)
	for i := 0; i < 6; i++ {
		for {
			r, w := utf8.DecodeRuneInString(upGrid[start:])
			c := string(r)
			switch i {
			case 0:
				if unicode.IsLetter(r) {
					n++
					start += w
					j = strings.Index(abc, c)
					if j == -1 {
						return gP, alphaError
					}
					gP.pos0 = j
				}
			case 1:
				if unicode.IsLetter(r) {
					n++
					start += w
					j = strings.Index(abc, c)
					if j == -1 {
						return gP, alphaError
					}
					gP.pos1 = j
				}
			case 2:
				if unicode.IsNumber(r) {
					n++
					start += w
					j = strings.Index(oneTwo, c)
					if j == -1 {
						return gP, numberError
					}
					gP.pos2 = j
				}
			case 3:
				if unicode.IsNumber(r) {
					n++
					start += w
					j = strings.Index(oneTwo, c)
					if j == -1 {
						return gP, numberError
					}
					gP.pos3 = j
				}
			case 4:
				if unicode.IsLetter(r) {
					n++
					start += w
					j = strings.Index(abc, c)
					if j == -1 {
						return gP, alphaError
					}
					gP.pos4 = j
				}
			case 5:
				if unicode.IsLetter(r) {
					n++
					start += w
					j = strings.Index(abc, c)
					if j == -1 {
						return gP, alphaError
					}
					gP.pos5 = j
				}
			}
			break
		}
	}
	if n != 6 {
		return gP, fmt.Errorf("scanning of the grid string did not parse correctly %d", n)
	}
	return gP, nil
}

func monthIndex(m string) int {

	var months = []string{"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "Novenmber", "December"}

	for i, month := range months {
		if m == month {
			return i
		}
	}
	return -1

}
