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

BUILD_DIR=$(PWD)/bin
SCHEMA_DIR=../../infrastructure/schema/
build:
	echo "Building profile service"
	mkdir -p ${BUILD_DIR} && GOOS=linux GOARCH=386 go build -o ${BUILD_DIR}/profile-service main.go

build-docker:
	echo "Building docker container"
	docker build . -t profile-service

test-unit:
	echo "Running unit tests"
	go test -short ./...

test-integration:
	echo "Running integration tests"
	docker build . -t profile-service \
		&& mkdir -p test_data \
		&& grep -v '^--*' ${SCHEMA_DIR}/schema.sql >test_data/schema.sql \
		&& go test --tags=integration ./...

test-all: test-unit test-integration

clean:
	echo "Running cleanup"
	rm bin/*
	docker rmi -f profile-service
