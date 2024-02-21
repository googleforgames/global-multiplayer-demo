#!/bin/bash

# Script to install applications for the Game Client VM 

export DEBIAN_FRONTEND=noninteractive

wget -O "/usr/share/keyrings/xpra.asc" https://xpra.org/xpra.asc
cd /etc/apt/sources.list.d 
wget https://raw.githubusercontent.com/Xpra-org/xpra/master/packaging/repos/bookworm/xpra.sources
sudo apt update -y
cd -

curl -OL https://github.com/LizardByte/Sunshine/releases/download/v0.21.0/sunshine-debian-bookworm-amd64.deb
sudo apt install -f ./sunshine-debian-bookworm-amd64.deb -y
sudo apt install xfce4 xfce4-goodies tightvncserver dbus-x11 gcc make linux-headers-$(uname -r) software-properties-common mesa-utils golang libgl1-mesa-dev xorg-dev xpra -y 

wget https://developer.download.nvidia.com/compute/cuda/12.3.1/local_installers/cuda_12.3.1_545.23.08_linux.run
sudo sh cuda_12.3.1_545.23.08_linux.run --silent

sudo apt update -y

#########
#
# Startup script to install Chrome remote desktop and a desktop environment.
#
# See environmental variables at then end of the script for configuration
#

function install_desktop_env {
  PACKAGES="desktop-base xscreensaver dbus-x11"

  if [[ "$INSTALL_XFCE" != "yes" && "$INSTALL_CINNAMON" != "yes" ]] ; then
    # neither XFCE nor cinnamon specified; install both
    INSTALL_XFCE=yes
    INSTALL_CINNAMON=yes
  fi

  if [[ "$INSTALL_XFCE" = "yes" ]] ; then
    PACKAGES="$PACKAGES xfce4"
    echo "exec xfce4-session" > /etc/chrome-remote-desktop-session
    [[ "$INSTALL_FULL_DESKTOP" = "yes" ]] && \
      PACKAGES="$PACKAGES task-xfce-desktop"
  fi

  if [[ "$INSTALL_CINNAMON" = "yes" ]] ; then
    PACKAGES="$PACKAGES cinnamon-core"
    echo "exec cinnamon-session-cinnamon2d" > /etc/chrome-remote-desktop-session
    [[ "$INSTALL_FULL_DESKTOP" = "yes" ]] && \
      PACKAGES="$PACKAGES task-cinnamon-desktop"
  fi

  DEBIAN_FRONTEND=noninteractive \
    apt-get install --assume-yes $PACKAGES $EXTRA_PACKAGES

  systemctl disable lightdm.service
}

function download_and_install { # args URL FILENAME
  curl -L -o "$2" "$1"
  apt-get install --assume-yes --fix-broken "$2"
}

function is_installed {  # args PACKAGE_NAME
  dpkg-query --list "$1" | grep -q "^ii" 2>/dev/null
  return $?
}

# Configure the following environmental variables as required:
INSTALL_XFCE=yes
INSTALL_CINNAMON=yes
INSTALL_CHROME=yes
INSTALL_FULL_DESKTOP=yes

# Any additional packages that should be installed on startup can be added here
EXTRA_PACKAGES="less bzip2 zip unzip tasksel wget"

apt-get update

! is_installed chrome-remote-desktop && \
  download_and_install \
    https://dl.google.com/linux/direct/chrome-remote-desktop_current_amd64.deb \
    /tmp/chrome-remote-desktop_current_amd64.deb

install_desktop_env

[[ "$INSTALL_CHROME" = "yes" ]] && \
  ! is_installed google-chrome-stable && \
  download_and_install \
    https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb \
    /tmp/google-chrome-stable_current_amd64.deb

echo "Chrome remote desktop installation completed"

# Delete instance startup script do it does not re-run after a boot
gcloud compute instances remove-metadata --keys startup-script --zone us-central1-a game-client-vm
