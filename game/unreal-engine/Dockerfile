# Copyright 2023 Google LLC All Rights Reserved.
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

# Credits to https://github.com/adamrehn/ue4-docker/blob/master/src/ue4docker/dockerfiles and
# https://unrealcontainers.com/docs/obtaining-images/write-your-own#writing-dockerfiles-for-linux-containers
# for how to accomplish this in a Dockerfile

# Split prerequisites into it's own layer, so we can reuse it.
FROM nvidia/opengl:1.0-glvnd-devel-ubuntu22.04 as prerequisites

# Disable interactive prompts during package installation

ENV DEBIAN_FRONTEND=noninteractive
# build prerequisites
RUN apt-get update && apt-get install -y --no-install-recommends \
		build-essential \
		ca-certificates \
		curl \
		git \
		git-lfs \
		python3 \
		python3-dev \
		python3-pip \
		shared-mime-info \
		software-properties-common \
		sudo \
		tzdata \
		unzip \
		xdg-user-dirs \
		zip \
    libasound2 \
    libatk1.0-0 \
    libatk-bridge2.0-0 \
    libcairo2 \
    libfontconfig1 \
    libfreetype6 \
    libgbm1 \
    libglu1 \
    libnss3 \
    libnspr4 \
    libpango-1.0-0 \
    libpangocairo-1.0-0 \
    libsm6 \
    libxcomposite1 \
    libxcursor1 \
    libxdamage1 \
    libxi6 \
    libxkbcommon-x11-0 \
    libxrandr2 \
    libxrender1 \
    libxss1 \
    libxtst6 \
    libxv1 \
    x11-xkb-utils \
    xauth \
    xfonts-base \
    xkb-data && \
	rm -rf /var/lib/apt/lists/*

# Unreal refuses to run as the root user, so create a non-root user with no password and allow them to run commands using sudo
RUN useradd --create-home --home /home/ue4 --shell /bin/bash --uid 1000 ue4 && \
	passwd -d ue4 && \
	usermod -a -G audio,video,sudo ue4
USER ue4

# Build Unreal Engine from source
FROM prerequisites as source

COPY --chown=ue4:ue4 ./UnrealEngine /home/ue4/UnrealEngine

WORKDIR /home/ue4/UnrealEngine
RUN whoami && ls -l .. && ls -l && ./Setup.sh -no-cache && sudo rm -rf /var/lib/apt/lists/* && ./GenerateProjectFiles.sh && make

# Do some cleanup before moving to the final layer
RUN rm -r /home/ue4/UnrealEngine/Engine/Documentation && \
    cd /home/ue4/UnrealEngine/Engine/Binaries && find -name 'Win64' -exec rm -r {} \; || true  && find -name 'Mac' -exec rm -r {} \; || true

# Copy over to the final image only what is needed
FROM prerequisites

COPY --from=source --chown=ue4:ue4 /home/ue4/UnrealEngine/Engine /home/ue4/UnrealEngine/Engine
WORKDIR /home/ue4/UnrealEngine
