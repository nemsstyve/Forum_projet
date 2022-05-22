@REM ECHO off allows to print only the return of the command
@ECHO OFF
ECHO.
ECHO ---------------------BUILDING IMAGE DOCKER---------------------------
docker build -t forum .
ECHO. 

ECHO ---------------------RUNNING DOCKER ON 8080--------------------------
docker run -d --name Forum -p 8080:8080 forum
ECHO. 
ECHO ---------------------DOCKER IMAGE LIST-------------------------------
docker images
ECHO. 
ECHO ----------------------CONTAINER LIST--------------------------------
docker container ls
