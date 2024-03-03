#!/bin/bash

#
# Copyright 2024 Google LLC All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# Script to install applications for the Game Client VM

# Credits
# https://cloud.google.com/architecture/chrome-desktop-remote-on-compute-engine
# https://cloud.google.com/compute/docs/gpus/install-grid-drivers#debianubuntu

set -euxo pipefail

export DEBIAN_FRONTEND=noninteractive

# if we've already run this script, exit successfully
if [ -f /opt/game-client/init.lock ]; then
  echo "Remote Game Client VM startup script already ran - exiting"
  exit 0
fi

# Latest version of Xpra
sudo wget -O "/usr/share/keyrings/xpra.asc" https://xpra.org/xpra.asc
cd /etc/apt/sources.list.d 
sudo wget https://raw.githubusercontent.com/Xpra-org/xpra/master/packaging/repos/bookworm/xpra.sources
sudo apt-get update --assume-yes
cd -

# Installs:
# * xfce4 desktop environment
# * remoting libraries
# * graphics support and utilities
# * general handy utilities
sudo apt-get install --assume-yes \
  xfce4 desktop-base dbus-x11 xscreensaver \
  tightvncserver xrdp xpra \
  build-essential "linux-headers-$(uname -r)" mesa-utils \
  less bzip2 zip unzip tasksel wget

# Chrome Remote Desktop session
sudo bash -c 'echo "exec /etc/X11/Xsession /usr/bin/xfce4-session" > /etc/chrome-remote-desktop-session'

# Install Nvidia Virtual Workstation Drivers
curl -O https://storage.googleapis.com/nvidia-drivers-us-public/GRID/vGPU16.4/NVIDIA-Linux-x86_64-535.161.07-grid.run
sudo bash NVIDIA-Linux-x86_64-535.161.07-grid.run --silent

# Disable GSP Firmware
echo "options nvidia NVreg_EnableGpuFirmware=0" | sudo tee /etc/modprobe.d/nvidia.conf

function download_and_install { # args URL FILENAME
  curl -L -o "$2" "$1"
  sudo apt-get install --assume-yes --fix-broken "$2"
}

function is_installed {  # args PACKAGE_NAME
  dpkg-query --list "$1" | grep -q "^ii" 2>/dev/null
  return $?
}

! is_installed sunshine && \
  download_and_install \
    https://github.com/LizardByte/Sunshine/releases/download/v0.21.0/sunshine-debian-bookworm-amd64.deb \
    /tmp/sunshine-debian-bookworm-amd64.deb

! is_installed chrome-remote-desktop && \
  download_and_install \
    https://dl.google.com/linux/direct/chrome-remote-desktop_current_amd64.deb \
    /tmp/chrome-remote-desktop_current_amd64.deb

! is_installed google-chrome-stable && \
  download_and_install \
    https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb \
    /tmp/google-chrome-stable_current_amd64.deb

# Write lock file, to stop this file rerunning, and Terraform recreating the VM if the script gets deleted.
sudo mkdir /opt/game-client
sudo touch /opt/game-client/init.lock

# Grab the script that will pull down the latest Client.
project=$(curl http://metadata.google.internal/computeMetadata/v1/project/project-id -H Metadata-Flavor:Google)
storage_bucket="gs://$project-release-artifacts"

sudo gsutil cp "$storage_bucket/update-client.sh" /opt/game-client/update-client.sh
sudo chmod o+x /opt/game-client/update-client.sh

# Rebooting to disable GSP Firmware because Nvidia says so.
echo "Remote Game Client VM startup script completed - rebooting"
sudo reboot
