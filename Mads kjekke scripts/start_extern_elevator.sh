#!/bin/bash
clear
echo "Username set to 'Student':"

echo "Type in the last byte of the IP to the elvator you want to connect to:"
read IP
echo "Connecting to 129.241.187."$IP
scp -rq /home/student/LC/SanntidsHeis/src/testmain student@129.241.187.$IP:~/Videos/testmain
ssh student@129.241.187.$IP "~/Videos/testmain"

echo Elevator script started


# scp -rq /home/student/LC/SanntidsHeis/src/testmain student@129.241.187.148:~/LC/testmain