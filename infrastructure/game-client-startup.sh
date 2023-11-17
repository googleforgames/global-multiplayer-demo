#!/bin/bash

# Script to install applications for the Game Client VM 

export DEBIAN_FRONTEND=noninteractive

curl -OL https://github.com/LizardByte/Sunshine/releases/download/v0.21.0/sunshine-debian-bookworm-amd64.deb
sudo apt install -f ./sunshine-debian-bookworm-amd64.deb -y
sudo apt install xfce4 xfce4-goodies -y
sudo apt install tightvncserver -y
sudo apt install dbus-x11 -y

sudo apt update -y
sudo apt upgrade -y
