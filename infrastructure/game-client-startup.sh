#!/bin/bash

# Script to install applications for the Game Client VM 

export DEBIAN_FRONTEND=noninteractive

curl -1sLf 'https://dl.cloudsmith.io/public/moonlight-game-streaming/moonlight-l4t/setup.deb.sh' | sudo -E bash
sudo apt install moonlight-qt
sudo apt install xfce4 xfce4-goodies -y
sudo apt install tightvncserver -y
sudo apt install dbus-x11 -y

sudo apt update -y
sudo apt upgrade -y
