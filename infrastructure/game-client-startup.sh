#!/bin/bash

# Script to install applications for the Game Client VM 

export DEBIAN_FRONTEND=noninteractive

cd /etc/apt/sources.list.d 
wget https://raw.githubusercontent.com/Xpra-org/xpra/master/packaging/repos/bookworm/xpra.sources
sudo apt update -y
cd -

curl -OL https://github.com/LizardByte/Sunshine/releases/download/v0.21.0/sunshine-debian-bookworm-amd64.deb
sudo apt install -f ./sunshine-debian-bookworm-amd64.deb -y
sudo apt install xfce4 xfce4-goodies tightvncserver dbus-x11 xpra -y

sudo apt upgrade -y
