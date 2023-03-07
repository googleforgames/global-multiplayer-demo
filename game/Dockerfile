# Copyright 2023 Google LLC
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
# Build a container for the Unreal Engine dedicated game server
#

FROM ghcr.io/epicgames/unreal-engine:dev-slim-5.1.0  as builder
ARG SERVER_CONFIG=Development

COPY --chown=ue4:ue4 . /tmp/project

RUN /home/ue4/UnrealEngine/Engine/Build/BatchFiles/RunUAT.sh BuildCookRun \
	-Server -NoClient -ServerConfig=${SERVER_CONFIG} \
	-Project=/tmp/project/Droidshooter.uproject \
	-UTF8Output -NoDebugInfo -AllMaps -NoP4 -Build -Cook -Stage -Pak -Package -Archive \
	-ArchiveDirectory=/tmp/project/Packaged \
	-Platform=Linux

FROM gcr.io/distroless/cc-debian11:nonroot
COPY --from=builder --chown=nonroot:nonroot /tmp/project/Packaged/LinuxServer /home/nonroot/project
EXPOSE 7777/udp
ENTRYPOINT ["/home/nonroot/project/Droidshooter/Binaries/Linux/DroidshooterServer", "Droidshooter"]
