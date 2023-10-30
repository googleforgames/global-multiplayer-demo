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

# Script to download and update the Game Client for Google Cloud Storage

set -euxo pipefail

project=$(curl http://metadata.google.internal/computeMetadata/v1/project/project-id -H Metadata-Flavor:Google)
storage_bucket="gs://$project-release-artifacts"
latest_client=$(gsutil ls -l "$storage_bucket/*.zip" | sort -k 2 -r | head -n 2 | tail -n 1 | awk '{print $3}')

mkdir -p ~/Desktop/Client || true

gsutil cp "$latest_client" ~/Desktop/Client/Client.zip
cd ~/Desktop/Client/
unzip Client.zip
