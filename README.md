# dishmaster
Program for controlling a dish antenna


Dishmaster is designed to work with a RoboClaw motor controller.  Before using it to
track celestial objects, it needs to be calibrated.  Here are the steps for calibration:

1. Find the approximate current azymuth and elevation of the antenna
2. Start the Dishmaster program
3. On the setup tab, enter the required information and update the applicattion
4. Enter the approximate azimuth and elevation into the targett azimuth and elevation windows on the operate tab
5. Push the recalibratte button - now the dish is approximately calibrated
6. Find the location of the moon or the sun (assuming they are visible) and move to it
6. Use the adjust buttons (they adjustt by half degree increments) to zero in on the target body
7. Make sure that you wait until the current position has cought up with the target position
8. Enter these values into the target windows again and push recalibrate
9. Now the Dishmaster is fully calibrated

System data is kept in file master.json that is in the project root directory.  Below is an example of this file.

{"grid":"FN20nh","lat":40.3125,"lon":-74.90833333333333,"parkAz":90,"parkEl":20,"maxAz":315,"minAz":45,"maxEl":90,"minEL":0}

Lattitude and longitude are calculated from grid.  Park, min, and max azimuth and elevation are entered via the setup window.



